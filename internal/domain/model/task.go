package model

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Task struct {
	Id               uuid.UUID
	Header           string
	Text             string
	Deadline         time.Time
	ProgressStatus   ProgressStatus
	IsUrgent         bool
	IsImportant      bool
	OwnerId          int32
	ParentId         *uuid.UUID
	PossibleDeadline time.Time
	ExternalImages   []string
	Weight           int32
}

type ProgressStatus string

const (
	ProgressStatusInProgress ProgressStatus = "in progress"
	ProgressStatusCanceled   ProgressStatus = "canceled"
	ProgressStatusDone       ProgressStatus = "done"
)

func ProgressStatusFromString(s string) (ProgressStatus, error) {
	switch s {
	case "in progress":
		return ProgressStatusInProgress, nil
	case "canceled":
		return ProgressStatusCanceled, nil
	case "done":
		return ProgressStatusDone, nil
	default:
		return "", fmt.Errorf("unknown progress status: %s", s)
	}
}
