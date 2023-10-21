package models

import "errors"

var (
	ErrResourceNotFound      = errors.New("resource not found")
	ErrInvalidResourceSyntax = errors.New("resource syntax invalid")
)
