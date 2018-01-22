package shodan

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"fmt"
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
		StreamChan:     make(chan HostData),
		Client:         client,
	}
}

// NewEnvClient creates new Shodan client using environment variable
// SHODAN_KEY as the token.
func NewEnvClient(client *http.Client) *Client {
	return NewClient(client, os.Getenv("SHODAN_KEY"))
}

// SetDebug toggles the debug mode
func (c *Client) SetDebug(debug bool) {
	c.Debug = debug
}

func (c *Client) buildURL(base, path string, params interface{}) string {
	baseURL, err := url.Parse(base + path)
	if err != nil {
		panic(fmt.Sprintf("Error: %s. This must never happen!", err))
	}

	qs, err := query.Values(params)
	if err != nil {
		panic(fmt.Sprintf("Error: %s. BaseURL: %s. This must never happen!", err, baseURL.String()))
	}

	qs.Add("key", c.Token)

	baseURL.RawQuery = qs.Encode()

	return baseURL.String()
}

func (c *Client) buildBaseURL(path string, params interface{}) string {
	return c.buildURL(c.BaseURL, path, params)
}

func (c *Client) buildExploitBaseURL(path string, params interface{}) string {
	return c.buildURL(c.ExploitBaseURL, path, params)
}

func (c *Client) buildStreamBaseURL(path string, params interface{}) string {
	return c.buildURL(c.StreamBaseURL, path, params)
}

func (c *Client) sendRequest(method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	if c.Debug {
		if command, err := http2curl.GetCurlCommand(req); err == nil {
			log.Printf("shodan.sendRequest: %s\n", command)
		}
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, getErrorFromResponse(res)
	}

	return res, nil
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

func (c *Client) executeRequest(method, path string, destination interface{}, body io.Reader) error {
	res, err := c.sendRequest(method, path, body)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if destination == nil {
		return nil
	}

	return c.parseResponse(destination, res.Body)
}

func (c *Client) executeStreamRequest(method, path string, ch chan []byte) error {
	res, err := c.sendRequest(method, path, nil)
	if err != nil {
		return err
	}

	go func() {
		reader := bufio.NewReader(res.Body)

		for {
			chunk, err := reader.ReadBytes('\n')
			if err != nil {
				res.Body.Close()
				close(ch)
				break
			}

			ch <- chunk
		}
	}()

	return nil
}
