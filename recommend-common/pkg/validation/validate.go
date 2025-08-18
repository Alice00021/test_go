package validation

import (
	"github.com/go-playground/validator/v10"
	"reflect"
)

type ResultValidation struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
	Value  string `json:"value"`
	Param  string `json:"param"`
}

func NewResultvalidation(field string, reason string, value string, param string) *ResultValidation {
	return &ResultValidation{
		Field:  field,
		Reason: reason,
		Value:  value,
		Param:  param,
	}
}

type Validation struct {
	v *validator.Validate
}

func NewValidation() *Validation {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {})
}
