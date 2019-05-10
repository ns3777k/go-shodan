package shodan

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetOrganization(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(organizationPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "org")) //nolint:errcheck
	})

	org, err := client.GetOrganization(context.TODO())
	orgExpected := &Organization{
		ID:          "OIj8DS0lks9",
		Name:        "ACME",
		Created:     "2018-05-20T01:35:20.604000",
		UpgradeType: "dev",
		Domains:     []string{"acme.com"},
		Members:     []*OrganizationMember{},
		Admins: []*OrganizationMember{
			{Username: "acme-admin", Email: "admin@acme.com"},
		},
	}

	assert.Nil(t, err)
	assert.EqualValues(t, orgExpected, org)
}

func TestClient_AddMemberToOrganization(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	user := "username"
	path := fmt.Sprintf(organizationMemberPath, user)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		fmt.Fprint(w, `{"success": true}`) //nolint:errcheck
	})

	r, err := client.AddMemberToOrganization(context.TODO(), user, nil)

	assert.Nil(t, err)
	assert.True(t, r)
}

func TestClient_AddMemberToOrganizationWithNotifications(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	user := "testing_user"
	path := fmt.Sprintf(organizationMemberPath, user)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "true", q["notify"][0])
		fmt.Fprint(w, `{"success": true}`) //nolint:errcheck
	})

	r, err := client.AddMemberToOrganization(context.TODO(), user, &AddMemberToOrganizationOptions{Notify: true})

	assert.Nil(t, err)
	assert.True(t, r)
}

func TestClient_RemoveMemberFromOrganization(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	user := "removing_user"
	path := fmt.Sprintf(organizationMemberPath, user)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		fmt.Fprint(w, `{"success": true}`) //nolint:errcheck
	})

	r, err := client.RemoveMemberFromOrganization(context.TODO(), user)

	assert.Nil(t, err)
	assert.True(t, r)
}
