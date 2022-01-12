package shodan

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"testing"
)

func TestClient_GeoPing(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	ip := "127.0.0.1"
	path := fmt.Sprintf(geonetPingPath, ip)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "geonet/ping")) //nolint:errcheck
	})

	pingResult, err := client.GeoPing(context.TODO(), net.ParseIP(ip))
	expectedPingResult := &PingResult{
		IP:              "127.0.0.1",
		IsAlive:         true,
		MinRTT:          0.271,
		MaxRTT:          0.307,
		AvgRTT:          0.289,
		RTTs:            []float64{0.3066062927246094, 0.2884864807128906, 0.27060508728027344},
		PacketsSent:     3,
		PacketsReceived: 3,
		PacketLoss:      0.0,
		FromLocation:    &Location{City: "Santa Clara", Country: "US", Coords: "37.3924,-121.9623"},
		Error:           "",
	}

	assert.Nil(t, err)
	assert.Equal(t, expectedPingResult, pingResult)
}

func TestClient_GeoPings(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	ip := "127.0.0.1"
	path := fmt.Sprintf(geonetPingsPath, ip)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "geonet/ping_multiple")) //nolint:errcheck
	})

	pingResult, err := client.GeoPings(context.TODO(), net.ParseIP(ip))
	expectedPingResult := []*PingResult{
		{
			IP:              "127.0.0.1",
			IsAlive:         true,
			MinRTT:          0.271,
			MaxRTT:          0.307,
			AvgRTT:          0.289,
			RTTs:            []float64{0.3066062927246094, 0.2884864807128906, 0.27060508728027344},
			PacketsSent:     3,
			PacketsReceived: 3,
			PacketLoss:      0.0,
			FromLocation:    &Location{City: "Santa Clara", Country: "US", Coords: "37.3924,-121.9623"},
			Error:           "",
		},
		{
			IP:              "127.0.0.1",
			IsAlive:         true,
			MinRTT:          0.197,
			MaxRTT:          0.268,
			AvgRTT:          0.229,
			RTTs:            []float64{0.19741058349609375, 0.2682209014892578, 0.22220611572265625},
			PacketsSent:     3,
			PacketsReceived: 3,
			PacketLoss:      0.0,
			FromLocation:    &Location{City: "London", Country: "GB", Coords: "51.5085,-0.1257"},
			Error:           "",
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, expectedPingResult, pingResult)
}

func TestClient_GeoDNSQuery(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	hostname := "test.com"
	path := fmt.Sprintf(geonetDNSPath, hostname)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "A", r.URL.Query().Get("rtype"))

		w.Write(getStub(t, "geonet/dns")) //nolint:errcheck
	})

	queryResult, err := client.GeoDNSQuery(context.TODO(), hostname, &DNSQueryOptions{RecordType: "A"})
	expectedQueryResult := &DNSQueryResult{
		Answers:      []*DNSQueryResultAnswers{{Type: "A", Value: "67.225.146.248"}},
		FromLocation: &Location{City: "Santa Clara", Country: "US", Coords: "37.3924,-121.9623"},
		Error:        "",
	}

	assert.Nil(t, err)
	assert.Equal(t, expectedQueryResult, queryResult)
}

func TestClient_GeoDNSQueries(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	hostname := "test.com"
	path := fmt.Sprintf(geonetDNSsPath, hostname)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "A", r.URL.Query().Get("rtype"))

		w.Write(getStub(t, "geonet/dns_multiple")) //nolint:errcheck
	})

	queryResult, err := client.GeoDNSQueries(context.TODO(), hostname, &DNSQueryOptions{RecordType: "A"})
	expectedQueryResult := []*DNSQueryResult{
		{
			Answers:      []*DNSQueryResultAnswers{{Type: "A", Value: "67.225.146.248"}},
			FromLocation: &Location{City: "Santa Clara", Country: "US", Coords: "37.3924,-121.9623"},
			Error:        "",
		},
		{
			Answers:      []*DNSQueryResultAnswers{{Type: "A", Value: "67.225.146.248"}},
			FromLocation: &Location{City: "Clifton", Country: "US", Coords: "40.8344,-74.1377"},
			Error:        "",
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, expectedQueryResult, queryResult)
}
