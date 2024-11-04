package api

import (
	"github.com/katatrina/simplebank/validator"
)

type badRequest struct {
	Message         string                      `json:"message"`
	FieldViolations []*validator.FieldViolation `json:"field_violations"`
}

func fieldViolation(field string, err error) *validator.FieldViolation {
	return &validator.FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func invalidArgumentError(violations []*validator.FieldViolation) *badRequest {
	return &badRequest{
		Message:         "Invalid argument(s) for the request",
		FieldViolations: violations,
	}
}
