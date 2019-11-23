package shodan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/google/go-querystring/query"
)

const (
	baseURL        = "https://api.shodan.io"
	exploitBaseURL = "https://exploits.shodan.io/api"
	streamBaseURL  = "https://stream.shodan.io"
)

// Client represents Shodan HTTP client
type Client struct {
	m *sync.Mutex

	Token          string
	BaseURL        string
	ExploitBaseURL string
	StreamBaseURL  string
	Debug          bool
	Client         *http.Client
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
		m:              &sync.Mutex{},
	}
}

// NewEnvClient creates new Shodan client using environment variable
// SHODAN_KEY as the token.
func NewEnvClient(client *http.Client) *Client {
	return NewClient(client, os.Getenv("SHODAN_KEY"))
}

// SetDebug toggles the debug mode.
func (c *Client) SetDebug(debug bool) {
	c.m.Lock()
	defer c.m.Unlock()

	c.Debug = debug
}

// NewExploitRequest prepares new request to exploit shodan api.
func (c *Client) NewExploitRequest(
	method string,
	path string,
	params interface{},
	body io.Reader,
) (*http.Request, error) {
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

func (c *Client) dumpRequest(req *http.Request) {
	var body string

	if req.Body != nil {
		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("[DEBUG] ns3777k/go-shodan: failed to read body: %s\n", err)
		}

		req.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
		body = string(bodyBytes)
	}

	message := fmt.Sprintf("%s %s %s", req.Method, req.URL.String(), body)
	log.Printf("[DEBUG] ns3777k/go-shodan: client request: %s\n", message)
}

func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	if c.Debug {
		c.dumpRequest(req)
	}

	return c.Client.Do(req)
}

// Do executes common (non-streaming) request.
func (c *Client) Do(ctx context.Context, req *http.Request, destination interface{}) error {
	resp, err := c.do(ctx, req)
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
