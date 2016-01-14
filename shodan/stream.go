package shodan

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

const (
	bannersPath      = "/shodan/banners"
	bannersPortsPath = "/shodan/ports/%s"
)

func (c *Client) readBannersResponse(rawChan chan []byte) {
	for {
		var banner HostData
		res := <-rawChan

		if res == nil {
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

// GetBannersByPorts returns only banner data for the list of specified hosts. This stream provides a filtered,
// bandwidth-saving view of the Banners stream in case you are only interested in a specific list of ports.
func (c *Client) GetBannersByPorts(ports []int) error {
	stringifiedPorts := make([]string, len(ports))
	for _, port := range ports {
		stringifiedPorts = append(stringifiedPorts, strconv.Itoa(port))
	}

	path := fmt.Sprintf(bannersPortsPath, strings.Join(stringifiedPorts, ","))
	url, err := c.buildStreamBaseURL(path, nil)
	if err != nil {
		return err
	}

	rawChan := make(chan []byte)
	go c.readBannersResponse(rawChan)
	go c.executeStreamRequest("GET", url, rawChan)

	return nil
}

// GetBanners provides ALL of the data that Shodan collects. Use this stream if you need access to everything and / or
// want to store your own Shodan database locally. If you only care about specific ports, please use the Ports stream.
func (c *Client) GetBanners() error {
	url, err := c.buildStreamBaseURL(bannersPath, nil)
	if err != nil {
		return err
	}

	rawChan := make(chan []byte)
	go c.readBannersResponse(rawChan)
	go c.executeStreamRequest("GET", url, rawChan)

	return nil
}
