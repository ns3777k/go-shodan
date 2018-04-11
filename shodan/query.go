package shodan

import (
	"context"
)

const (
	queryTagsPath   = "/shodan/query/tags"
	querySearchPath = "/shodan/query/search"
	queryPath       = "/shodan/query"
)

// QueryTagsMatch represents a matched tag.
type QueryTagsMatch struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// QueryTags represents matched tags.
type QueryTags struct {
	Total   int               `json:"total"`
	Matches []*QueryTagsMatch `json:"matches"`
}

// QuerySearchMatch is a match of QuerySearch.
type QuerySearchMatch struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Query       string   `json:"query"`
	Votes       int      `json:"votes"`
	Timestamp   string   `json:"timestamp"`
	Tags        []string `json:"tags"`
}

// QuerySearch is the results of querying saved search queries.
type QuerySearch struct {
	Total   int                 `json:"total"`
	Matches []*QuerySearchMatch `json:"matches"`
}

// QueryTagsOptions represents options for GetQueryTags.
type QueryTagsOptions struct {
	// The number of tags to return (default: 10).
	Size int `url:"size,omitempty"`
}

// SearchQueryOptions is options for SearchQueries.
type SearchQueryOptions struct {
	// What to search for in the directory of saved search queries.
	Query string `url:"query"`

	// Page number to iterate over results; each page contains 10 items.
	Page int `url:"page,omitempty"`
}

// QueryOptions represents query options for fetching saved queries.
type QueryOptions struct {
	// Page number to iterate over results; each page contains 10 items.
	Page int `url:"page,omitempty"`

	// Sort the list based on a property. Possible values are: votes, timestamp.
	Sort string `url:"sort,omitempty"`

	// Whether to sort the list in ascending or descending order. Possible values are: asc, desc.
	Order string `url:"order,omitempty"`
}

// GetQueryTags obtains a list of popular tags for the saved search queries in Shodan.
func (c *Client) GetQueryTags(ctx context.Context, options *QueryTagsOptions) (*QueryTags, error) {
	var queryTags QueryTags

	req, err := c.NewRequest("GET", queryTagsPath, options, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &queryTags); err != nil {
		return nil, err
	}

	return &queryTags, nil
}

// GetQueries obtains a list of search queries that users have saved in Shodan.
func (c *Client) GetQueries(ctx context.Context, options *QueryOptions) (*QuerySearch, error) {
	var querySearch QuerySearch

	req, err := c.NewRequest("GET", queryPath, options, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &querySearch); err != nil {
		return nil, err
	}

	return &querySearch, nil
}

// SearchQueries searches the directory of search queries that users have saved in Shodan.
func (c *Client) SearchQueries(ctx context.Context, options *SearchQueryOptions) (*QuerySearch, error) {
	if options == nil || options.Query == "" {
		return nil, ErrInvalidQuery
	}

	var querySearch QuerySearch

	req, err := c.NewRequest("GET", querySearchPath, options, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &querySearch); err != nil {
		return nil, err
	}

	return &querySearch, nil
}
