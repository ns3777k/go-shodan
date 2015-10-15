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

func (c *Client) readBannersResponse(ch chan HostData, internalChan chan []byte) {
	for {
		var banner HostData
		res := <-internalChan

		if res == nil {
			close(ch)
			break
		}

		buf := bytes.NewBuffer(res)
		if err := c.parseResponse(&banner, buf); err != nil {
			close(ch)
			break
		}

		ch <- banner
	}
}

// GetBannersByPorts returns only banner data for the list of specified hosts. This stream provides a filtered,
// bandwidth-saving view of the Banners stream in case you are only interested in a specific list of ports.
func (c *Client) GetBannersByPorts(ports []int, ch chan HostData) error {
	stringifiedPorts := make([]string, len(ports))
	for _, port := range ports {
		stringifiedPorts = append(stringifiedPorts, strconv.Itoa(port))
	}

	path := fmt.Sprintf(bannersPortsPath, strings.Join(stringifiedPorts, ","))
	url, err := c.buildStreamBaseURL(path, nil)
	if err != nil {
		return err
	}

	internalChan := make(chan []byte)
	go c.readBannersResponse(ch, internalChan)
	go c.executeStreamRequest("GET", url, internalChan)

	return nil
}

// GetBanners provides ALL of the data that Shodan collects. Use this stream if you need access to everything and / or
// want to store your own Shodan database locally. If you only care about specific ports, please use the Ports stream.
func (c *Client) GetBanners(ch chan HostData) error {
	url, err := c.buildStreamBaseURL(bannersPath, nil)
	if err != nil {
		return err
	}

	internalChan := make(chan []byte)
	go c.readBannersResponse(ch, internalChan)
	go c.executeStreamRequest("GET", url, internalChan)

	return nil
}
