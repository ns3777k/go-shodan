package shodan

import (
	"fmt"
	neturl "net/url"
	"strconv"
	"strings"
)

// ScanStatusState is an alias to string that represents a scan state.
type ScanStatusState string

const (
	scanStatusPath   = "/shodan/scan/%s"
	scanPath         = "/shodan/scan"
	scanInternetPath = "/shodan/scan/internet"

	// ScanStatusSubmitting is "SUBMITTING"
	ScanStatusSubmitting ScanStatusState = "SUBMITTING"

	// ScanStatusQueue is "QUEUE"
	ScanStatusQueue ScanStatusState = "QUEUE"

	// ScanStatusProcessing is "PROCESSING"
	ScanStatusProcessing ScanStatusState = "PROCESSING"

	// ScanStatusDone is "DONE"
	ScanStatusDone ScanStatusState = "DONE"
)

// ScanStatus is a current scan status.
type ScanStatus struct {
	ID     string          `json:"id"`
	Count  int             `json:"count"`
	Status ScanStatusState `json:"status"`
}

// CrawlScanStatus is the response to a scan request.
type CrawlScanStatus struct {
	ID          string `json:"id"`
	Count       int    `json:"count"`
	CreditsLeft int    `json:"credits_left"`
}

// Scan requests Shodan to crawl a network.
// This method uses API scan credits: 1 IP consumes 1 scan credit. You must have a paid API plan (either one-time
// payment or subscription) in order to use this method.
func (c *Client) Scan(ip []string) (*CrawlScanStatus, error) {
	url := c.buildBaseURL(scanPath, nil)

	var crawlScanStatus CrawlScanStatus
	body := neturl.Values{}
	body.Add("ips", strings.Join(ip, ","))

	err := c.executeRequest("POST", url, &crawlScanStatus, strings.NewReader(body.Encode()))

	return &crawlScanStatus, err
}

// ScanInternet requests Shodan to crawl the Internet for a specific port.
// This method is restricted to security researchers and companies with a Shodan Data license. To apply for access to
// this method as a researcher, please email jmath@shodan.io with information about your project. Access is restricted
// to prevent abuse.
func (c *Client) ScanInternet(port int, protocol string) (string, error) {
	url := c.buildBaseURL(scanInternetPath, nil)

	crawlScanInternetStatus := new(struct {
		ID string `json:"id"`
	})

	body := neturl.Values{}
	body.Add("port", strconv.Itoa(port))
	body.Add("protocol", protocol)

	err := c.executeRequest("POST", url, crawlScanInternetStatus, strings.NewReader(body.Encode()))

	return crawlScanInternetStatus.ID, err
}

// GetScanStatus checks the progress of a previously submitted scan request.
func (c *Client) GetScanStatus(id string) (*ScanStatus, error) {
	path := fmt.Sprintf(scanStatusPath, id)
	url := c.buildBaseURL(path, nil)

	var scanStatus ScanStatus
	err := c.executeRequest("GET", url, &scanStatus, nil)

	return &scanStatus, err
}
