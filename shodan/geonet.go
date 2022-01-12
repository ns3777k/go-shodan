package shodan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

const (
	geonetPingPath  = "/api/ping/%s"
	geonetPingsPath = "/api/geoping/%s"
	geonetDNSPath   = "/api/dns/%s"
	geonetDNSsPath  = "/api/geodns/%s"
)

// Location from where the request was made.
type Location struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Coords  string `json:"latlon"`
}

// PingResult contains details about ping request.
type PingResult struct {
	IP              string    `json:"ip"`
	IsAlive         bool      `json:"is_alive"`
	MinRTT          float64   `json:"min_rtt"`
	MaxRTT          float64   `json:"max_rtt"`
	AvgRTT          float64   `json:"avg_rtt"`
	RTTs            []float64 `json:"rtts"`
	PacketsSent     int       `json:"packets_sent"`
	PacketsReceived int       `json:"packets_received"`
	PacketLoss      float64   `json:"packet_loss"`
	FromLocation    *Location `json:"from_loc"`
	Error           string    `json:"error"`
}

// DNSQueryResult is a response to the dns query.
type DNSQueryResult struct {
	Answers      []*DNSQueryResultAnswers `json:"answers"`
	FromLocation *Location                `json:"from_loc"`
	Error        string                   `json:"error"`
}

// DNSQueryResultAnswers contains dns record type and its value.
type DNSQueryResultAnswers struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// DNSQueryOptions customizes dns query.
type DNSQueryOptions struct {
	RecordType string `url:"rtype"`
}

func getGeoNetErrorFromResponse(r *http.Response) error {
	errorResponse := new(struct {
		Error string `json:"detail"`
	})

	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ErrBodyRead
	}

	if err := json.Unmarshal(message, errorResponse); err != nil {
		return err
	}

	return errors.New(errorResponse.Error)
}

// NewGeoNetRequest prepares new request to geonet shodan api.
func (c *Client) NewGeoNetRequest(path string, params interface{}) (*http.Request, error) {
	u, err := url.Parse(c.GeoNetBaseURL + path)
	if err != nil {
		return nil, err
	}

	qs, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	u.RawQuery = qs.Encode()

	return http.NewRequest("GET", u.String(), nil)
}

// GeoPing executes a ping command from one location.
func (c *Client) GeoPing(ctx context.Context, ip net.IP) (*PingResult, error) {
	var pingResult PingResult

	req, err := c.NewGeoNetRequest(fmt.Sprintf(geonetPingPath, ip.String()), nil)
	if err != nil {
		return nil, err
	}

	if err := c.DoWithErrorHandling(ctx, req, &pingResult, getGeoNetErrorFromResponse); err != nil {
		return nil, err
	}

	return &pingResult, nil
}

// GeoPings executes a ping command from multiple locations.
func (c *Client) GeoPings(ctx context.Context, ip net.IP) ([]*PingResult, error) {
	pingResults := make([]*PingResult, 0)

	req, err := c.NewGeoNetRequest(fmt.Sprintf(geonetPingsPath, ip.String()), nil)
	if err != nil {
		return nil, err
	}

	if err := c.DoWithErrorHandling(ctx, req, &pingResults, getGeoNetErrorFromResponse); err != nil {
		return nil, err
	}

	return pingResults, nil
}

// GeoDNSQuery queries a record for a hostname from one location.
func (c *Client) GeoDNSQuery(ctx context.Context, hostname string, options *DNSQueryOptions) (*DNSQueryResult, error) {
	var queryResult DNSQueryResult

	req, err := c.NewGeoNetRequest(fmt.Sprintf(geonetDNSPath, hostname), options)
	if err != nil {
		return nil, err
	}

	if err := c.DoWithErrorHandling(ctx, req, &queryResult, getGeoNetErrorFromResponse); err != nil {
		return nil, err
	}

	return &queryResult, nil
}

// GeoDNSQueries queries a record for a hostname from multiple locations.
func (c *Client) GeoDNSQueries(
	ctx context.Context,
	hostname string,
	options *DNSQueryOptions,
) ([]*DNSQueryResult, error) {
	queryResult := make([]*DNSQueryResult, 0)

	req, err := c.NewGeoNetRequest(fmt.Sprintf(geonetDNSsPath, hostname), options)
	if err != nil {
		return nil, err
	}

	if err := c.DoWithErrorHandling(ctx, req, &queryResult, getGeoNetErrorFromResponse); err != nil {
		return nil, err
	}

	return queryResult, nil
}
