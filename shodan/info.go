package shodan

import (
	"context"
)

const (
	infoPath = "/api-info"
)

// APIInfo holds API information.
type APIInfo struct {
	QueryCredits int    `json:"query_credits"`
	ScanCredits  int    `json:"scan_credits"`
	Telnet       bool   `json:"telnet"`
	Plan         string `json:"plan"`
	HTTPS        bool   `json:"https"`
	Unlocked     bool   `json:"unlocked"`
	UnlockedLeft int    `json:"unlocked_left"`
}

// GetAPIInfo returns information about the API plan belonging to the given API key.
func (c *Client) GetAPIInfo(ctx context.Context) (*APIInfo, error) {
	var apiInfo APIInfo

	req, err := c.NewRequest("GET", infoPath, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &apiInfo); err != nil {
		return nil, err
	}

	return &apiInfo, nil
}
