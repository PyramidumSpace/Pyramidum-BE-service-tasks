package update

import (
	"context"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/domain/model"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/grpc/mapper"
	slogattr "github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/lib/log/slog/attr"
	"github.com/google/uuid"
	proto "github.com/pyramidum-space/protos/gen/go/tasks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type HandlerFunc = func(ctx context.Context, req *proto.UpdateRequest) (*proto.UpdateResponse, error)

type TaskUpdater interface {
	UpdateTaskContext(context.Context, *model.Task) error
}

func MakeUpdateHandler(log *slog.Logger, provider TaskUpdater) HandlerFunc {
	const op = "grpc.handlers.update.MakeUpdateHandler"

	log = slog.With(
		log, slog.String("op", op),
	)

	return func(ctx context.Context, req *proto.UpdateRequest) (*proto.UpdateResponse, error) {
		id, err := uuid.FromBytes(req.Task.Id)
		if err != nil {
			log.Error("invalid id", slog.Any("id", req.Task.Id))
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		var parentId *uuid.UUID
		if req.Task.ParentId == nil {
			parentId = nil
		} else {
			parentIdTemp, err := uuid.FromBytes(req.Task.ParentId)
			if err != nil {
				log.Error("invalid parent id", slog.Any("id", req.Task.ParentId))
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
			parentId = &parentIdTemp
		}

		progressStatus, err := mapper.ProtoProgressStatusToModelProgressStatus(req.Task.ProgressStatus)
		if err != nil {
			log.Error("invalid progress status", slog.Any("status", req.Task.ProgressStatus))
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		mTask := model.Task{
			Id:               id,
			Header:           req.Task.Header,
			Text:             req.Task.Text,
			ExternalImages:   req.Task.ExternalImages,
			Deadline:         req.Task.Deadline.AsTime(),
			ProgressStatus:   progressStatus,
			IsUrgent:         req.Task.IsUrgent,
			IsImportant:      req.Task.IsImportant,
			OwnerId:          req.Task.OwnerId,
			ParentId:         parentId,
			PossibleDeadline: req.Task.PossibleDeadline.AsTime(),
			Weight:           req.Task.Weight,
		}

		if err := provider.UpdateTaskContext(ctx, &mTask); err != nil {
			log.Error("error updating task", slogattr.Err(err))
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &proto.UpdateResponse{
			TaskId: req.Task.Id,
		}, nil
	}
}
