package validator

import (
	"github.com/go-playground/validator/v10"
)

type EchoValidator struct {
	validator *validator.Validate
}

func NewEchoValidator() *EchoValidator {
	return &EchoValidator{
		validator: validator.New(),
	}
}

func (ev *EchoValidator) Validate(i interface{}) error {
	return ev.validator.Struct(i)
}
