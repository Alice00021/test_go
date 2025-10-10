package persistent

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"test_go/internal/entity"
	"test_go/pkg/postgres"
)

type OperationCommandsRepo struct {
	*postgres.Postgres
}

func NewOperationCommandsRepo(pg *postgres.Postgres) *OperationCommandsRepo {
	return &OperationCommandsRepo{pg}
}

func (r *OperationCommandsRepo) Create(ctx context.Context, operationId int64, commandId int64) error {
	op := "OperationCommandsRepo - Create"

	sql, args, err := r.Builder.
		Insert("operation_commands").
		Columns("operation_id", "command_id").
		Values(operationId, commandId).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s - r.Builder: %w", op, err)
	}

	client := r.GetClient(ctx)
	_, err = client.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s - client.Exec: %w", op, err)
	}

	return nil
}

func (r *OperationCommandsRepo) DeleteByOperationId(ctx context.Context, operationId int64) error {
	op := "OperationCommandsRepo - Delete"

	sql, args, err := r.Builder.
		Update("operation_commands").
		Set("deleted_at", "NOW()").
		Where("deleted_at IS NULL").
		Where(squirrel.Eq{"operation_id": operationId}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s - r.Builder: %w", op, err)
	}

	client := r.GetClient(ctx)
	tag, err := client.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s - client.Exec: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return entity.ErrOperationNotFound
	}
	return nil
}
