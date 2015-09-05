package shodan

const (
	queryTagsPath = "/shodan/query/tags"
	querySearchPath = "/shodan/query/search"
	queryPath = "/shodan/query"
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
	Size int `url:"size,omitempty"`
}

type SearchQueryOptions struct {
	Query string `url:"query"`
	Page  int    `url:"page,omitempty"`
}

type QueryOptions struct {
	Page  int    `url:"page,omitempty"`
	Sort  string `url:"sort,omitempty"`
	Order string `url:"order,omitempty"`
}

func (c *Client) GetQueryTags(options *QueryTagsOptions) (*QueryTags, error) {
	url, err := c.buildUrl(queryTagsPath, options)
	if err != nil {
		return nil, err
	}

	var queryTags QueryTags
	err = c.executeRequest("GET", url, &queryTags)

	return &queryTags, err
}

func (c *Client) GetQueries(options *QueryOptions) (*QuerySearch, error) {
	url, err := c.buildUrl(queryPath, options)
	if err != nil {
		return nil, err
	}

	var querySearch QuerySearch
	err = c.executeRequest("GET", url, &querySearch)

	return &querySearch, err
}

func (c *Client) SearchQueries(options *SearchQueryOptions) (*QuerySearch, error) {
	url, err := c.buildUrl(querySearchPath, options)
	if err != nil {
		return nil, err
	}

	var querySearch QuerySearch
	err = c.executeRequest("GET", url, &querySearch)

	return &querySearch, err
}
