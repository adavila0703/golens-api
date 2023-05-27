package api

import "github.com/pkg/errors"

var (
	ErrMissingFields     = errors.New("Missing required fields.")
	ErrCannotParseClaims = errors.New("Could not parse claims.")
	ErrTokenExpired      = errors.New("Token has expired")
	BadRequest           = errors.New("Bad request")
)
