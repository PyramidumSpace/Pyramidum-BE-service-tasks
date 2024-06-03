package patch

import (
	"context"
	"log/slog"
	"time"

	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/domain/model"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/grpc/mapper"
	slogattr "github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/lib/log/slog/attr"
	"github.com/google/uuid"
	proto "github.com/pyramidum-space/protos/gen/go/tasks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HandlerFunc = func(ctx context.Context, req *proto.PatchRequest) (*proto.PatchResponse, error)

type TaskUpdater interface {
	PatchTaskContext(
		ctx context.Context,
		id uuid.UUID,
		header *string,
		text *string,
		externalImages []string,
		deadline time.Time,
		progressStatus *model.ProgressStatus,
		isUrgent *bool,
		isImportant *bool,
		ownerId *int32,
		parentId uuid.UUID,
		possibleDeadline time.Time,
		weight *int32,
	) error
}

func MakePatchHandler(log *slog.Logger, provider TaskUpdater) HandlerFunc {
	const op = "grpc.handlers.patch.MakePatchHandler"

	log = log.With(
		slog.String("op", op),
	)

	return func(ctx context.Context, req *proto.PatchRequest) (*proto.PatchResponse, error) {
		id, err := uuid.FromBytes(req.TaskId)
		if err != nil {
			log.Error("invalid id", slog.Any("id", req.TaskId))
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		var parentId uuid.UUID
		if req.ParentId == nil {
			parentId = uuid.Nil
		} else {
			parentIdTemp, err := uuid.FromBytes(req.ParentId)
			if err != nil {
				log.Error("invalid parent id", slog.Any("id", req.ParentId))
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
			parentId = parentIdTemp
		}

		var progressStatusPtr *model.ProgressStatus
		if req.ProgressStatus != nil {
			progressStatus, err := mapper.ProtoProgressStatusToModelProgressStatus(*req.ProgressStatus)
			if err != nil {
				log.Error("invalid progress status", slog.Any("status", req.ProgressStatus))
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
			progressStatusPtr = &progressStatus
		} else {
			progressStatusPtr = nil
		}

		var deadline time.Time
		if req.Deadline == nil {
			deadline = time.Time{}
		} else {
			deadline = req.Deadline.AsTime()
		}

		var possibleDeadline time.Time
		if req.PossibleDeadline == nil {
			possibleDeadline = time.Time{}
		} else {
			possibleDeadline = req.PossibleDeadline.AsTime()
		}

		var externalImages []string
		if req.ExternalImages != nil {
			externalImages = req.ExternalImages.ExternalImages
		} else {
			externalImages = nil
		}

		err = provider.PatchTaskContext(
			ctx,
			id,
			req.Header,
			req.Text,
			externalImages,
			deadline,
			progressStatusPtr,
			req.IsUrgent,
			req.IsImportant,
			req.OwnerId,
			parentId,
			possibleDeadline,
			req.Weight,
		)

		if err != nil {
			log.Error("error updating task", slogattr.Err(err))
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &proto.PatchResponse{
			TaskId: req.TaskId,
		}, nil
	}
}
