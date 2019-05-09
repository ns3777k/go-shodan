package shodan

import (
	"context"
	"fmt"
	"net"
)

const (
	honeyscorePath = "/labs/honeyscore/%s"
)

// CalcHoneyScore calculates a honeypot probability score ranging from
// 0 (not a honeypot) to 1.0 (is a honeypot).
func (c *Client) CalcHoneyScore(ctx context.Context, ip net.IP) (float64, error) { //nolint:interfacer
	var score float64

	path := fmt.Sprintf(honeyscorePath, ip.String())
	req, err := c.NewRequest("GET", path, nil, nil)
	if err != nil {
		return 0.0, err
	}

	if err := c.Do(ctx, req, &score); err != nil {
		return 0.0, err
	}

	return score, err
}
