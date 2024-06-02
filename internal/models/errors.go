package models

import (
	"errors"
)

var (
	ErrResourceNotFound      = errors.New("resource not found")
	ErrInvalidResourceSyntax = errors.New("resource syntax invalid")
	ErrRepository            = errors.New("repository error")
	ErrInvalidId             = errors.New("invalid id")
)
