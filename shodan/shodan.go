package shodan

import (
	"net/http"
	"io"
	"encoding/json"
	"net/url"
	"io/ioutil"
	"errors"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	baseURL = "https://api.shodan.io"
	exploitBaseURL = "https://exploits.shodan.io/api"
	streamBaseURL = "https://stream.shodan.io"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func getErrorFromResponse (r *http.Response) error {
	var errorResponse ErrorResponse
	message, err := ioutil.ReadAll(r.Body)
	if err == nil {
		if err := json.Unmarshal(message, &errorResponse); err == nil {
			return errors.New(errorResponse.Error)
		} else {
			return errors.New(strings.TrimSpace(string(message)))
		}
	}

	return ErrBodyRead
}

type Client struct {
	Token string

	client *http.Client
}

func NewClient(token string) *Client {
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	client := &http.Client{Transport: transport}

	return &Client{
		Token: token,
		client: client,
	}
}

func (c *Client) buildURL(base, path string, params interface{}) (string, error) {
	baseURL, err := url.Parse(base + path)
	if err != nil {
		return "", err
	}

	qs, err := query.Values(params)
	if err != nil {
		return baseURL.String(), err
	}

	qs.Add("key", c.Token)

	baseURL.RawQuery = qs.Encode()

	return baseURL.String(), nil
}

func (c *Client) buildBaseURL(path string, params interface{}) (string, error) {
	return c.buildURL(baseURL, path, params)
}

func (c *Client) buildExploitBaseURL(path string, params interface{}) (string, error) {
	return c.buildURL(exploitBaseURL, path, params)
}

func (c *Client) buildStreamBaseURL(path string, params interface{}) (string, error) {
	return c.buildURL(streamBaseURL, path, params)
}

func (c *Client) executeRequest(method, path string, destination interface{}) error {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return getErrorFromResponse(res)
	}

	if w, ok := destination.(io.Writer); ok {
		io.Copy(w, res.Body)
	} else {
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(destination)
	}

	return err
}
