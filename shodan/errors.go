package shodan

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	// ErrInvalidQuery is returned when query is not valid.
	ErrInvalidQuery = errors.New("query is invalid")

	// ErrBodyRead is returned when response's body cannot be read.
	ErrBodyRead = errors.New("could not read error response")
)

func getErrorFromResponse(r *http.Response) error {
	errorResponse := new(struct {
		Error string `json:"error"`
	})
	message, err := ioutil.ReadAll(r.Body)
	if err == nil {
		if err := json.Unmarshal(message, errorResponse); err == nil {
			return errors.New(errorResponse.Error)
		}

		return errors.New(strings.TrimSpace(string(message)))
	}

	return ErrBodyRead
}
