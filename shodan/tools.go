package shodan

import (
	"bytes"
	"net"
	"strings"
	"context"
)

const (
	ipPath      = "/tools/myip"
	headersPath = "/tools/httpheaders"
)

// GetMyIP returns your current IP address as seen from the Internet
// API key for this method is unnecessary.
func (c *Client) GetMyIP(ctx context.Context) (net.IP, error) {
	var ip bytes.Buffer

	req, err := c.NewRequest("GET", ipPath, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &ip); err != nil {
		return nil, err
	}

	return net.ParseIP(strings.Trim(ip.String(), "\"")), nil
}

// GetHTTPHeaders shows the HTTP headers that your client sends
// when connecting to a webserver.
func (c *Client) GetHTTPHeaders(ctx context.Context) (map[string]string, error) {
	var headers map[string]string

	req, err := c.NewRequest("GET", headersPath, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &headers); err != nil {
		return nil, err
	}

	return headers, nil
}
