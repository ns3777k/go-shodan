package shodan

import (
	"bytes"
	"strings"
)

const (
	ipPath      = "/tools/myip"
	headersPath = "/tools/httpheaders"
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

// GetHttpHeaders shows the HTTP headers that your client sends
// when connecting to a webserver.
func (c *Client) GetHttpHeaders() (map[string]string, error) {
	url, err := c.buildBaseURL(headersPath, nil)
	if err != nil {
		return nil, err
	}

	var headers map[string]string
	err = c.executeRequest("GET", url, &headers, nil)

	return headers, err
}
