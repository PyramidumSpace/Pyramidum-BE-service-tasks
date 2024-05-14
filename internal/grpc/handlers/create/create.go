package create

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	proto "github.com/pyramidum-space/protos/gen/go/tasks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"time"
)

type HandlerFunc = func(ctx context.Context, req *proto.CreateRequest) (*proto.CreateResponse, error)

type TaskCreator interface {
	CreateTask(
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

	log = slog.With(
		log, slog.String("op", op),
	)

	return func(ctx context.Context, req *proto.CreateRequest) (*proto.CreateResponse, error) {
		progressStatus, err := progressStatusToString(req.ProgressStatus)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		parentUUID, err := uuid.FromBytes(req.ParentId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		id, err := creator.CreateTask(
			req.Header,
			req.Text,
			req.ExternalImages,
			protobufTimestampToTime(req.Deadline),
			progressStatus,
			req.IsUrgent,
			req.IsImportant,
			req.OwnerId,
			parentUUID,
			protobufTimestampToTime(req.PossibleDeadline),
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

func protobufTimestampToTime(ts *timestamppb.Timestamp) time.Time {
	return time.Unix(ts.Seconds, int64(ts.Nanos))
}

func progressStatusToString(status proto.ProgressStatus) (string, error) {
	switch status {
	case proto.ProgressStatus_PROGRESS_STATUS_CANCELED:
		return "canceled", nil
	case proto.ProgressStatus_PROGRESS_STATUS_IN_PROGRESS:
		return "in progress", nil
	case proto.ProgressStatus_PROGRESS_STATUS_DONE:
		return "done", nil
	default:
		return "", fmt.Errorf("unknown progress status: %d", status)
	}
}
