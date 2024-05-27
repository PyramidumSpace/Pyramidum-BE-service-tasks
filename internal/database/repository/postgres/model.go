package postgres

import (
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/domain/model"
	"github.com/google/uuid"
	"time"
)

type taskTable struct {
	Id               uuid.UUID
	Header           string
	Text             string
	Deadline         time.Time
	ProgressStatus   string
	IsUrgent         bool
	IsImportant      bool
	OwnerId          int32
	ParentId         *uuid.UUID
	PossibleDeadline time.Time
	Weight           int32
}

type externalImageTable struct {
	Id     int32
	Url    string
	TaskId uuid.UUID
}

type progressStatus string

const (
	progressStatusInProgress progressStatus = "in progress"
	progressStatusCanceled   progressStatus = "canceled"
	progressStatusDone       progressStatus = "done"
)

func progressStatusFromModelProgressStatus(s model.ProgressStatus) (progressStatus, error) {
	switch s {
	case model.ProgressStatusInProgress:
		return progressStatusInProgress, nil
	case model.ProgressStatusCanceled:
		return progressStatusCanceled, nil
	case model.ProgressStatusDone:
		return progressStatusDone, nil
	default:
		return "", nil
	}
}
