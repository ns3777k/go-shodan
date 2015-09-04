package shodan

import (
	"strings"
)

type DnsResolved map[string]string
type DnsReversed map[string][]string

const (
	resolvePath = "/dns/resolve"
	reversePath = "/dns/reverse"
)

func (c *Client) GetDnsResolve(hostnames []string) (*DnsResolved, error) {
	params := QueryStringParams{"hostnames": strings.Join(hostnames, ",")}
	url, err := c.buildUrl(resolvePath, params)
	if err != nil {
		return nil, err
	}

	var dnsResolved DnsResolved
	err = c.executeRequest("GET", url, &dnsResolved)

	return &dnsResolved, err
}

func (c *Client) GetDnsReverse(ips []string) (*DnsReversed, error) {
	params := QueryStringParams{"ips": strings.Join(ips, ",")}
	url, err := c.buildUrl(reversePath, params)
	if err != nil {
		return nil, err
	}

	var dnsReversed DnsReversed
	err = c.executeRequest("GET", url, &dnsReversed)

	return &dnsReversed, err
}
