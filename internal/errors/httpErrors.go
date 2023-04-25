package errors

import (
	"errors"
	"fmt"
)

type internalServerError struct {
	cause error
}

type unauthorizedError struct {
	cause error
}

type badRequestError struct {
	cause error
}

type notFoundError struct {
	cause error
}

type forbiddenError struct {
	cause error
}

func (e internalServerError) Error() string {
	return e.cause.Error()
}

func (e unauthorizedError) Error() string {
	return e.cause.Error()
}

func (e badRequestError) Error() string {
	return e.cause.Error()
}

func (e notFoundError) Error() string {
	return e.cause.Error()
}

func (e forbiddenError) Error() string {
	return e.cause.Error()
}

func NewInternalServerError(e error) internalServerError {
	return internalServerError{cause: e}
}

func NewUnauthorizedError(e error) unauthorizedError {
	return unauthorizedError{cause: e}
}

func NewBadRequestError(e error) badRequestError {
	return badRequestError{cause: e}
}

func NewNotFoundError(e error) notFoundError {
	return notFoundError{cause: e}
}

func NewForbiddenError(e error) forbiddenError {
	return forbiddenError{cause: e}
}

func NewInternalServerErrorWrapped(e error, message string) internalServerError {
	err := fmt.Errorf("%v: %w ", message, e)
	return internalServerError{cause: err}
}

func NewUnauthorizedErrorWrapped(e error, message string) unauthorizedError {
	err := fmt.Errorf("%v: %w ", message, e)
	return unauthorizedError{cause: err}
}

func NewBadRequestErrorWrapped(e error, message string) badRequestError {
	err := fmt.Errorf("%v: %w ", message, e)
	return badRequestError{cause: err}
}

func NewNotFoundErrorWrapped(e error, message string) notFoundError {
	err := fmt.Errorf("%v: %w ", message, e)
	return notFoundError{cause: err}
}

func NewForbiddenErrorWrapped(e error, message string) forbiddenError {
	err := fmt.Errorf("%v: %w ", message, e)
	return forbiddenError{cause: err}
}

func IsInternalServerError(err error) bool {
	return errors.As(err, &internalServerError{})
}

func IsUnauthorizedError(err error) bool {
	return errors.As(err, &unauthorizedError{})
}

func IsBadRequestError(err error) bool {
	return errors.As(err, &badRequestError{})
}

func IsNotFoundError(err error) bool {
	return errors.As(err, &notFoundError{})
}

func IsForbiddenError(err error) bool {
	return errors.As(err, &forbiddenError{})
}
