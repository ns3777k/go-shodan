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
func (c *Client) CalcHoneyScore(ip string) (float64, error) {
	var score float64

	if parsedIP := net.ParseIP(ip); parsedIP == nil {
		return 0.0, &net.ParseError{
			Type: "IP address",
			Text: ip,
		}
	}

	path := fmt.Sprintf(honeyscorePath, ip)
	url := c.buildBaseURL(path, nil)
	err := c.executeRequest("GET", url, &score, nil)

	return score, err
}
