package apperrors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrAlreadyExists      = errors.New("already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInternal           = errors.New("internal error")
	ErrNotImplemented     = errors.New("not implemented")

	//user
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserNotFound     = errors.New("user not found")

	//token
	ErrUnexpectedStatusCode  = errors.New("unexpected status code")
	ErrTokenGeneration       = errors.New("token generation failed")
	ErrTokenNotFound         = errors.New("token not found in header")
	ErrInvalidOrExpiredToken = errors.New("invalid signature or token expired")
)

type NotFoundError struct {
	Entity string
	ID     any
}

func NewNotFoundError(entity string, id any) NotFoundError {
	return NotFoundError{
		Entity: entity,
		ID:     id,
	}
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s with id %v not found", e.Entity, e.ID)
}

// Unwrap, чтобы errors.Is(err, ErrNotFound) работало и для структуры NotFoundError т.к. errors.Is вызывает err.Unwrap()
func (e NotFoundError) Unwrap() error {
	return ErrNotFound
}
