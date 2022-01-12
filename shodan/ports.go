package shodan

import (
	"context"
)

const (
	portsPath = "/shodan/ports"
)

// GetPorts returns a list of port numbers that the crawlers are looking for.
func (c *Client) GetPorts(ctx context.Context) ([]int, error) {
	var ports []int

	req, err := c.NewRequest("GET", portsPath, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &ports); err != nil {
		return nil, err
	}

	return ports, nil
}
