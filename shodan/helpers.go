package shodan

import (
	"encoding/json"
	"strconv"
)

type genericSuccessResponse struct {
	Success bool `json:"success"`
}

// IntString is string with custom unmarshaling.
type IntString string

// UnmarshalJSON handles either a string or a number
// and casts it to string.
func (v *IntString) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		*v = IntString(s)
		return nil
	}

	var n int
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}

	*v = IntString(strconv.Itoa(n))

	return nil
}

// String method just returns string out of IntString.
func (v *IntString) String() string {
	return string(*v)
}
