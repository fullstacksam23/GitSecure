package models

type ScanJob struct {
	JobID      string `json:"job_id"`
	Repo       string `json:"repo"`
	Status     string `json:"status"`
	CommitHash string `json:"commit_hash"`
}

type BatchJob struct {
	BatchID        string `json:"batch_id"`
	Language       string `json:"language"`
	Status         string `json:"status"`
	RepoCount      int    `json:"repo_count"`
	CompletedRepos int    `json:"completed_repos"`
}
