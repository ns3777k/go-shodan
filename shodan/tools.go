package shodan

import (
	"bytes"
	"strings"
)

const (
	ipPath = "/tools/myip"
)

func (c *Client) GetMyIP() (string, error) {
	url, err := c.buildUrl(ipPath, nil)
	if err != nil {
		return "", err
	}

	var ip bytes.Buffer
	err = c.executeRequest("GET", url, &ip)

	return strings.Trim(ip.String(), "\""), err
}
