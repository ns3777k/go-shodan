package shodan

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetAPIInfo(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(infoPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "info")) //nolint:errcheck
	})

	info, err := client.GetAPIInfo(context.TODO())
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
