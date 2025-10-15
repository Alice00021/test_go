package entity

type Operation struct {
	Entity
	Name        string
	Description string
	AverageTime int64
	Commands    []*Command
}

type UpdateOperationInput struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Commands    []*CommandInput `json:"commands"`
}

type CommandInput struct {
	SystemName string
	Address    Address
}

type CreateOperationInput struct {
	Name        string
	Description string
	AverageTime int64
	Commands    []*CommandInput
}

func (e *Operation) SumAverageTime() error {
	var total int64
	for _, c := range e.Commands {
		total += c.AverageTime
	}
	e.AverageTime = total
	return nil
}
