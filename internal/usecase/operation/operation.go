package operation

import (
	"context"
	"errors"
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
	op := "OperationUseCase - CreateOperation"

	var operation entity.Operation
	if err := uc.RunInTransaction(ctx, func(txCtx context.Context) error {
		e := &entity.Operation{
			Name:        inp.Name,
			Description: inp.Description,
			Commands:    []*entity.Command{},
		}

		systemNames := make([]string, 0, len(inp.Commands))
		for _, c := range inp.Commands {
			systemNames = append(systemNames, c.SystemName)
		}

		commands, err := uc.cRepo.GetBySystemNames(txCtx, systemNames)
		if err != nil {
			return fmt.Errorf("%s - uc.cRepo.GetBySystemNames: %w", op, err)
		}

		for _, command := range commands {
			for _, inputCommand := range inp.Commands {
				if command.SystemName == inputCommand.SystemName {
					command.DefaultAddress = inputCommand.Address
					e.Commands = append(e.Commands, command)
					break
				}
			}
		}

		if err := entity.ValidateCommands(e.Commands); err != nil {
			if errors.Is(err, entity.ErrCommandDuplicateAddress) {
				return err
			}
			return fmt.Errorf("%s - entity.ValidateCommands: %w", op, err)
		}

		if err := e.SumAverageTime(); err != nil {
			return fmt.Errorf("%s - e.SumAverageTime(): %w", op, err)
		}

		res, err := uc.opRepo.Create(txCtx, e)
		if err != nil {
			return fmt.Errorf("%s - uc.opRepo.Create: %w", op, err)
		}

		for _, command := range e.Commands {
			if err := uc.opcRepo.Create(txCtx, res.ID, command.ID, command.DefaultAddress); err != nil {
				return fmt.Errorf("%s - uc.opсRepo.Create: %w", op, err)
			}
		}

		operation = *res
		operation.Commands = e.Commands
		return nil
	}); err != nil {
		return nil, fmt.Errorf("%s - uc.RunInTransaction: %w", op, err)
	}

	return &operation, nil
}

func (uc *UseCase) UpdateOperation(ctx context.Context, inp entity.UpdateOperationInput) error {
	op := "OperationUseCase - UpdateOperation"

	return uc.RunInTransaction(ctx, func(txCtx context.Context) error {
		currentCommandIds, err := uc.opcRepo.GetCommandIdsByOperation(txCtx, inp.ID)
		if err != nil {
			return fmt.Errorf("%s - uc.opcRepo.GetCommandIdsByOperation: %w", op, err)
		}

		currentCommandsMap := make(map[int64]struct{}, len(currentCommandIds))
		for _, id := range currentCommandIds {
			currentCommandsMap[id] = struct{}{}
		}

		systemNames := make([]string, 0, len(inp.Commands))
		for _, c := range inp.Commands {
			systemNames = append(systemNames, c.SystemName)
		}

		commands, err := uc.cRepo.GetBySystemNames(txCtx, systemNames)
		if err != nil {
			return fmt.Errorf("%s - uc.cRepo.GetBySystemNames: %w", op, err)
		}

		commandMap := make(map[string]*entity.Command)
		for _, cmd := range commands {
			commandMap[cmd.SystemName] = cmd
		}

		var updatedCommands []*entity.Command
		newCommandIds := make([]int64, 0, len(inp.Commands))

		for _, commandInput := range inp.Commands {
			cmd, ok := commandMap[commandInput.SystemName]
			if !ok {
				return entity.ErrCommandNotFound
			}

			newCommand := *cmd
			newCommand.DefaultAddress = commandInput.Address

			updatedCommands = append(updatedCommands, &newCommand)
			newCommandIds = append(newCommandIds, newCommand.ID)
		}

		if err := entity.ValidateCommands(updatedCommands); err != nil {
			if errors.Is(err, entity.ErrCommandDuplicateAddress) {
				return err
			}
			return fmt.Errorf("%s - entity.ValidateCommands: %w", op, err)
		}

		for i, commandInput := range inp.Commands {
			commandID := updatedCommands[i].ID
			if _, exists := currentCommandsMap[commandID]; exists {
				if err := uc.opcRepo.Update(txCtx, inp.ID, commandID, commandInput.Address); err != nil {
					return fmt.Errorf("%s - uc.opcRepo.UpdateAddress: %w", op, err)
				}
			} else {
				if err := uc.opcRepo.Create(txCtx, inp.ID, commandID, commandInput.Address); err != nil {
					return fmt.Errorf("%s - uc.opcRepo.Create: %w", op, err)
				}
			}
		}

		if err := uc.opcRepo.DeleteIfNotInOperationIds(txCtx, inp.ID, newCommandIds); err != nil {
			return fmt.Errorf("%s - uc.opcRepo.DeleteIfNotInOperationIds: %w", op, err)
		}
		var totalTime int64
		for _, cmd := range updatedCommands {
			totalTime += cmd.AverageTime
		}

		op := &entity.Operation{
			Entity:      entity.Entity{ID: inp.ID},
			Name:        inp.Name,
			Description: inp.Description,
			AverageTime: totalTime,
		}

		if err := uc.opRepo.Update(txCtx, op); err != nil {
			return fmt.Errorf("%s - uc.opRepo.Update: %w", op, err)
		}

		return nil
	})
}

func (uc *UseCase) DeleteOperation(ctx context.Context, id int64) error {
	op := "OperationUseCase - DeleteOperation"

	return uc.RunInTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.opRepo.DeleteById(txCtx, id); err != nil {
			return fmt.Errorf("%s - uc.opRepo.DeleteById: %w", op, err)
		}
		if err := uc.opcRepo.DeleteByOperationId(txCtx, id); err != nil {
			return fmt.Errorf("%s - uc.opсRepo.DeleteByOperationId: %w", op, err)
		}
		return nil
	})
}
