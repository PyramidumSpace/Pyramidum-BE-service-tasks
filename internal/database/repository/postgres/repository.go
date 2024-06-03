package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/domain/model"
	"github.com/google/uuid"
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
		Columns("url", "task_id")

	for _, v := range extImgs {
		stmt = stmt.Values(v.Url, v.TaskId)
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

func (r *Repository) TaskContext(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	const op = "repository.Task"

	task := taskTable{}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
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
		QueryRowContext(ctx).
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
		QueryContext(ctx)
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

	images := make([]string, 0, len(extImgs))
	for _, v := range extImgs {
		images = append(images, v.Url)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	progressStatus, err := model.ProgressStatusFromString(task.ProgressStatus)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &model.Task{
		Id:               task.Id,
		Header:           task.Header,
		Text:             task.Text,
		Deadline:         task.Deadline,
		ProgressStatus:   progressStatus,
		IsUrgent:         task.IsUrgent,
		IsImportant:      task.IsImportant,
		OwnerId:          task.OwnerId,
		ParentId:         task.ParentId,
		PossibleDeadline: task.PossibleDeadline,
		Weight:           task.Weight,
		ExternalImages:   images,
	}, nil
}

func (r *Repository) TasksContext(
	ctx context.Context,
	ownerId int32,
	search *string,
	deadlineFrom time.Time,
	deadlineTo time.Time,
	possibleDeadlineFrom time.Time,
	possibleDeadlineTo time.Time,
	progressStatus string,
	isUrgent *bool,
	isImportant *bool,
	weightFrom *int32,
	weightTo *int32,
) ([]*model.Task, error) {
	const op = "repository.Tasks"

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := r.pgsq.Select("id", "header", "text", "deadline", "progress_status", "is_urgent", "is_important", "owner_id", "parent_id", "possible_deadline", "weight").
		From("task").
		Where(sq.Eq{"owner_id": ownerId})

	if (deadlineFrom != time.Time{}) {
		query = query.Where(sq.GtOrEq{"deadline": deadlineFrom})
	}

	if (deadlineTo != time.Time{}) {
		query = query.Where(sq.LtOrEq{"deadline": deadlineTo})
	}

	if (possibleDeadlineFrom != time.Time{}) {
		query = query.Where(sq.GtOrEq{"possible_deadline": possibleDeadlineFrom})
	}

	if (possibleDeadlineTo != time.Time{}) {
		query = query.Where(sq.GtOrEq{"possible_deadline": possibleDeadlineTo})
	}

	query = query.Where(sq.Eq{"progress_status": progressStatus})

	if isUrgent != nil {
		query = query.Where(sq.Eq{"is_urgent": *isUrgent})
	}

	if isImportant != nil {
		query = query.Where(sq.Eq{"is_important": *isImportant})
	}

	if weightFrom != nil {
		query = query.Where(sq.GtOrEq{"weight": *weightFrom})
	}

	if weightTo != nil {
		query = query.Where(sq.LtOrEq{"weight": *weightTo})
	}

	rows, err := query.
		RunWith(tx).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	tasks := make([]*model.Task, 0)
	for rows.Next() {
		task := taskTable{}
		err = rows.Scan(
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

		if task.ParentId == nil {
			task.ParentId = &uuid.Nil
		}

		progressStatus, err := model.ProgressStatusFromString(task.ProgressStatus)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		tasks = append(tasks, &model.Task{
			Id:               task.Id,
			Header:           task.Header,
			Text:             task.Text,
			Deadline:         task.Deadline,
			ProgressStatus:   progressStatus,
			IsUrgent:         task.IsUrgent,
			IsImportant:      task.IsImportant,
			OwnerId:          task.OwnerId,
			ParentId:         task.ParentId,
			PossibleDeadline: task.PossibleDeadline,
			Weight:           task.Weight,
		})
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	images := make([]string, 0)
	for _, v := range tasks {
		rows, err := r.pgsq.Select("id", "url", "task_id").
			From("external_image").
			Where(sq.Eq{"task_id": v.Id}).
			RunWith(tx).
			QueryContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		for rows.Next() {
			var extImg externalImageTable
			err = rows.Scan(&extImg.Id, &extImg.Url, &extImg.TaskId)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			images = append(images, extImg.Url)
		}

		err = rows.Close()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for i := range tasks {
		tasks[i].ExternalImages = images
	}

	return tasks, nil
}

func (r *Repository) UpdateTaskContext(ctx context.Context, task *model.Task) error {
	const op = "repository.UpdateTask"

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = r.pgsq.Update("task").
		Set("header", task.Header).
		Set("text", task.Text).
		Set("deadline", task.Deadline).
		Set("progress_status", task.ProgressStatus).
		Set("is_urgent", task.IsUrgent).
		Set("is_important", task.IsImportant).
		Set("owner_id", task.OwnerId).
		Set("parent_id", task.ParentId).
		Set("possible_deadline", task.PossibleDeadline).
		Set("weight", task.Weight).
		Where(sq.Eq{"id": task.Id}).
		RunWith(tx).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.pgsq.Delete("external_image").
		Where(sq.Eq{"task_id": task.Id}).
		RunWith(tx).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for i := range task.ExternalImages {
		_, err = r.pgsq.Insert("external_image").
			Columns("url", "task_id").
			Values(task.ExternalImages[i], task.Id).
			RunWith(tx).
			ExecContext(ctx)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repository) PatchTaskContext(
	ctx context.Context,
	id uuid.UUID,
	header *string,
	text *string,
	externalImages []string,
	deadline time.Time,
	progressStatus model.ProgressStatus,
	isUrgent *bool,
	isImportant *bool,
	ownerId *int32,
	parentId uuid.UUID,
	possibleDeadline time.Time,
	weight *int32,
) error {
	const op = "repository.PatchTask"

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := r.pgsq.Update("task")

	if header != nil {
		query = query.Set("header", header)
	}

	if text != nil {
		query = query.Set("text", text)
	}

	if (deadline != time.Time{}) {
		query = query.Set("deadline", deadline)
	}

	query = query.Set("progress_status", progressStatus)

	if isUrgent != nil {
		query = query.Set("is_urgent", isUrgent)
	}

	if isImportant != nil {
		query = query.Set("is_important", isImportant)
	}

	if ownerId != nil {
		query = query.Set("owner_id", ownerId)
	}

	if parentId != uuid.Nil {
		query = query.Set("parent_id", parentId)
	}

	if (possibleDeadline != time.Time{}) {
		query = query.Set("possible_deadline", possibleDeadline)
	}

	if weight != nil {
		query = query.Set("weight", weight)
	}

	_, err = query.
		Where(sq.Eq{"id": id}).
		RunWith(tx).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if externalImages != nil {
		_, err = r.pgsq.Delete("external_image").
			Where(sq.Eq{"task_id": id}).
			RunWith(tx).
			ExecContext(ctx)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		for i := range externalImages {
			_, err = r.pgsq.Insert("external_image").
				Columns("url", "task_id").
				Values(externalImages[i], id).
				RunWith(tx).
				ExecContext(ctx)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
