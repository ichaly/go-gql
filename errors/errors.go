package errors

type QueryError struct {
	Err           error                  `json:"-"` // Err holds underlying if available
	Message       string                 `json:"message"`
	Locations     []Location             `json:"locations,omitempty"`
	Path          []interface{}          `json:"path,omitempty"`
	Rule          string                 `json:"-"`
	ResolverError error                  `json:"-"`
	Extensions    map[string]interface{} `json:"extensions,omitempty"`
}

type Location struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}
