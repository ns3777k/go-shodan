package shodan

import (
	"testing"
	"net/http"

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
		HTTPS: true,
		Unlocked: true,
		UnlockedLeft: 9999,
		Telnet: false,
		ScanCredits: 254,
		Plan: "basic",
		QueryCredits: 2341,
	}

	assert.Nil(t, err);
	assert.EqualValues(t, infoExpected, info)
}
