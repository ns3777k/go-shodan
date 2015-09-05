package shodan

import (
	"strings"
)

type DnsResolved map[string]string
type DnsReversed map[string][]string

type dnsResolveOptions struct {
	Hostnames string `url:"hostnames"`
}

type dnsReverseOptions struct {
	Ips string `url:"ips"`
}

const (
	resolvePath = "/dns/resolve"
	reversePath = "/dns/reverse"
)

func (c *Client) GetDnsResolve(hostnames []string) (*DnsResolved, error) {
	options := &dnsResolveOptions{
		Hostnames: strings.Join(hostnames, ","),
	}
	url, err := c.buildUrl(resolvePath, options)
	if err != nil {
		return nil, err
	}

	var dnsResolved DnsResolved
	err = c.executeRequest("GET", url, &dnsResolved)

	return &dnsResolved, err
}

func (c *Client) GetDnsReverse(ips []string) (*DnsReversed, error) {
	options := &dnsReverseOptions{
		Ips: strings.Join(ips, ","),
	}
	url, err := c.buildUrl(reversePath, options)
	if err != nil {
		return nil, err
	}

	var dnsReversed DnsReversed
	err = c.executeRequest("GET", url, &dnsReversed)

	return &dnsReversed, err
}
