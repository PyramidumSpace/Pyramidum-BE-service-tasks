package model

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	Id               uuid.UUID
	Header           string
	Text             string
	Deadline         time.Time
	ProgressStatus   string
	IsUrgent         bool
	IsImportant      bool
	OwnerId          int32
	ParentId         uuid.UUID
	PossibleDeadline time.Time
	ExternalImages   []string
	Weight           int32
}
