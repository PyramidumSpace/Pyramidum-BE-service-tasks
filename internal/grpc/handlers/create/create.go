package create

import (
	"context"
	"log/slog"
	"time"

	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/grpc/mapper"
	"github.com/google/uuid"
	proto "github.com/pyramidum-space/protos/gen/go/tasks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HandlerFunc = func(ctx context.Context, req *proto.CreateRequest) (*proto.CreateResponse, error)

type TaskCreator interface {
	CreateTaskContext(
		ctx context.Context,
		header string,
		text string,
		externalImages []string,
		deadline time.Time,
		progressStatus string,
		isUrgent bool,
		isImportant bool,
		ownerId int32,
		parentId uuid.UUID,
		possibleDeadline time.Time,
		weight int32,
	) (uuid.UUID, error)
}

func MakeCreateHandler(log *slog.Logger, creator TaskCreator) HandlerFunc {
	const op = "grpc.handlers.create.MakeCreateHandler"

	log = log.With(
		slog.String("op", op),
	)

	return func(ctx context.Context, req *proto.CreateRequest) (*proto.CreateResponse, error) {
		progressStatus, err := mapper.ProtoProgressStatusToString(req.ProgressStatus)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		var parentUUID uuid.UUID
		if req.ParentId == nil {
			parentUUID = uuid.Nil
		} else {
			parentUUID, err = uuid.FromBytes(req.ParentId)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}

		id, err := creator.CreateTaskContext(
			ctx,
			req.Header,
			req.Text,
			req.ExternalImages,
			req.Deadline.AsTime(),
			progressStatus,
			req.IsUrgent,
			req.IsImportant,
			req.OwnerId,
			parentUUID,
			req.PossibleDeadline.AsTime(),
			req.Weight,
		)

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		response := proto.CreateResponse{}
		response.TaskId, err = id.MarshalBinary()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &response, nil
	}
}

func bytesArrayToUUIDArray(ids [][]byte) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, 0, len(ids))

	for _, id := range ids {
		uuidFromBytes, err := uuid.FromBytes(id)
		if err != nil {
			return nil, err
		}
		uuids = append(uuids, uuidFromBytes)
	}

	return uuids, nil
}
