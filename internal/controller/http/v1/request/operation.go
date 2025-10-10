package request

import "test_go/internal/entity"

type CreateOperationRequest struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Commands    []*entity.CreateCommand `json:"commands"`
}

func (req *CreateOperationRequest) ToEntity() entity.CreateOperationInput {
	return entity.CreateOperationInput{
		Name:        req.Name,
		Description: req.Description,
		Commands:    req.Commands,
	}
}

type UpdateOperationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (req *UpdateOperationRequest) ToEntity() entity.UpdateOperationInput {
	return entity.UpdateOperationInput{
		Name:        req.Name,
		Description: req.Description,
	}
}
