package shodan

import (
	"net"
	"strings"
)

const (
	resolvePath = "/dns/resolve"
	reversePath = "/dns/reverse"
)

// GetDNSResolve looks up the IP address for the provided list of hostnames
func (c *Client) GetDNSResolve(hostnames []string) (map[string]*string, error) {
	url := c.buildBaseURL(resolvePath, struct {
		Hostnames string `url:"hostnames"`
	}{strings.Join(hostnames, ",")})

	dnsResolved := make(map[string]*string)
	err := c.executeRequest("GET", url, &dnsResolved, nil)

	return dnsResolved, err
}

// GetDNSReverse looks up the hostnames that have been defined for the given list of IP addresses
func (c *Client) GetDNSReverse(ip []string) (map[string]*[]string, error) {
	for _, ipAddress := range ip {
		if parsedIP := net.ParseIP(ipAddress); parsedIP == nil {
			return nil, &net.ParseError{
				Type: "IP address",
				Text: ipAddress,
			}
		}
	}

	url := c.buildBaseURL(reversePath, struct {
		IP string `url:"ips"`
	}{strings.Join(ip, ",")})

	dnsReversed := make(map[string]*[]string)
	err := c.executeRequest("GET", url, &dnsReversed, nil)

	return dnsReversed, err
}
