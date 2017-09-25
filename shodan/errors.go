package shodan

import (
	"errors"
)

var (
	// ErrInvalidQuery is returned when query is not valid.
	ErrInvalidQuery = errors.New("query is invalid")

	// ErrBodyRead is returned when response's body cannot be read.
	ErrBodyRead = errors.New("could not read error response")
)
