package shodan

import (
	"errors"
)

var (
	ErrInvalidQuery = errors.New("Query is invalid")
	ErrBodyRead     = errors.New("Could not read error response")
)
