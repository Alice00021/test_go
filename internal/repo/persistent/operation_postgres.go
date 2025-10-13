package persistent

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"test_go/internal/entity"
	"test_go/pkg/postgres"
)

type OperationRepo struct {
	*postgres.Postgres
}

func NewOperationRepo(pg *postgres.Postgres) *OperationRepo {
	return &OperationRepo{pg}
}

func (r *OperationRepo) Create(ctx context.Context, e *entity.Operation) (*entity.Operation, error) {
	op := "OperationRepo - Create"

	sql, args, err := r.Builder.
		Insert("operations").
		Columns(
			"name, description, average_time").
		Values(
			e.Name, e.Description, e.AverageTime).
		Suffix(`RETURNING id`).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s - r.Builder: %w", op, err)
	}

	client := r.GetClient(ctx)

	var id int64
	err = client.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("%s - client.QueryRow: %w", op, err)
	}

	return r.GetById(ctx, id)
}

func (r *OperationRepo) GetById(ctx context.Context, id int64) (*entity.Operation, error) {
	op := "OperationRepo - GetById"

	sql, args, err := r.Builder.
		Select(
			"op.id", "op.created_at", "op.updated_at", "op.deleted_at",
			"op.name", "op.description", "op.average_time",
			"c.system_name", "c.default_address",
		).
		From("operations op").
		LeftJoin("operation_commands opc ON opc.operation_id = op.id").
		LeftJoin("commands c ON opc.command_id = c.id").
		Where("op.deleted_at IS NULL").
		Where(squirrel.Eq{"op.id": id}).
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

	var e entity.Operation
	e.Commands = []*entity.Command{}

	for rows.Next() {
		var (
			systemName, defaultAddress *string
		)

		if err := rows.Scan(
			&e.ID, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt,
			&e.Name, &e.Description, &e.AverageTime,
			&systemName, &defaultAddress,
		); err != nil {
			return nil, fmt.Errorf("%s - rows.Scan: %w", op, err)
		}

		if systemName != nil && defaultAddress != nil {
			e.Commands = append(e.Commands, &entity.Command{
				SystemName:     *systemName,
				DefaultAddress: entity.Address(*defaultAddress),
			})
		}
	}

	if len(e.Commands) == 0 {
		e.Commands = nil
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s - rows.Err: %w", op, err)
	}

	return &e, nil
}

func (r *OperationRepo) Update(ctx context.Context, e *entity.Operation) error {
	op := "OperationRepo - Update"

	sqlBuilder := r.Builder.
		Update("operations").
		Set("name", e.Name).
		Set("description", e.Description).
		Where(squirrel.Eq{"id": e.ID})

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

func (r *OperationRepo) DeleteById(ctx context.Context, id int64) error {
	op := "OperationRepo - Delete"

	builder := r.Builder.
		Update("operations").
		Set("deleted_at", "NOW()").
		Where("deleted_at IS NULL").
		Where(squirrel.Eq{"id": id})

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
