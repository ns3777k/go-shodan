package shodan

import (
	"strings"
)

type DNSResolved map[string]string
type DNSReversed map[string][]string

const (
	resolvePath = "/dns/resolve"
	reversePath = "/dns/reverse"
)

// GetDNSResolve looks up the IP address for the provided list of hostnames
func (c *Client) GetDNSResolve(hostnames []string) (DNSResolved, error) {
	url, err := c.buildBaseURL(resolvePath, struct {
		Hostnames string `url:"hostnames"`
	}{strings.Join(hostnames, ",")})
	if err != nil {
		return nil, err
	}

	dnsResolved := make(DNSResolved)
	err = c.executeRequest("GET", url, &dnsResolved)

	return dnsResolved, err
}

// GetDNSReverse looks up the hostnames that have been defined for the given list of IP addresses
func (c *Client) GetDNSReverse(ip []string) (DNSReversed, error) {
	url, err := c.buildBaseURL(reversePath, struct {
		IP string `url:"ips"`
	}{strings.Join(ip, ",")})
	if err != nil {
		return nil, err
	}

	dnsReversed := make(DNSReversed)
	err = c.executeRequest("GET", url, &dnsReversed)

	return dnsReversed, err
}
