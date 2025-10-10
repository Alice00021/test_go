package operation

import (
	"context"
	"fmt"
	"test_go/internal/entity"
	"test_go/internal/repo"
	"test_go/pkg/logger"
	"test_go/pkg/transactional"
)

type UseCase struct {
	transactional.Transactional
	opRepo  repo.OperationRepo
	opсRepo repo.OperationCommandsRepo
	cRepo   repo.CommandRepo
	l       logger.Interface
}

func New(
	t transactional.Transactional,
	opRepo repo.OperationRepo,
	opCmdRepo repo.OperationCommandsRepo,
	cmdRepo repo.CommandRepo,
	l logger.Interface,
) *UseCase {
	return &UseCase{
		Transactional: t,
		opRepo:        opRepo,
		opсRepo:       opCmdRepo,
		cRepo:         cmdRepo,
		l:             l,
	}
}

func (uc *UseCase) CreateOperation(ctx context.Context, inp entity.CreateOperationInput) (*entity.Operation, error) {
	opName := "CommandUseCase - CreateOperation"

	var createdOp *entity.Operation

	if err := uc.RunInTransaction(ctx, func(txCtx context.Context) error {

		operation := &entity.Operation{
			Name:        inp.Name,
			Description: inp.Description,
			Commands:    make([]*entity.Command, len(inp.Commands)),
		}

		for i, c := range inp.Commands {
			cmd, err := uc.cRepo.GetBySystemName(txCtx, c.SystemName)
			if err != nil {
				return fmt.Errorf("%s - cRepo.GetBySystemName: %w", opName, err)
			}
			cmd.DefaultAddress = c.Address
			operation.Commands[i] = cmd
		}

		operation.SumAverageTime()

		operationId, err := uc.opRepo.Create(txCtx, operation)
		if err != nil {
			return fmt.Errorf("%s - opRepo.Create: %w", opName, err)
		}

		for _, command := range operation.Commands {
			if err := uc.opсRepo.Create(txCtx, operationId.ID, command.ID); err != nil {
				return fmt.Errorf("%s - opсRepo.Create: %w", opName, err)
			}
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("%s - uc.RunInTransaction: %w", opName, err)
	}

	return createdOp, nil
}

func (uc *UseCase) UpdateOperation(ctx context.Context, inp entity.UpdateOperationInput) error {
	opName := "UseCase - UpdateOperation"

	return uc.RunInTransaction(ctx, func(txCtx context.Context) error {
		e := &entity.Operation{
			Entity:      entity.Entity{ID: inp.ID},
			Name:        inp.Name,
			Description: inp.Description,
		}
		if err := uc.opRepo.Update(txCtx, e); err != nil {
			return fmt.Errorf("%s - opRepo.Update: %w", opName, err)
		}
		return nil
	})
}

func (uc *UseCase) DeleteOperation(ctx context.Context, id int64) error {
	opName := "UseCase - DeleteOperation"

	return uc.RunInTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.opRepo.DeleteById(txCtx, id); err != nil {
			return fmt.Errorf("%s - opRepo.DeleteById: %w", opName, err)
		}
		if err := uc.opсRepo.DeleteByOperationId(txCtx, id); err != nil {
			return fmt.Errorf("%s - opсRepo.DeleteByOperationId: %w", opName, err)
		}
		return nil
	})
}
