package models

const MimeTypeSearchResultV1 = "application/vnd.searchResult.v1+json"

// Schema: searchResult.v1
type SearchResultV1 struct {
	Total int                     `json:"total"`
	Data  []*EnrichedPersonDataV1 `json:"data,omitempty"`
}
