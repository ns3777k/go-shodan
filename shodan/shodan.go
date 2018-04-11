package shodan

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"bufio"
	"bytes"
	"context"
	"github.com/google/go-querystring/query"
	"github.com/moul/http2curl"
	"log"
	"os"
)

const (
	baseURL        = "https://api.shodan.io"
	exploitBaseURL = "https://exploits.shodan.io/api"
	streamBaseURL  = "https://stream.shodan.io"
)

func getErrorFromResponse(r *http.Response) error {
	errorResponse := new(struct {
		Error string `json:"error"`
	})
	message, err := ioutil.ReadAll(r.Body)
	if err == nil {
		if err := json.Unmarshal(message, errorResponse); err == nil {
			return errors.New(errorResponse.Error)
		}

		return errors.New(strings.TrimSpace(string(message)))
	}

	return ErrBodyRead
}

// Client represents Shodan HTTP client
type Client struct {
	Token          string
	BaseURL        string
	ExploitBaseURL string
	StreamBaseURL  string
	StreamChan     chan HostData
	Debug          bool

	Client *http.Client
}

// NewClient creates new Shodan client
func NewClient(client *http.Client, token string) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	return &Client{
		Token:          token,
		BaseURL:        baseURL,
		ExploitBaseURL: exploitBaseURL,
		StreamBaseURL:  streamBaseURL,
		Client:         client,
	}
}

// NewEnvClient creates new Shodan client using environment variable
// SHODAN_KEY as the token.
func NewEnvClient(client *http.Client) *Client {
	return NewClient(client, os.Getenv("SHODAN_KEY"))
}

// SetDebug toggles the debug mode.
func (c *Client) SetDebug(debug bool) {
	c.Debug = debug
}

// NewExploitRequest prepares new request to exploit shodan api.
func (c *Client) NewExploitRequest(method string, path string, params interface{}, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(c.ExploitBaseURL + path)
	if err != nil {
		return nil, err
	}

	return c.newRequest(method, u, params, body)
}

// NewRequest prepares new request to common shodan api.
func (c *Client) NewRequest(method string, path string, params interface{}, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, err
	}

	return c.newRequest(method, u, params, body)
}

// NewStreamingRequest prepares new request to streaming api.
func (c *Client) NewStreamingRequest(method string, path string, params interface{}, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(c.StreamBaseURL + path)
	if err != nil {
		return nil, err
	}

	return c.newRequest(method, u, params, body)
}

func (c *Client) newRequest(method string, u *url.URL, params interface{}, body io.Reader) (*http.Request, error) {
	qs, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	qs.Add("key", c.Token)

	u.RawQuery = qs.Encode()

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	return req, nil
}

// DoStream executes streaming request.
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

// Do executes common (non-streaming) request.
func (c *Client) Do(ctx context.Context, req *http.Request, destination interface{}) error {
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
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return getErrorFromResponse(resp)
	}

	if destination == nil {
		return nil
	}

	return c.parseResponse(destination, resp.Body)
}

func (c *Client) parseResponse(destination interface{}, body io.Reader) error {
	var err error

	if w, ok := destination.(io.Writer); ok {
		_, err = io.Copy(w, body)
	} else {
		decoder := json.NewDecoder(body)
		err = decoder.Decode(destination)
	}

	return err
}
