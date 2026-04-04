package models

type ScanJob struct {
	JobID      string `json:"job_id"`
	Repo       string `json:"repo"`
	Status     string `json:"status"`
	CommitHash string `json:"commit_hash"`
}
