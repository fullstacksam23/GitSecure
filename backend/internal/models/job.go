package models

type ScanJob struct {
	JobID      string `json:"job_id"`
	BatchID    string `json:"batch_id"`
	Repo       string `json:"repo"`
	Status     string `json:"status"`
	CommitHash string `json:"commit_hash"`
	RepoID     int    `json:"repo_id"`
	JobType    string `json:"job_type"`
}

type BatchJob struct {
	BatchID        string `json:"batch_id"`
	Language       string `json:"language"`
	Status         string `json:"status"`
	RepoCount      int    `json:"repo_count"`
	CompletedRepos int    `json:"completed_repos"`
	TotalRepos     int    `json:"total_repos"` //the total repos returned by github
}
