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
	opcRepo repo.OperationCommandsRepo
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
		opcRepo:       opCmdRepo,
		cRepo:         cmdRepo,
		l:             l,
	}
}

func (uc *UseCase) CreateOperation(ctx context.Context, inp entity.CreateOperationInput) (*entity.Operation, error) {
	opName := "OperationUseCase - CreateOperation"

	var operation entity.Operation
	if err := uc.RunInTransaction(ctx, func(txCtx context.Context) error {
		e := &entity.Operation{
			Name:        inp.Name,
			Description: inp.Description,
			Commands:    []*entity.Command{},
		}

		for _, c := range inp.Commands {
			command, err := uc.cRepo.GetBySystemName(txCtx, c.SystemName)
			if err != nil {
				return fmt.Errorf("%s - cRepo.GetBySystemName: %w", opName, err)
			}
			command.DefaultAddress = c.Address
			e.Commands = append(e.Commands, command)
		}

		if err := e.SumAverageTime(); err != nil {
			return fmt.Errorf("%s - operation.SumAverageTime(): %w", opName, err)
		}

		res, err := uc.opRepo.Create(txCtx, e)
		if err != nil {
			return fmt.Errorf("%s - opRepo.Create: %w", opName, err)
		}

		for _, command := range e.Commands {
			if err := uc.opcRepo.Create(txCtx, res.ID, command.ID); err != nil {
				return fmt.Errorf("%s - opсRepo.Create: %w", opName, err)
			}
		}

		operation = *res
		operation.Commands = e.Commands
		return nil
	}); err != nil {
		return nil, fmt.Errorf("%s - uc.RunInTransaction: %w", opName, err)
	}

	return &operation, nil
}

func (uc *UseCase) UpdateOperation(ctx context.Context, inp entity.UpdateOperationInput) error {
	opName := "OperationUseCase - UpdateOperation"

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
	opName := "OperationUseCase - DeleteOperation"

	return uc.RunInTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.opRepo.DeleteById(txCtx, id); err != nil {
			return fmt.Errorf("%s - opRepo.DeleteById: %w", opName, err)
		}
		if err := uc.opcRepo.DeleteByOperationId(txCtx, id); err != nil {
			return fmt.Errorf("%s - opсRepo.DeleteByOperationId: %w", opName, err)
		}
		return nil
	})
}
