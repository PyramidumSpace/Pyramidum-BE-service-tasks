package postgres

import (
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
