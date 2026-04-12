package models

type EcosystemBatchListItem struct {
	BatchID        string `json:"batch_id"`
	Language       string `json:"language"`
	Status         string `json:"status"`
	RepoCount      int64  `json:"repo_count"`
	TotalRepos     int64  `json:"total_repos"`
	CompletedRepos int64  `json:"completed_repos"`
	CreatedAt      string `json:"created_at"`
	CompletedAt    string `json:"completed_at"`
}

type EcosystemBatchDetail struct {
	BatchID        string `json:"batch_id"`
	Language       string `json:"language"`
	Status         string `json:"status"`
	RepoCount      int64  `json:"repo_count"`
	TotalRepos     int64  `json:"total_repos"`
	CompletedRepos int64  `json:"completed_repos"`
	CreatedAt      string `json:"created_at"`
	CompletedAt    string `json:"completed_at"`
}

type EcosystemBatchListResponse struct {
	Items      []EcosystemBatchListItem `json:"items"`
	Pagination Pagination               `json:"pagination"`
}

type EcosystemRepoListItem struct {
	ID                 int64             `json:"id"`
	BatchID            string            `json:"batch_id"`
	RepoName           string            `json:"repo_name"`
	Stars              int64             `json:"stars"`
	Rank               int64             `json:"rank"`
	ScanStatus         string            `json:"scan_status"`
	JobID              string            `json:"job_id"`
	VulnerabilityCount int64             `json:"vulnerability_count"`
	TopSeverity        string            `json:"top_severity"`
	RiskScore          float64           `json:"risk_score"`
	SeverityCounts     SeverityBreakdown `json:"severity_counts"`
}

type EcosystemRepoListResponse struct {
	Items      []EcosystemRepoListItem `json:"items"`
	Pagination Pagination              `json:"pagination"`
}

type EcosystemBatchSummary struct {
	BatchID              string                `json:"batch_id"`
	TotalRepositories    int64                 `json:"total_repositories"`
	TotalVulnerabilities int64                 `json:"total_vulnerabilities"`
	SeverityBreakdown    SeverityBreakdown     `json:"severity_breakdown"`
	AverageRiskScore     float64               `json:"average_risk_score"`
	MostVulnerableRepo   *EcosystemRepoSummary `json:"most_vulnerable_repo"`
}

type EcosystemRepoSummary struct {
	RepoName           string                `json:"repo_name"`
	JobID              string                `json:"job_id"`
	ScanStatus         string                `json:"scan_status"`
	VulnerabilityCount int64                 `json:"vulnerability_count"`
	TopSeverity        string                `json:"top_severity"`
	RiskScore          float64               `json:"risk_score"`
	SeverityCounts     SeverityBreakdown     `json:"severity_counts"`
	TopVulnerabilities []VulnerabilityRecord `json:"top_vulnerabilities"`
}
