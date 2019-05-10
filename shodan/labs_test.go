package shodan

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_CalcHoneyScore(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	ip := net.ParseIP("192.168.0.1")

	mux.HandleFunc(fmt.Sprintf(honeyscorePath, ip), func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		fmt.Fprint(w, `0.5`)
	})

	score, err := client.CalcHoneyScore(context.TODO(), ip)
	assert.Nil(t, err)
	assert.Equal(t, 0.5, score)
}
