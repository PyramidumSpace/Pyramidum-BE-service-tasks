package mapper

import (
	"fmt"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/domain/model"
	proto "github.com/pyramidum-space/protos/gen/go/tasks"
)

func ProtoProgressStatusToString(status proto.ProgressStatus) (string, error) {
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

func StringToProgressStatus(s string) (proto.ProgressStatus, error) {
	switch s {
	case "in progress":
		return proto.ProgressStatus_PROGRESS_STATUS_IN_PROGRESS, nil
	case "canceled":
		return proto.ProgressStatus_PROGRESS_STATUS_CANCELED, nil
	case "done":
		return proto.ProgressStatus_PROGRESS_STATUS_DONE, nil
	default:
		return 0, fmt.Errorf("unknown progress status: %s", s)
	}
}

func ProtoProgressStatusToModelProgressStatus(status proto.ProgressStatus) (model.ProgressStatus, error) {
	switch status {
	case proto.ProgressStatus_PROGRESS_STATUS_CANCELED:
		return model.ProgressStatusCanceled, nil
	case proto.ProgressStatus_PROGRESS_STATUS_IN_PROGRESS:
		return model.ProgressStatusInProgress, nil
	case proto.ProgressStatus_PROGRESS_STATUS_DONE:
		return model.ProgressStatusDone, nil
	default:
		return "", fmt.Errorf("unknown progress status: %d", status)
	}
}

func ModelProgressStatusToProtoProgressStatus(status model.ProgressStatus) (proto.ProgressStatus, error) {
	switch status {
	case model.ProgressStatusCanceled:
		return proto.ProgressStatus_PROGRESS_STATUS_CANCELED, nil
	case model.ProgressStatusInProgress:
		return proto.ProgressStatus_PROGRESS_STATUS_IN_PROGRESS, nil
	case model.ProgressStatusDone:
		return proto.ProgressStatus_PROGRESS_STATUS_DONE, nil
	default:
		return proto.ProgressStatus_PROGRESS_STATUS_IN_PROGRESS, fmt.Errorf("unknown progress status: %s", status)
	}
}
