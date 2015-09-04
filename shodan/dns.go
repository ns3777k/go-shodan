package shodan

import (
	"strings"
)

type DnsResolveResult map[string]string
type DnsReverseResult map[string][]string

const (
	resolvePath = "/dns/resolve"
	reversePath = "/dns/reverse"
)

func (c *Client) GetDnsResolve(hostnames []string) (*DnsResolveResult, error) {
	params := map[string]string{"hostnames": strings.Join(hostnames, ",")}
	url, err := c.buildUrl(resolvePath, params)
	if err != nil {
		return nil, err
	}

	var dnsResolved DnsResolveResult
	err = c.executeRequest("GET", url, &dnsResolved)

	return &dnsResolved, err
}

func (c *Client) GetDnsReverse(ips []string) (*DnsReverseResult, error) {
	params := map[string]string{"ips": strings.Join(ips, ",")}
	url, err := c.buildUrl(reversePath, params)
	if err != nil {
		return nil, err
	}

	var dnsReversed DnsReverseResult
	err = c.executeRequest("GET", url, &dnsReversed)

	return &dnsReversed, err
}
