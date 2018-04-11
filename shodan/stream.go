package shodan

import (
	"strconv"
	"fmt"
	"strings"
	"context"
	"net/http"
	"github.com/moul/http2curl"
	"log"
	"bufio"
	"bytes"
)

const (
	bannersPath        = "/shodan/banners"
	bannersAlertPath   = "/shodan/alert/%s"
	bannersAlertsPath  = "/shodan/alert"
	bannersPortsPath   = "/shodan/ports/%s"
	bannersCountryPath = "/shodan/countries/%s"
	bannersASNPath     = "/shodan/asn/%s"
)

func (c *Client) DoStream(ctx context.Context, req *http.Request) (*http.Response, error) {
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	if c.Debug {
		if command, err := http2curl.GetCurlCommand(req); err == nil {
			log.Printf("shodan client request: %s\n", command)
		}
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, getErrorFromResponse(resp)
	}

	return resp, nil
}

func (c *Client) handleResponseStream(resp *http.Response, ch chan *HostData) {
	reader := bufio.NewReader(resp.Body)

	for {
		banner := new(HostData)

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

		if err := c.parseResponse(banner, bytes.NewBuffer(chunk)); err != nil {
			resp.Body.Close()
			close(ch)
			break
		}

		ch <- banner
	}
}

// GetBannersByASN provides a filtered, bandwidth-saving view of the Banners stream in case
// you are only interested in devices located in certain ASNs.
func (c *Client) GetBannersByASN(ctx context.Context, asn []string, ch chan *HostData) error {
	path := fmt.Sprintf(bannersASNPath, strings.Join(asn, ","))

	req, err := c.NewStreamingRequest("GET", path, nil, nil)
	if err != nil {
		return err
	}

	resp, err := c.DoStream(ctx, req)
	if err != nil {
		return err
	}

	go c.handleResponseStream(resp, ch)

	return nil
}

// GetBannersByCountries provides a filtered, bandwidth-saving view of the Banners
// stream in case you are only interested in devices located in certain countries.
func (c *Client) GetBannersByCountries(ctx context.Context, countries []string, ch chan *HostData) error {
	strCountries := make([]string, 0)
	for _, country := range countries {
		strCountries = append(strCountries, strings.ToUpper(country))
	}

	path := fmt.Sprintf(bannersCountryPath, strings.Join(strCountries, ","))

	req, err := c.NewStreamingRequest("GET", path, nil, nil)
	if err != nil {
		return err
	}

	resp, err := c.DoStream(ctx, req)
	if err != nil {
		return err
	}

	go c.handleResponseStream(resp, ch)

	return nil
}

// GetBannersByPorts returns only banner data for the list of specified hosts.
// This stream provides a filtered, bandwidth-saving view of the Banners stream
// in case you are only interested in a specific list of ports.
func (c *Client) GetBannersByPorts(ctx context.Context, ports []int, ch chan *HostData) error {
	strPorts := make([]string, 0)
	for _, port := range ports {
		strPorts = append(strPorts, strconv.Itoa(port))
	}

	path := fmt.Sprintf(bannersPortsPath, strings.Join(strPorts, ","))

	req, err := c.NewStreamingRequest("GET", path, nil, nil)
	if err != nil {
		return err
	}

	resp, err := c.DoStream(ctx, req)
	if err != nil {
		return err
	}

	go c.handleResponseStream(resp, ch)

	return nil
}

// GetBannersByAlert subscribes to banners discovered on the IP range defined
// in a specific network alert.
func (c *Client) GetBannersByAlert(ctx context.Context, id string, ch chan *HostData) error {
	path := fmt.Sprintf(bannersAlertPath, id)

	req, err := c.NewStreamingRequest("GET", path, nil, nil)
	if err != nil {
		return err
	}

	resp, err := c.DoStream(ctx, req)
	if err != nil {
		return err
	}

	go c.handleResponseStream(resp, ch)

	return nil
}

// GetBannersByAlerts subscribes to banners discovered on all IP ranges described
// in the network alerts.
func (c *Client) GetBannersByAlerts(ctx context.Context, ch chan *HostData) error {
	req, err := c.NewStreamingRequest("GET", bannersAlertsPath, nil, nil)
	if err != nil {
		return err
	}

	resp, err := c.DoStream(ctx, req)
	if err != nil {
		return err
	}

	go c.handleResponseStream(resp, ch)

	return nil
}

// GetBanners provides ALL of the data that Shodan collects. Use this stream
// if you need access to everything and / or want to store your own Shodan database
// locally. If you only care about specific ports, please use the Ports stream.
func (c *Client) GetBanners(ctx context.Context, ch chan *HostData) error {
	req, err := c.NewStreamingRequest("GET", bannersPath, nil, nil)
	if err != nil {
		return err
	}

	resp, err := c.DoStream(ctx, req)
	if err != nil {
		return err
	}

	go c.handleResponseStream(resp, ch)

	return nil
}
