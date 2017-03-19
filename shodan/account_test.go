package shodan

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetAccountProfile(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(profilePath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "profile"))
	})

	account, err := client.GetAccountProfile()
	accountExpected := &Profile{
		Member:  true,
		Name:    "",
		Credits: 40,
		Created: "2015-09-03T12:44:29.278000",
	}

	assert.Nil(t, err)
	assert.IsType(t, accountExpected, account)
	assert.EqualValues(t, accountExpected, account)
}
