package gofindit

type SearchQuery struct {
	Limit  int                         `json:"limit"`
	Skip   int                         `json:"skip"`
	Sort   string                      `json:"sort"`
	Fields map[string]SearchQueryField `json:"fields"`
}

type SearchQueryField struct {
	Type  string `json:"type"` // "match", "partial", "range"
	Value any    `json:"value"`
}
