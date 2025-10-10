package command

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"test_go/internal/entity"
	"test_go/internal/repo"
	"test_go/pkg/logger"
	"test_go/pkg/transactional"
)

type useCase struct {
	transactional.Transactional
	repo repo.CommandRepo
	l    logger.Interface
}

func New(t transactional.Transactional,
	repo repo.CommandRepo,
	l logger.Interface,
) *useCase {
	return &useCase{
		Transactional: t,
		repo:          repo,
		l:             l,
	}
}

func (uc *useCase) UpdateCommands(ctx context.Context, fileHeader *multipart.FileHeader) error {
	op := "CommandUseCase - UpdateCommands"

	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("%s - fileHeader.Open: %w", op, err)
	}
	defer file.Close()

	var rawCommands []map[string]interface{}
	if err := json.NewDecoder(file).Decode(&rawCommands); err != nil {
		return fmt.Errorf("%s - json.Decode: %w", op, err)
	}

	var commands []entity.Command

	for _, rawCmd := range rawCommands {
		nameArray, ok := rawCmd["name"].([]interface{})
		if !ok {
			return fmt.Errorf("%s - invalid name format", op)
		}
		var name string
		for _, n := range nameArray {
			nMap, ok := n.(map[string]interface{})
			if !ok {
				continue
			}
			locale, ok := nMap["locale"].(string)
			if !ok || locale != "en" {
				continue
			}
			value, ok := nMap["value"].(string)
			if ok {
				name = value
				break
			}
		}
		rawCmd["name"] = name

		cmdBytes, err := json.Marshal(rawCmd)
		if err != nil {
			return fmt.Errorf("%s - json.Marshal: %w", op, err)
		}
		var cmd entity.Command
		if err := json.Unmarshal(cmdBytes, &cmd); err != nil {
			return fmt.Errorf("%s - json.Unmarshal: %w", op, err)
		}
		commands = append(commands, cmd)
	}

	if err := uc.RunInTransaction(ctx, func(txCtx context.Context) error {
		for _, c := range commands {
			existCommand, err := uc.repo.GetBySystemName(txCtx, string(c.SystemName))
			if err != nil {
				if err == entity.ErrCommandNotFound {
					if _, err := uc.repo.Create(txCtx, &c); err != nil {
						return fmt.Errorf("%s - repo.Create: %w", op, err)
					}
					continue
				}
				return fmt.Errorf("%s - repo.GetBySystemName: %w", op, err)
			}

			if existCommand.DefaultAddress != c.DefaultAddress {
				if err := uc.repo.Update(txCtx, existCommand.ID, c.DefaultAddress); err != nil {
					return fmt.Errorf("%s - repo.Update: %w", op, err)
				}
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("%s - uc.RunInTransaction: %w", op, err)
	}
	return nil
}
