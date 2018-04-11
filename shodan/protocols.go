package shodan

import (
	"context"
)

const (
	protocolsPath = "/shodan/protocols"
)

// GetProtocols returns an object containing all the protocols that can be
// used when launching an Internet scan.
func (c *Client) GetProtocols(ctx context.Context) (map[string]string, error) {
	var protocols map[string]string

	req, err := c.NewRequest("GET", protocolsPath, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &protocols); err != nil {
		return nil, err
	}

	return protocols, err
}
