package shodan

import (
	"bytes"
	"net"
	"strings"
)

const (
	ipPath      = "/tools/myip"
	headersPath = "/tools/httpheaders"
)

// GetMyIP returns your current IP address as seen from the Internet
// API key for this method is unnecessary
func (c *Client) GetMyIP() (net.IP, error) {
	url := c.buildBaseURL(ipPath, nil)

	var ip bytes.Buffer
	err := c.executeRequest("GET", url, &ip, nil)

	return net.ParseIP(strings.Trim(ip.String(), "\"")), err
}

// GetHTTPHeaders shows the HTTP headers that your client sends
// when connecting to a webserver.
func (c *Client) GetHTTPHeaders() (map[string]string, error) {
	url := c.buildBaseURL(headersPath, nil)

	var headers map[string]string
	err := c.executeRequest("GET", url, &headers, nil)

	return headers, err
}
