package shodan

const (
	hostPath             = "/shodan/host"
	hostCountPath        = "/shodan/host/count"
	hostSearchPath       = "/shodan/host/search"
	hostSearchTokensPath = "/shodan/host/search/tokens"
)

type HostServicesOptions struct {
	History bool `url:"history,omitempty"`
	Minify  bool `url:"minify,omitempty"`
}

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

type HostData struct {
	Product      string                 `json:"product"`
	Hostnames    []string               `json:"hostnames"`
	Version      string                 `json:"version"`
	Title        string                 `json:"title"`
	IPLong       int                    `json:"ip"`
	IP           string                 `json:"ip_str"`
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

type Host struct {
	OS              string      `json:"os"`
	Ports           []int       `json:"ports"`
	IPLong          int         `json:"ip"`
	IP              string      `json:"ip_str"`
	ISP             string      `json:"isp"`
	Hostnames       []string    `json:"hostnames"`
	Organization    string      `json:"org"`
	Vulnerabilities []string    `json:"vulns"`
	ASN             string      `json:"asn"`
	LastUpdate      string      `json:"last_update"`
	Data            []*HostData `json:"data"`
	HostLocation
}

type HostQueryOptions struct {
	Query  string `url:"query"`
	Facets string `url:"facets,omitempty"`
	Minify bool   `url:"minify,omitempty"`
	Page   int    `url:"page,omitempty"`
}

type HostMatch struct {
	Total   int                 `json:"total"`
	Facets  map[string][]*Facet `json:"facets"`
	Matches []*HostData         `json:"matches"`
}

type HostQueryTokens struct {
	Filters []string `json:"filters"`
	String  string   `json:"string"`
	Errors  []string `json:"errors"`
	// FIXME: should it really be interface{} ?
	Attributes map[string]interface{} `json:"attributes"`
}

// GetServicesForHost returns all services that have been found on the given host IP
func (c *Client) GetServicesForHost(ip string, options *HostServicesOptions) (*Host, error) {
	url, err := c.buildBaseURL(hostPath+"/"+ip, options)
	if err != nil {
		return nil, err
	}

	var host Host
	err = c.executeRequest("GET", url, &host, nil)

	return &host, err
}

// GetServicesCountForHost behaves identical to "/shodan/host/search" with the only difference that this method
// does not return any host results, it only returns the total number of results that matched the query and any facet
// information that was requested. As a result this method does not consume query credits
func (c *Client) GetHostsCountForQuery(options *HostQueryOptions) (*HostMatch, error) {
	url, err := c.buildBaseURL(hostCountPath, options)
	if err != nil {
		return nil, err
	}

	var found HostMatch
	err = c.executeRequest("GET", url, &found, nil)

	return &found, err
}

// GetHostsForQuery searches Shodan using the same query syntax as the website and use facets to get summary
// information for different properties. This method may use API query credits depending on usage. If any of the
// following criteria are met, your account will be deducated 1 query credit:
// 1. The search query contains a filter
// 2. Accessing results past the 1st page using the "page". For every 100 results past the 1st page 1 query credit is
// deducted
func (c *Client) GetHostsForQuery(options *HostQueryOptions) (*HostMatch, error) {
	url, err := c.buildBaseURL(hostSearchPath, options)
	if err != nil {
		return nil, err
	}

	var found HostMatch
	err = c.executeRequest("GET", url, &found, nil)

	return &found, err
}

// This method lets you determine which filters are being used by the query string and what parameters were provided
// to the filters.
func (c *Client) BreakQueryIntoTokens(query string) (*HostQueryTokens, error) {
	url, err := c.buildBaseURL(hostSearchTokensPath, struct {
		Query string `url:"query"`
	}{Query: query})
	if err != nil {
		return nil, err
	}

	var tokens HostQueryTokens
	err = c.executeRequest("GET", url, &tokens, nil)

	return &tokens, err
}
