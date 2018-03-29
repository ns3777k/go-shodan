package shodan

import (
	"net"
	"encoding/json"
	"strconv"
	"context"
)

const (
	hostPath             = "/shodan/host"
	hostCountPath        = "/shodan/host/count"
	hostSearchPath       = "/shodan/host/search"
	hostSearchTokensPath = "/shodan/host/search/tokens"
)

//HostServicesOptions is options for querying services.
type HostServicesOptions struct {
	History bool `url:"history,omitempty"`
	Minify  bool `url:"minify,omitempty"`
}

// HostLocation is the location of the host.
type HostLocation struct {
	City         string  `json:"city"`
	RegionCode   string  `json:"region_code"`
	AreaCode     int     `json:"area_code"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Country      string  `json:"country_name"`
	CountryCode  string  `json:"country_code"`
	CountryCode3 string  `json:"country_code3"`
	Postal       string  `json:"postal_code"`
	DMA          int     `json:"dma_code"`
}

// HostVersion is string with custom unmarshaling.
type HostVersion string

// UnmarshalJSON handles either a string or a number
// and casts it to string.
func (v *HostVersion) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		*v = HostVersion(s)
		return nil
	}

	var n int
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}

	*v = HostVersion(strconv.Itoa(n))

	return nil
}

func (v *HostVersion) String() string {
	return string(*v)
}

// HostData is all services that have been found on the given host IP.
type HostData struct {
	Product      string                 `json:"product"`
	Hostnames    []string               `json:"hostnames"`
	Version      HostVersion            `json:"version"`
	Title        string                 `json:"title"`
	IP           net.IP                 `json:"ip_str"`
	OS           string                 `json:"os"`
	Organization string                 `json:"org"`
	ISP          string                 `json:"isp"`
	CPE          []string               `json:"cpe"`
	Data         string                 `json:"data"`
	ASN          string                 `json:"asn"`
	Port         int                    `json:"port"`
	HTML         string                 `json:"html"`
	Banner       string                 `json:"banner"`
	Link         string                 `json:"link"`
	Transport    string                 `json:"transport"`
	Domains      []string               `json:"domains"`
	Timestamp    string                 `json:"timestamp"`
	DeviceType   string                 `json:"devicetype"`
	Location     *HostLocation          `json:"location"`
	ShodanData   map[string]interface{} `json:"_shodan"`
	Opts         map[string]interface{} `json:"opts"`
}

// Host is the all information about the host.
type Host struct {
	OS              string      `json:"os"`
	Ports           []int       `json:"ports"`
	IP              net.IP      `json:"ip_str"`
	ISP             string      `json:"isp"`
	Hostnames       []string    `json:"hostnames"`
	Organization    string      `json:"org"`
	Vulnerabilities []string    `json:"vulns"`
	ASN             string      `json:"asn"`
	LastUpdate      string      `json:"last_update"`
	Data            []*HostData `json:"data"`
	HostLocation
}

// HostQueryOptions is Shodan search query options.
type HostQueryOptions struct {
	Query  string `url:"query"`
	Facets string `url:"facets,omitempty"`
	Minify bool   `url:"minify,omitempty"`
	Page   int    `url:"page,omitempty"`
}

// HostMatch is the search results with all matched hosts.
type HostMatch struct {
	Total   int                 `json:"total"`
	Facets  map[string][]*Facet `json:"facets"`
	Matches []*HostData         `json:"matches"`
}

// HostQueryTokens is filters are being used by the query string and what
// parameters were provided to the filters.
type HostQueryTokens struct {
	Filters []string `json:"filters"`
	String  string   `json:"string"`
	Errors  []string `json:"errors"`
	Attributes map[string]interface{} `json:"attributes"`
}

// GetServicesForHost returns all services that have been found on the given host IP
func (c *Client) GetServicesForHost(ctx context.Context, ip string, options *HostServicesOptions) (*Host, error) {
	var host Host

	req, err := c.NewRequest("GET", hostPath+"/"+ip, options, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &host); err != nil {
		return nil, err
	}

	return &host, nil
}

// GetHostsCountForQuery behaves identical to "/shodan/host/search" with the only difference that this method
// does not return any host results, it only returns the total number of results that matched the query and any facet
// information that was requested. As a result this method does not consume query credits
func (c *Client) GetHostsCountForQuery(ctx context.Context, options *HostQueryOptions) (*HostMatch, error) {
	var found HostMatch

	req, err := c.NewRequest("GET", hostCountPath, options, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &found); err != nil {
		return nil, err
	}

	return &found, nil
}

// GetHostsForQuery searches Shodan using the same query syntax as the website and use facets to get summary
// information for different properties. This method may use API query credits depending on usage. If any of the
// following criteria are met, your account will be deducated 1 query credit:
// 1. The search query contains a filter
// 2. Accessing results past the 1st page using the "page". For every 100 results past the 1st page 1 query credit is
// deducted
func (c *Client) GetHostsForQuery(ctx context.Context, options *HostQueryOptions) (*HostMatch, error) {
	var found HostMatch

	req, err := c.NewRequest("GET", hostSearchPath, options, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &found); err != nil {
		return nil, err
	}

	return &found, nil
}

// BreakQueryIntoTokens determines which filters are being used by the query string
// and what parameters were provided to the filters.
func (c *Client) BreakQueryIntoTokens(ctx context.Context, query string) (*HostQueryTokens, error) {
	var tokens HostQueryTokens

	options := struct {
		Query string `url:"query"`
	}{Query: query}
	req, err := c.NewRequest("GET", hostSearchTokensPath, options, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &tokens); err != nil {
		return nil, err
	}

	return &tokens, nil
}
