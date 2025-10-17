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

func (r *OperationCommandsRepo) Create(ctx context.Context, operationId int64, commands []*entity.OperationCommand) error {
	op := "OperationCommandsRepo - Create"

	builder := r.Builder.Insert("operation_commands").
		Columns("operation_id", "command_id", "address")

	for _, command := range commands {
		builder = builder.Values(operationId, command.ID, command.Address)
	}

	sql, args, err := builder.ToSql()
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

func (r *OperationCommandsRepo) Update(ctx context.Context, operationId int64, commandId int64, address entity.Address) error {
	op := "OperationCommandsRepo - Update"

	sqlBuilder := r.Builder.
		Update("operation_commands").
		Set("address", address).
		Where(squirrel.Eq{
			"operation_id": operationId,
			"command_id":   commandId})

	sql, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("%s - r.Builder: %w", op, err)
	}

	client := r.GetClient(ctx)
	if _, err = client.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("%s - client.Exec: %w", op, err)
	}

	return nil
}

func (r *OperationCommandsRepo) GetCommandIdsByOperation(ctx context.Context, operationID int64) (map[int64]struct{}, error) {
	op := "OperationCommandsRepo - GetCommandIdsByOperation"

	sql, args, err := r.Builder.
		Select("command_id").
		From("operation_commands").
		Where(squirrel.Eq{"operation_id": operationID}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s - r.Builder: %w", op, err)
	}
	client := r.GetClient(ctx)
	rows, err := client.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s - client.Query: %w", op, err)
	}
	defer rows.Close()

	commandIdsMap := make(map[int64]struct{})

	for rows.Next() {
		var commandId int64
		if err := rows.Scan(&commandId); err != nil {
			return nil, fmt.Errorf("%s - rows.Scan: %w", op, err)
		}
		commandIdsMap[commandId] = struct{}{}
	}
	return commandIdsMap, nil
}

func (r *OperationCommandsRepo) DeleteByOperationId(ctx context.Context, operationId int64) error {
	op := "OperationCommandsRepo - Delete"

	sql, args, err := r.Builder.
		Delete("operation_commands").
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

func (r *OperationCommandsRepo) DeleteIfNotInOperationIds(ctx context.Context, operationID int64, commandIds []int64) error {
	op := "OperationCommandsRepo - DeleteIfNotInAccountIds"

	builder := r.Builder.
		Delete("operation_commands").
		Where(squirrel.Eq{"operation_id": operationID})

	if len(commandIds) > 0 {
		builder = builder.Where(squirrel.NotEq{"command_id": commandIds})
	}

	sql, args, err := builder.ToSql()
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
