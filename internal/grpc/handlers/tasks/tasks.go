package tasks

import (
	"context"
	"log/slog"

	"time"

	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/domain/model"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/grpc/mapper"
	slogattr "github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/lib/log/slog/attr"
	proto "github.com/pyramidum-space/protos/gen/go/tasks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type HandlerFunc = func(ctx context.Context, req *proto.TasksRequest) (*proto.TasksResponse, error)

type TaskProvider interface {
	TasksContext(
		ctx context.Context,
		ownerId int32,
		search *string,
		deadlineFrom time.Time,
		deadlineTo time.Time,
		possibleDeadlineFrom time.Time,
		possibleDeadlineTo time.Time,
		progressStatus *string,
		isUrgent *bool,
		isImportant *bool,
		weightFrom *int32,
		weightTo *int32,
	) ([]*model.Task, error)
}

func MakeTasksHandler(log *slog.Logger, provider TaskProvider) HandlerFunc {
	const op = "grpc.handlers.tasks.MakeTasksHandler"

	log = log.With(
		slog.String("op", op),
	)

	return func(ctx context.Context, req *proto.TasksRequest) (*proto.TasksResponse, error) {
		var progressStatusPtr *string
		if req.ProgressStatus != nil {
			progressStatus, err := mapper.ProtoProgressStatusToString(*req.ProgressStatus)
			if err != nil {
				log.Error("invalid progress status", slog.Any("status", req.ProgressStatus))
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
			progressStatusPtr = &progressStatus
		} else {
			progressStatusPtr = nil
		}

		var deadlineFrom time.Time
		if req.DeadlineFrom == nil {
			deadlineFrom = time.Time{}
		} else {
			deadlineFrom = req.DeadlineFrom.AsTime()
		}

		var deadlineTo time.Time
		if req.DeadlineTo == nil {
			deadlineTo = time.Time{}
		} else {
			deadlineTo = req.DeadlineTo.AsTime()
		}

		var possibleDeadlineFrom time.Time
		if req.PossibleDeadlineFrom == nil {
			possibleDeadlineFrom = time.Time{}
		} else {
			possibleDeadlineFrom = req.PossibleDeadlineFrom.AsTime()
		}

		var possibleDeadlineTo time.Time
		if req.PossibleDeadlineTo == nil {
			possibleDeadlineTo = time.Time{}
		} else {
			possibleDeadlineTo = req.PossibleDeadlineTo.AsTime()
		}

		tasks, err := provider.TasksContext(
			ctx,
			req.OwnerId,
			req.Search,
			deadlineFrom,
			deadlineTo,
			possibleDeadlineFrom,
			possibleDeadlineTo,
			progressStatusPtr,
			req.IsUrgent,
			req.IsImportant,
			req.WeightFrom,
			req.WeightTo,
		)
		if err != nil {
			log.Error("error getting tasks", slogattr.Err(err))
			return nil, status.Error(codes.Internal, err.Error())
		}

		protoTasks := make([]*proto.Task, 0, len(tasks))

		for _, v := range tasks {
			idBytes, err := v.Id.MarshalBinary()
			if err != nil {
				log.Error("error converting task", slogattr.Err(err))
				return nil, status.Error(codes.Internal, err.Error())
			}

			protoProgressStatus, err := mapper.ModelProgressStatusToProtoProgressStatus(v.ProgressStatus)
			if err != nil {
				log.Error("error converting task", slogattr.Err(err))
				return nil, status.Error(codes.Internal, err.Error())
			}

			parentIdBytes, err := v.ParentId.MarshalBinary()
			if err != nil {
				log.Error("error converting task", slogattr.Err(err))
				return nil, status.Error(codes.Internal, err.Error())
			}

			protoTasks = append(protoTasks, &proto.Task{
				Id:               idBytes,
				Header:           v.Header,
				Text:             v.Text,
				ExternalImages:   v.ExternalImages,
				Deadline:         timestamppb.New(v.Deadline),
				ProgressStatus:   protoProgressStatus,
				IsUrgent:         v.IsUrgent,
				IsImportant:      v.IsImportant,
				OwnerId:          v.OwnerId,
				ParentId:         parentIdBytes,
				PossibleDeadline: timestamppb.New(v.PossibleDeadline),
				Weight:           v.Weight,
			})
		}

		return &proto.TasksResponse{Tasks: protoTasks}, nil
	}
}
