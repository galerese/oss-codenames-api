package domain_util

import "github.com/pkg/errors"

type GenericError struct {
	Message string
}

func (e *GenericError) Error() string {
	return e.Message
}

type StateValidationError struct {
	GenericError
}

func NewStateValidationError(msg string) *StateValidationError {
	return &StateValidationError{GenericError{
		Message: msg,
	}}
}

type UnexpectedError struct {
	GenericError
}

func NewUnexpectedError(originalError error, msg string) *UnexpectedError {
	return &UnexpectedError{GenericError{
		Message: errors.Wrap(originalError, msg).Error(),
	}}
}

type InvalidActionError struct {
	GenericError
}

func NewInvalidActionError(msg string) *InvalidActionError {
	return &InvalidActionError{GenericError{
		Message: msg,
	}}
}

type InvalidParameterError struct {
	GenericError
}

func NewInvalidParameterError(msg string) *InvalidParameterError {
	return &InvalidParameterError{GenericError{
		Message: msg,
	}}
}
