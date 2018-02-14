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
func (c *Client) GetDNSResolve(hostnames []string) (map[string]*net.IP, error) {
	url := c.buildBaseURL(resolvePath, struct {
		Hostnames string `url:"hostnames"`
	}{strings.Join(hostnames, ",")})

	dnsResolved := make(map[string]*net.IP)
	err := c.executeRequest("GET", url, &dnsResolved, nil)

	return dnsResolved, err
}

// GetDNSReverse looks up the hostnames that have been defined for the given list of IP addresses
func (c *Client) GetDNSReverse(ip []net.IP) (map[string]*[]string, error) {
	ips := make([]string, 0)

	for _, ipAddress := range ip {
		ips = append(ips, ipAddress.String())
	}

	url := c.buildBaseURL(reversePath, struct {
		IP string `url:"ips"`
	}{strings.Join(ips, ",")})

	dnsReversed := make(map[string]*[]string)
	err := c.executeRequest("GET", url, &dnsReversed, nil)

	return dnsReversed, err
}
