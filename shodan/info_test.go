package shodan

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetAPIInfo(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(infoPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "info"))
	})

	info, err := client.GetAPIInfo()
	infoExpected := &APIInfo{
		HTTPS:        true,
		Unlocked:     true,
		UnlockedLeft: 9999,
		Telnet:       false,
		ScanCredits:  254,
		Plan:         "basic",
		QueryCredits: 2341,
	}

	assert.Nil(t, err)
	assert.IsType(t, infoExpected, info)
	assert.EqualValues(t, infoExpected, info)
}

func TestClient_GetAPIInfo_invalidBaseURL(t *testing.T) {
	client := NewClient(nil, testClientToken)
	client.BaseURL = ":/1232.22"
	_, err := client.GetAPIInfo()
	assert.NotNil(t, err)
}
