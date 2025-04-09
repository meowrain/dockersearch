package models

type SearchResponse struct {
	NumPages   int      `json:"num_pages"`
	NumResults int      `json:"num_results"`
	PageSize   int      `json:"page_size"`
	Page       int      `json:"page"`
	Query      string   `json:"query"`
	Results    []Result `json:"results"`
}

type Result struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PullCount   int    `json:"pull_count"`
	StarCount   int    `json:"star_count"`
	IsTrusted   bool   `json:"is_trusted"`
	IsOfficial  bool   `json:"is_official"`
	IsAutomated bool   `json:"is_automated"`
}
