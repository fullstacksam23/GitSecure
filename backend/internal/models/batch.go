package models

type EcosystemRepo struct {
	ID       int64   `json:"id,omitempty"` // Database generates this
	BatchID  *string `json:"batch_id"`
	RepoName string  `json:"repo_name"`
	Stars    int     `json:"stars"`
	RepoRank int     `json:"repo_rank"`
}
