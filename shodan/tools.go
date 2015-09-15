package shodan

import (
	"bytes"
	"strings"
)

const (
	ipPath = "/tools/myip"
)

// GetMyIP returns your current IP address as seen from the Internet
// API key for this method is unnecessary
func (c *Client) GetMyIP() (string, error) {
	url, err := c.buildBaseURL(ipPath, nil)
	if err != nil {
		return "", err
	}

	var ip bytes.Buffer
	err = c.executeRequest("GET", url, &ip, nil)

	return strings.Trim(ip.String(), "\""), err
}
