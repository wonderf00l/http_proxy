package errors

import "fmt"

type Type uint8

type Layer string

const (
	Repo     Layer = "Repository"
	Delivery Layer = "Delivery"
)

const (
	_ Type = iota
	ErrNotFound
	ErrAlreadyExists
	ErrInvalidInput
)

type DeclaredError interface {
	Type() Type
}

// general application errors

type InternalError struct {
	Message string
	Layer   string
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("Internal error occured. Message: '%s'. Layer: %s", e.Message, e.Layer)
}
