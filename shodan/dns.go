package shodan

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	dnsPath     = "/dns/domain/%s"
	resolvePath = "/dns/resolve"
	reversePath = "/dns/reverse"
)

type DomainDNSInfo struct {
	Domain     string              `json:"domain"`
	Tags       []string            `json:"tags"`
	Data       []*SubdomainDNSInfo `json:"data"`
	Subdomains []string            `json:"subdomains"`
}

type SubdomainDNSInfo struct {
	Subdomain string    `json:"subdomain"`
	Type      string    `json:"type"`
	Value     string    `json:"value"`
	LastSeen  time.Time `json:"last_seen"`
}

// GetDomain returns all the subdomains and other DNS entries for the given domain. Uses 1 query credit per lookup.
func (c *Client) GetDomain(ctx context.Context, domain string) (*DomainDNSInfo, error) {
	var info DomainDNSInfo
	path := fmt.Sprintf(dnsPath, domain)
	req, err := c.NewRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

// GetDNSResolve looks up the IP address for the provided list of hostnames
func (c *Client) GetDNSResolve(ctx context.Context, hostnames []string) (map[string]*net.IP, error) {
	dnsResolved := make(map[string]*net.IP)
	params := struct {
		Hostnames string `url:"hostnames"`
	}{strings.Join(hostnames, ",")}

	req, err := c.NewRequest("GET", resolvePath, params, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &dnsResolved); err != nil {
		return nil, err
	}

	return dnsResolved, nil
}

// GetDNSReverse looks up the hostnames that have been defined for the given list of IP addresses
func (c *Client) GetDNSReverse(ctx context.Context, ip []net.IP) (map[string]*[]string, error) {
	ips := make([]string, 0)

	for _, ipAddress := range ip {
		ips = append(ips, ipAddress.String())
	}

	dnsReversed := make(map[string]*[]string)
	params := struct {
		IP string `url:"ips"`
	}{strings.Join(ips, ",")}

	req, err := c.NewRequest("GET", reversePath, params, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &dnsReversed); err != nil {
		return nil, err
	}

	return dnsReversed, nil
}
