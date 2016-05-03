package shodan

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_CalcHoneyScore(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	ip := "192.168.0.1"

	mux.HandleFunc(fmt.Sprintf(honeyscorePath, ip), func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		fmt.Fprint(w, `0.5`)
	})

	score, err := client.CalcHoneyScore(ip)
	assert.Nil(t, err)
	assert.Equal(t, 0.5, score)
}
