package shodan

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

const (
	bannersPath       = "/shodan/banners"
	bannersAlertPath  = "/shodan/alert/%s"
	bannersAlertsPath = "/shodan/alert"
	bannersPortsPath  = "/shodan/ports/%s"
)

func (c *Client) readBannersResponse(rawChan chan []byte) {
	for {
		var banner HostData
		res, ok := <-rawChan

		if !ok {
			close(c.StreamChan)
			break
		}

		buf := bytes.NewBuffer(res)
		if err := c.parseResponse(&banner, buf); err != nil {
			close(c.StreamChan)
			break
		}

		c.StreamChan <- banner
	}
}

func (c *Client) beginStreaming(path string) {
	url := c.buildStreamBaseURL(path, nil)
	rawChan := make(chan []byte)

	go c.readBannersResponse(rawChan)
	go c.executeStreamRequest("GET", url, rawChan)
}

// GetBannersByPorts returns only banner data for the list of specified hosts.
// This stream provides a filtered, bandwidth-saving view of the Banners stream
// in case you are only interested in a specific list of ports.
func (c *Client) GetBannersByPorts(ports []int) {
	stringifiedPorts := make([]string, 0)
	for _, port := range ports {
		stringifiedPorts = append(stringifiedPorts, strconv.Itoa(port))
	}

	path := fmt.Sprintf(bannersPortsPath, strings.Join(stringifiedPorts, ","))
	c.beginStreaming(path)
}

// GetBannersByAlert subscribes to banners discovered on the IP range defined
// in a specific network alert.
func (c *Client) GetBannersByAlert(id string) {
	path := fmt.Sprintf(bannersAlertPath, id)
	c.beginStreaming(path)
}

// GetBannersByAlerts subscribes to banners discovered on all IP ranges described
// in the network alerts.
func (c *Client) GetBannersByAlerts() {
	c.beginStreaming(bannersAlertsPath)
}

// GetBanners provides ALL of the data that Shodan collects. Use this stream
// if you need access to everything and / or want to store your own Shodan database
// locally. If you only care about specific ports, please use the Ports stream.
func (c *Client) GetBanners() {
	c.beginStreaming(bannersPath)
}
