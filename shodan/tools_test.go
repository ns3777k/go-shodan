package shodan

import (
	"testing"
	"net/http"
	"strconv"
	"fmt"
)

func TestClient_GetMyIP(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	testIP := "192.168.22.34"

	mux.HandleFunc(ipPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, strconv.Quote(testIP))
	})

	ip, err := client.GetMyIP()
	if err != nil {
		t.Errorf("Client executeRequest returned error %v", err)
	}

	if ip != testIP {
		t.Errorf("IPs are not equal, actual %v expected %s", ip, testIP)
	}
}
