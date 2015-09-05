package shodan

import (
	"net/http"
	"io"
	"encoding/json"
	"net/url"

	"github.com/google/go-querystring/query"
)

const (
	baseUrl = "https://api.shodan.io"
)

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

func (c *Client) buildUrl(path string, params interface{}) (string, error) {
	baseUrl, err := url.Parse(baseUrl + path)
	if err != nil {
		return "", err
	}

	qs, err := query.Values(params)
	if err != nil {
		return baseUrl.String(), err
	}

	qs.Add("key", c.Token)

	baseUrl.RawQuery = qs.Encode()

	return baseUrl.String(), nil
}

func (c *Client) executeRequest(method, path string, v interface{}) error {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if w, ok := v.(io.Writer); ok {
		io.Copy(w, res.Body)
	} else {
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(v)
	}

	return err
}
