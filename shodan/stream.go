package shodan

import (
	"context"
	"fmt"
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
