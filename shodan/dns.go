package shodan

import (
	"strings"
)

type DNSResolved map[string]string
type DNSReversed map[string][]string

type dnsResolveOptions struct {
	Hostnames string `url:"hostnames"`
}

type dnsReverseOptions struct {
	IP string `url:"ips"`
}

const (
	resolvePath = "/dns/resolve"
	reversePath = "/dns/reverse"
)

// GetDNSResolve looks up the IP address for the provided list of hostnames
func (c *Client) GetDNSResolve(hostnames []string) (DNSResolved, error) {
	options := &dnsResolveOptions{
		Hostnames: strings.Join(hostnames, ","),
	}
	url, err := c.buildURL(resolvePath, options)
	if err != nil {
		return nil, err
	}

	dnsResolved := make(DNSResolved)
	err = c.executeRequest("GET", url, &dnsResolved)

	return dnsResolved, err
}

// GetDNSReverse looks up the hostnames that have been defined for the given list of IP addresses
func (c *Client) GetDNSReverse(ip []string) (DNSReversed, error) {
	options := &dnsReverseOptions{
		IP: strings.Join(ip, ","),
	}
	url, err := c.buildURL(reversePath, options)
	if err != nil {
		return nil, err
	}

	dnsReversed := make(DNSReversed)
	err = c.executeRequest("GET", url, &dnsReversed)

	return dnsReversed, err
}
