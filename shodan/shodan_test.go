package shodan

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"io/ioutil"
	"fmt"

	"github.com/stretchr/testify/assert"
)

const (
	testClientToken = "TEST_TOKEN"
	stubsDir = "stubs"
)

var (
	mux *http.ServeMux
	server *httptest.Server
	client *Client
)

func setUpTestServe() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client = NewClient(testClientToken)
	client.BaseURL = server.URL
	client.ExploitBaseURL = server.URL
	client.StreamBaseURL = server.URL
}

func getStub(t *testing.T, stubName string) []byte {
	stubPath := fmt.Sprintf("%s/%s.json", stubsDir, stubName)
	content, err := ioutil.ReadFile(stubPath)
	if err != nil {
		t.Errorf("getStub error %v", err)
	}

	return content
}

func tearDownTestServe() {
	server.Close()
}

func TestNewClient(t *testing.T) {
	client := NewClient(testClientToken)

	assert.Equal(t, testClientToken, client.Token)
}

func TestClient_buildURL_valid(t *testing.T) {
	client := NewClient(testClientToken)
	testOptions := struct {
		Page int `url:"page"`
		ShowAll bool `url:"show_all"`
	}{
		100,
		true,
	}
	testCases := []struct {
		path string
		params interface{}
		expected string
	}{
		{
			"/testing/test/1",
			nil,
			baseURL + "/testing/test/1?key=" + testClientToken,
		},
		{
			"/testing/test/2",
			testOptions,
			baseURL + "/testing/test/2?key=" + testClientToken + "&page=100&show_all=true",
		},
	}

	for _, caseParams := range testCases {
		url, err := client.buildURL(baseURL, caseParams.path, caseParams.params)

		assert.Nil(t, err);
		assert.Equal(t, caseParams.expected, url)
	}
}
