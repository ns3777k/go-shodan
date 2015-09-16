package shodan

import (
	"strings"
	neturl "net/url"
	"strconv"
)

const (
	scanPath = "/shodan/scan"
	scanInternetPath = "/shodan/scan/internet"
)

type CrawlScanStatus struct {
	ID          string `json:"id"`
	Count       int    `json:"count"`
	CreditsLeft int    `json:"credits_left"`
}

// Scan requests Shodan to crawl a network.
// This method uses API scan credits: 1 IP consumes 1 scan credit. You must have a paid API plan (either one-time
// payment or subscription) in order to use this method.
func (c *Client) Scan(ip []string) (*CrawlScanStatus, error) {
	url, err := c.buildBaseURL(scanPath, nil)
	if err != nil {
		return nil, err
	}

	var crawlScanStatus CrawlScanStatus
	body := neturl.Values{}
	body.Add("ips", strings.Join(ip, ","))

	err = c.executeRequest("POST", url, &crawlScanStatus, strings.NewReader(body.Encode()))
	return &crawlScanStatus, err
}

// ScanInternet requests Shodan to crawl the Internet for a specific port.
// This method is restricted to security researchers and companies with a Shodan Data license. To apply for access to
// this method as a researcher, please email jmath@shodan.io with information about your project. Access is restricted
// to prevent abuse.
func (c *Client) ScanInternet(port int, protocol string) (string, error) {
	url, err := c.buildBaseURL(scanInternetPath, nil)
	if err != nil {
		return "", err
	}

	crawlScanInternetStatus := new(struct {
		ID string `json:"id"`
	})

	body := neturl.Values{}
	body.Add("port", strconv.Itoa(port))
	body.Add("protocol", protocol)

	err = c.executeRequest("POST", url, crawlScanInternetStatus, strings.NewReader(body.Encode()))
	return crawlScanInternetStatus.ID, err
}
