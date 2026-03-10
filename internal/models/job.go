package models

type ScanJob struct {
	JobID string `json:"job_id"`
	Repo  string `json:"repo"`
}
