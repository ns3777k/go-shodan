package shodan

const (
	queryTagsPath   = "/shodan/query/tags"
	querySearchPath = "/shodan/query/search"
	queryPath       = "/shodan/query"
)

type QueryTagsMatch struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

type QueryTags struct {
	Total   int               `json:"total"`
	Matches []*QueryTagsMatch `json:"matches"`
}

type QuerySearchMatch struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Query       string   `json:"query"`
	Votes       int      `json:"votes"`
	Timestamp   string   `json:"timestamp"`
	Tags        []string `json:"tags"`
}

type QuerySearch struct {
	Total   int                 `json:"total"`
	Matches []*QuerySearchMatch `json:"matches"`
}

type QueryTagsOptions struct {
	// The number of tags to return (default: 10)
	Size int `url:"size,omitempty"`
}

type SearchQueryOptions struct {
	// What to search for in the directory of saved search queries
	Query string `url:"query"`

	// Page number to iterate over results; each page contains 10 items
	Page int `url:"page,omitempty"`
}

type QueryOptions struct {
	// Page number to iterate over results; each page contains 10 items
	Page int `url:"page,omitempty"`

	// Sort the list based on a property. Possible values are: votes, timestamp
	Sort string `url:"sort,omitempty"`

	// Whether to sort the list in ascending or descending order. Possible values are: asc, desc
	Order string `url:"order,omitempty"`
}

// GetQueryTags obtains a list of popular tags for the saved search queries in Shodan
func (c *Client) GetQueryTags(options *QueryTagsOptions) (*QueryTags, error) {
	url, err := c.buildBaseURL(queryTagsPath, options)
	if err != nil {
		return nil, err
	}

	var queryTags QueryTags
	err = c.executeRequest("GET", url, &queryTags, nil)

	return &queryTags, err
}

// GetQueries obtains a list of search queries that users have saved in Shodan
func (c *Client) GetQueries(options *QueryOptions) (*QuerySearch, error) {
	url, err := c.buildBaseURL(queryPath, options)
	if err != nil {
		return nil, err
	}

	var querySearch QuerySearch
	err = c.executeRequest("GET", url, &querySearch, nil)

	return &querySearch, err
}

// SearchQueries searches the directory of search queries that users have saved in Shodan
func (c *Client) SearchQueries(options *SearchQueryOptions) (*QuerySearch, error) {
	if options.Query == "" {
		return nil, ErrInvalidQuery
	}

	url, err := c.buildBaseURL(querySearchPath, options)
	if err != nil {
		return nil, err
	}

	var querySearch QuerySearch
	err = c.executeRequest("GET", url, &querySearch, nil)

	return &querySearch, err
}
