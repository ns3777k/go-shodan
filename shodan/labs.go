package shodan

import (
	"fmt"
	"net"
)

const (
	honeyscorePath = "/labs/honeyscore/%s"
)

// CalcHoneyScore calculates a honeypot probability score ranging from
// 0 (not a honeypot) to 1.0 (is a honeypot)
func (c *Client) CalcHoneyScore(ip net.IP) (float64, error) {
	var score float64

	path := fmt.Sprintf(honeyscorePath, ip.String())
	url := c.buildBaseURL(path, nil)
	err := c.executeRequest("GET", url, &score, nil)

	return score, err
}
