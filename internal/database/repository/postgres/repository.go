package postgres

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"time"
)

type Repository struct {
	db   *sql.DB
	pgsq sq.StatementBuilderType
}

func NewRepository(db *sql.DB) *Repository {
	pgsq := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &Repository{
		db:   db,
		pgsq: pgsq,
	}
}

func (r *Repository) CreateTask(
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
) (uuid.UUID, error) {
	const op = "repository.CreateTask"

	id, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	var parentIdPointer *uuid.UUID
	if parentId == uuid.Nil {
		parentIdPointer = nil
	} else {
		parentIdPointer = &parentId
	}

	task := taskTable{
		Id:               id,
		Header:           header,
		Text:             text,
		Deadline:         deadline,
		ProgressStatus:   progressStatus,
		IsUrgent:         isUrgent,
		IsImportant:      isImportant,
		OwnerId:          ownerId,
		ParentId:         parentIdPointer,
		PossibleDeadline: possibleDeadline,
		Weight:           weight,
	}

	extImgs := make([]externalImageTable, 0, len(externalImages))
	for _, v := range externalImages {
		extImgs = append(extImgs, externalImageTable{
			Url:    v,
			TaskId: task.Id,
		})
	}

	tx, err := r.db.Begin()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	_, err = r.pgsq.Insert("task").
		Columns(
			"id",
			"header",
			"text",
			"deadline",
			"progress_status",
			"is_urgent",
			"is_important",
			"owner_id",
			"parent_id",
			"possible_deadline",
			"weight").
		Values(
			task.Id,
			task.Header,
			task.Text,
			task.Deadline,
			task.ProgressStatus,
			task.IsUrgent,
			task.IsImportant,
			task.OwnerId,
			task.ParentId,
			task.PossibleDeadline,
			task.Weight).
		RunWith(tx).
		Exec()

	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt := r.pgsq.Insert("external_image").
		Columns("id", "url", "task_id")

	for _, v := range extImgs {
		stmt = stmt.Values(v.Id, v.Url, v.TaskId)
	}

	_, err = stmt.RunWith(tx).Exec()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return task.Id, nil
}
