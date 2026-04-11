package scanner

type Repo struct {
	Name        string   `json:"name"`
	FullName    string   `json:"full_name"`
	Stars       int      `json:"stargazers_count"`
	URL         string   `json:"html_url"`
	Description string   `json:"description"`
	Topics      []string `json:"topics"`

	DefaultBranch string `json:"default_branch"`
}

type BatchResponse struct {
	Items []Repo `json:"items"`
}

type CommitResponse struct {
	SHA string `json:"sha"`
}
