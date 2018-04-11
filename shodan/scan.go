package shodan

import (
	"context"
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
func (c *Client) Scan(ctx context.Context, ip []string) (*CrawlScanStatus, error) {
	var crawlScanStatus CrawlScanStatus

	body := neturl.Values{}
	body.Add("ips", strings.Join(ip, ","))

	req, err := c.NewRequest("POST", scanPath, nil, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &crawlScanStatus); err != nil {
		return nil, err
	}

	return &crawlScanStatus, nil
}

// ScanInternet requests Shodan to crawl the Internet for a specific port.
// This method is restricted to security researchers and companies with a Shodan Data license.
func (c *Client) ScanInternet(ctx context.Context, port int, protocol string) (string, error) {
	crawlScanInternetStatus := new(struct {
		ID string `json:"id"`
	})

	body := neturl.Values{}
	body.Add("port", strconv.Itoa(port))
	body.Add("protocol", protocol)

	req, err := c.NewRequest("POST", scanInternetPath, nil, strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}

	if err := c.Do(ctx, req, crawlScanInternetStatus); err != nil {
		return "", err
	}

	return crawlScanInternetStatus.ID, nil
}

// GetScanStatus checks the progress of a previously submitted scan request.
func (c *Client) GetScanStatus(ctx context.Context, id string) (*ScanStatus, error) {
	var scanStatus ScanStatus
	path := fmt.Sprintf(scanStatusPath, id)

	req, err := c.NewRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &scanStatus); err != nil {
		return nil, err
	}

	return &scanStatus, nil
}
