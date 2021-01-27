package shodan

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	bannersPath        = "/shodan/banners"
	bannersAlertPath   = "/shodan/alert/%s"
	bannersAlertsPath  = "/shodan/alert"
	bannersPortsPath   = "/shodan/ports/%s"
	bannersCountryPath = "/shodan/countries/%s"
	bannersASNPath     = "/shodan/asn/%s"
)

type StreamAlertMeta struct {
	ID            string
	Name          string
	Trigger       string
	Sha1Signature string
}

type StreamHostData struct {
	*HostData
	AlertMeta *StreamAlertMeta
}

// handleResponseStream reads response body, transforms it to *HostData and
// sends to channel.
func (c *Client) handleResponseStream(resp *http.Response, ch chan *StreamHostData) {
	reader := bufio.NewReader(resp.Body)

	for {
		chunk, err := reader.ReadBytes('\n')
		if err != nil {
			resp.Body.Close()
			close(ch)
			break
		}

		chunk = bytes.TrimRight(chunk, "\n\r")
		if len(chunk) == 0 {
			continue
		}

		banner := new(StreamHostData)
		if err := c.parseResponse(banner, bytes.NewBuffer(chunk)); err != nil {
			resp.Body.Close()
			close(ch)
			break
		}

		banner.AlertMeta = &StreamAlertMeta{
			ID:            resp.Header.Get("SHODAN-ALERT-ID"),
			Name:          resp.Header.Get("SHODAN-ALERT-NAME"),
			Trigger:       resp.Header.Get("SHODAN-ALERT-TRIGGER"),
			Sha1Signature: resp.Header.Get("SHODAN-SIGNATURE-SHA1"),
		}

		ch <- banner
	}
}

// NewStreamingRequest prepares new request to streaming api.
func (c *Client) NewStreamingRequest(path string, params interface{}) (*http.Request, error) {
	u, err := url.Parse(c.StreamBaseURL + path)
	if err != nil {
		return nil, err
	}

	return c.newRequest("GET", u, params, nil)
}

// DoStream executes streaming request.
func (c *Client) DoStream(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, getErrorFromResponse(resp)
	}

	return resp, nil
}

// startStreaming creates new streaming request and sends it.
func (c *Client) startStreaming(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewStreamingRequest(path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.DoStream(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetBannersByASN provides a filtered, bandwidth-saving view of the Banners stream in case
// you are only interested in devices located in certain ASNs.
func (c *Client) GetBannersByASN(ctx context.Context, asn []string, ch chan *StreamHostData) error {
	path := fmt.Sprintf(bannersASNPath, strings.Join(asn, ","))
	resp, err := c.startStreaming(ctx, path)
	if err != nil {
		return err
	}

	go c.handleResponseStream(resp, ch)

	return nil
}

// GetBannersByCountries provides a filtered, bandwidth-saving view of the Banners
// stream in case you are only interested in devices located in certain countries.
func (c *Client) GetBannersByCountries(ctx context.Context, countries []string, ch chan *StreamHostData) error {
	strCountries := make([]string, 0)
	for _, country := range countries {
		strCountries = append(strCountries, strings.ToUpper(country))
	}

	path := fmt.Sprintf(bannersCountryPath, strings.Join(strCountries, ","))
	resp, err := c.startStreaming(ctx, path)
	if err != nil {
		return err
	}

	go c.handleResponseStream(resp, ch)

	return nil
}

// GetBannersByPorts returns only banner data for the list of specified hosts.
// This stream provides a filtered, bandwidth-saving view of the Banners stream
// in case you are only interested in a specific list of ports.
func (c *Client) GetBannersByPorts(ctx context.Context, ports []int, ch chan *StreamHostData) error {
	strPorts := make([]string, 0)
	for _, port := range ports {
		strPorts = append(strPorts, strconv.Itoa(port))
	}

	path := fmt.Sprintf(bannersPortsPath, strings.Join(strPorts, ","))
	resp, err := c.startStreaming(ctx, path)
	if err != nil {
		return err
	}

	go c.handleResponseStream(resp, ch)

	return nil
}

// GetBannersByAlert subscribes to banners discovered on the IP range defined
// in a specific network alert.
func (c *Client) GetBannersByAlert(ctx context.Context, id string, ch chan *StreamHostData) error {
	path := fmt.Sprintf(bannersAlertPath, id)
	resp, err := c.startStreaming(ctx, path)
	if err != nil {
		return err
	}

	go c.handleResponseStream(resp, ch)

	return nil
}

// GetBannersByAlerts subscribes to banners discovered on all IP ranges described
// in the network alerts.
func (c *Client) GetBannersByAlerts(ctx context.Context, ch chan *StreamHostData) error {
	resp, err := c.startStreaming(ctx, bannersAlertsPath)
	if err != nil {
		return err
	}

	go c.handleResponseStream(resp, ch)

	return nil
}

// GetBanners provides ALL of the data that Shodan collects. Use this stream
// if you need access to everything and / or want to store your own Shodan database
// locally. If you only care about specific ports, please use the Ports stream.
func (c *Client) GetBanners(ctx context.Context, ch chan *StreamHostData) error {
	resp, err := c.startStreaming(ctx, bannersPath)
	if err != nil {
		return err
	}

	go c.handleResponseStream(resp, ch)

	return nil
}
