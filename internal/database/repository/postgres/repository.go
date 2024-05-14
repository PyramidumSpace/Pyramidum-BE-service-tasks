package postgres

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/domain/model"
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

func (r *Repository) CreateTaskContext(
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

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
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
		ExecContext(ctx)

	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt := r.pgsq.Insert("external_image").
		Columns("id", "url", "task_id")

	for _, v := range extImgs {
		stmt = stmt.Values(v.Id, v.Url, v.TaskId)
	}

	_, err = stmt.RunWith(tx).ExecContext(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return task.Id, nil
}

func (r *Repository) Task(id uuid.UUID) (*model.Task, error) {
	const op = "repository.Task"

	task := taskTable{}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	err = r.pgsq.Select("id", "header", "text", "deadline", "progress_status", "is_urgent", "is_important", "owner_id", "parent_id", "possible_deadline", "weight").
		From("task").
		Where(sq.Eq{"id": id}).
		RunWith(tx).
		QueryRow().
		Scan(
			&task.Id,
			&task.Header,
			&task.Text,
			&task.Deadline,
			&task.ProgressStatus,
			&task.IsUrgent,
			&task.IsImportant,
			&task.OwnerId,
			&task.ParentId,
			&task.PossibleDeadline,
			&task.Weight,
		)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	extImgs := make([]externalImageTable, 0)
	rows, err := r.pgsq.Select("id", "url", "task_id").
		From("external_image").
		Where(sq.Eq{"task_id": id}).
		RunWith(tx).
		Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var extImg externalImageTable
		err = rows.Scan(&extImg.Id, &extImg.Url, &extImg.TaskId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		extImgs = append(extImgs, extImg)
	}
	err = rows.Close()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var parentId uuid.UUID
	if task.ParentId == nil {
		parentId = uuid.Nil
	} else {
		parentId = *task.ParentId
	}

	images := make([]string, 0, len(extImgs))
	for _, v := range extImgs {
		images = append(images, v.Url)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &model.Task{
		Id:               task.Id,
		Header:           task.Header,
		Text:             task.Text,
		Deadline:         task.Deadline,
		ProgressStatus:   task.ProgressStatus,
		IsUrgent:         task.IsUrgent,
		IsImportant:      task.IsImportant,
		OwnerId:          task.OwnerId,
		ParentId:         parentId,
		PossibleDeadline: task.PossibleDeadline,
		Weight:           task.Weight,
		ExternalImages:   images,
	}, nil
}
