package getbyid

import (
	"context"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/domain/model"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/grpc/mapper"
	slogattr "github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/lib/log/slog/attr"
	proto "github.com/pyramidum-space/protos/gen/go/tasks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
)

type HandlerFunc = func(ctx context.Context, req *proto.GetByUserIDRequest) (*proto.GetByUserIDResponse, error)

type TaskProvider interface {
	TasksByUserIdContext(context.Context, int32) ([]*model.Task, error)
}

func MakeGetByUserIdHandler(log *slog.Logger, provider TaskProvider) HandlerFunc {
	const op = "grpc.handlers.getbyid.MakeGetByUserIdHandler"

	log = slog.With(
		log, slog.String("op", op),
	)

	return func(ctx context.Context, req *proto.GetByUserIDRequest) (*proto.GetByUserIDResponse, error) {
		tasks, err := provider.TasksByUserIdContext(ctx, req.UserId)
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

		return &proto.GetByUserIDResponse{Tasks: protoTasks}, nil
	}
}
