package shodan

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"io/ioutil"
	"fmt"
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

	if client.Token != testClientToken {
		t.Errorf("NewClient Token is %v, expected %v", client.Token, testClientToken)
	}
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
		if err != nil {
			t.Errorf("buildURL returned error %v", err)
		}

		if caseParams.expected != url {
			t.Errorf("buildURL returned invalid url, expected %v, actual %v", caseParams.expected, url)
		}
	}
}
