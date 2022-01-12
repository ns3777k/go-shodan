package shodan

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetQueryTags(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(queryTagsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "query_tags"))
	})

	queryTagsExpected := &QueryTags{
		Total: 3782,
		Matches: []*QueryTagsMatch{
			{
				Count: 76,
				Value: "webcam",
			},
			{
				Count: 68,
				Value: "scada",
			},
		},
	}
	queryTags, err := client.GetQueryTags(context.TODO(), new(QueryTagsOptions))

	assert.Nil(t, err)
	assert.EqualValues(t, queryTagsExpected, queryTags)
}

func TestClient_SearchQueries(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(querySearchPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "query_search_results"))
	})

	searchQueryExpected := &QuerySearch{
		Total: 2,
		Matches: []*QuerySearchMatch{
			{
				Votes:       2,
				Description: "apache servers US",
				Title:       "apache servers",
				Timestamp:   "2013-02-21T02:25:53.018000",
				Tags:        []string{"apache"},
				Query:       "apache country:US",
			},
			{
				Votes:       5,
				Description: "exacttouch ...smtp",
				Title:       "Centos apache",
				Timestamp:   "2010-03-07T15:47:13",
				Tags:        []string{},
				Query:       "country:in apache centos hostname:exacttouch.com",
			},
		},
	}
	searchQuery, err := client.SearchQueries(context.TODO(), &SearchQueryOptions{Query: "apache"})

	assert.Nil(t, err)
	assert.EqualValues(t, searchQueryExpected, searchQuery)
}

func TestClient_SearchQueries_nilOptions(t *testing.T) {
	_, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	_, err := client.SearchQueries(context.TODO(), nil)

	assert.NotNil(t, err)
	assert.IsType(t, ErrInvalidQuery, err)
}

func TestClient_SearchQueries_emptyQueryOption(t *testing.T) {
	_, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	_, err := client.SearchQueries(context.TODO(), &SearchQueryOptions{Query: ""})

	assert.NotNil(t, err)
	assert.IsType(t, ErrInvalidQuery, err)
}

func TestClient_GetQueries(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(queryPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "query_search_results"))
	})

	queriesExpected := &QuerySearch{
		Total: 2,
		Matches: []*QuerySearchMatch{
			{
				Votes:       2,
				Description: "apache servers US",
				Title:       "apache servers",
				Timestamp:   "2013-02-21T02:25:53.018000",
				Tags:        []string{"apache"},
				Query:       "apache country:US",
			},
			{
				Votes:       5,
				Description: "exacttouch ...smtp",
				Title:       "Centos apache",
				Timestamp:   "2010-03-07T15:47:13",
				Tags:        []string{},
				Query:       "country:in apache centos hostname:exacttouch.com",
			},
		},
	}
	queries, err := client.GetQueries(context.TODO(), new(QueryOptions))

	assert.Nil(t, err)
	assert.EqualValues(t, queriesExpected, queries)
}
