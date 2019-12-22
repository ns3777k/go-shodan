package shodan

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetNotifierProviders(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(notifierProviderPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "notifiers/providers")) //nolint:errcheck
	})

	providers, err := client.GetNotifierProviders(context.TODO())
	expected := map[string]*NotifierProvider{
		"pagerduty": {Required: []string{"routing_key"}},
		"email":     {Required: []string{"to"}},
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, providers)
}

func TestClient_GetNotifiers(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(notifierPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "notifiers/notifiers")) //nolint:errcheck
	})

	notifiers, err := client.GetNotifiers(context.TODO())
	expected := []*Notifier{
		{Description: "", Provider: "email", ID: "default", Args: map[string]string{"to": "ns3777k@gmail.com"}},
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, notifiers)
}

func TestClient_GetNotifier(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	id := "FKDSHFKS47D"
	mux.HandleFunc(notifierPath+"/"+id, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "notifiers/notifier")) //nolint:errcheck
	})

	notifier, err := client.GetNotifier(context.TODO(), id)
	expected := &Notifier{
		Description: "", Provider: "email", ID: id, Args: map[string]string{"to": "ns3777k@gmail.com"},
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, notifier)
}

func TestClient_DeleteNotifier(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	id := "default"
	mux.HandleFunc(notifierPath+"/"+id, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		fmt.Fprint(w, `{"success": true}`) //nolint:errcheck
	})

	isOK, err := client.DeleteNotifier(context.TODO(), id)

	assert.Nil(t, err)
	assert.True(t, isOK)
}

func TestClient_CreateNotifier(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(notifierPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		fmt.Fprint(w, `{"id": "GREY545DYUDYU3432", "success": true}`) //nolint:errcheck
	})

	notifier := &Notifier{
		Provider:    "email",
		Description: "test",
		Args:        nil,
	}

	isOK, err := client.CreateNotifier(context.TODO(), notifier)

	assert.Nil(t, err)
	assert.True(t, isOK)
	assert.Equal(t, "GREY545DYUDYU3432", notifier.ID)
}

func TestClient_UpdateNotifierArgs(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	id := "default"
	mux.HandleFunc(notifierPath+"/"+id, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		fmt.Fprint(w, `{"success": true}`) //nolint:errcheck
	})

	isOK, err := client.UpdateNotifierArgs(context.TODO(), id, map[string]string{"to": "test@test.local"})

	assert.Nil(t, err)
	assert.True(t, isOK)
}
