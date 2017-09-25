package shodan

// Facet is a property to get summary information on.
type Facet struct {
	Count int    `json:"count"`
	Value string `json:"value"`
}
