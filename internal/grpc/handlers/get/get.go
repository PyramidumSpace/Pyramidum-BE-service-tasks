package get

import (
	"context"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/domain/model"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/grpc/mapper"
	slogattr "github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/lib/log/slog/attr"
	"github.com/google/uuid"
	proto "github.com/pyramidum-space/protos/gen/go/tasks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
)

type HandlerFunc = func(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error)

type TaskProvider interface {
	TaskContext(context.Context, uuid.UUID) (*model.Task, error)
}

func MakeGetHandler(log *slog.Logger, provider TaskProvider) HandlerFunc {
	const op = "grpc.handlers.get.MakeGetHandler"

	log = slog.With(
		log, slog.String("op", op),
	)

	return func(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
		id, err := uuid.FromBytes(req.TaskId)
		if err != nil {
			log.Error("invalid id", slog.Any("id", req.TaskId))
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		task, err := provider.TaskContext(ctx, id)
		if err != nil {
			log.Error("error getting task", slogattr.Err(err))
			return nil, status.Error(codes.Internal, err.Error())
		}

		protoId, err := task.Id.MarshalBinary()
		if err != nil {
			log.Error("error converting task", slogattr.Err(err))
			return nil, status.Error(codes.Internal, err.Error())
		}

		progressStatus, err := mapper.ModelProgressStatusToProtoProgressStatus(task.ProgressStatus)
		if err != nil {
			log.Error("error converting task", slogattr.Err(err))
			return nil, status.Error(codes.Internal, err.Error())
		}

		parentId, err := task.ParentId.MarshalBinary()
		if err != nil {
			log.Error("error converting task", slogattr.Err(err))
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &proto.GetResponse{
			Task: &proto.Task{
				Id:               protoId,
				Header:           task.Header,
				Text:             task.Text,
				ExternalImages:   task.ExternalImages,
				Deadline:         timestamppb.New(task.Deadline),
				ProgressStatus:   progressStatus,
				IsUrgent:         task.IsUrgent,
				IsImportant:      task.IsImportant,
				OwnerId:          task.OwnerId,
				ParentId:         parentId,
				PossibleDeadline: timestamppb.New(task.PossibleDeadline),
				Weight:           task.Weight,
			},
		}, nil
	}
}
