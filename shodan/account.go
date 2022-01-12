package shodan

import (
	"context"
)

const (
	profilePath = "/account/profile"
)

// Profile holds account's information.
type Profile struct {
	Member  bool   `json:"member"`
	Credits int    `json:"credits"`
	Name    string `json:"display_name"`
	Created string `json:"created"`
}

// GetAccountProfile returns information about the Shodan account linked to the API key.
func (c *Client) GetAccountProfile(ctx context.Context) (*Profile, error) {
	var profile Profile

	req, err := c.NewRequest("GET", profilePath, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &profile); err != nil {
		return nil, err
	}

	return &profile, err
}
