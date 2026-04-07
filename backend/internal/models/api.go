package models

type SeverityBreakdown struct {
	Critical int64 `json:"critical"`
	High     int64 `json:"high"`
	Medium   int64 `json:"medium"`
	Low      int64 `json:"low"`
	Unknown  int64 `json:"unknown"`
}

type FacetCount struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

type PackageRiskItem struct {
	Package            string  `json:"package"`
	Ecosystem          string  `json:"ecosystem"`
	Risk               float64 `json:"risk"`
	VulnerabilityCount int64   `json:"vulnerability_count"`
}

type RepoSummary struct {
	Repo        string `json:"repo"`
	Status      string `json:"status"`
	LastJobID   string `json:"last_job_id"`
	LastScanned string `json:"last_scanned"`
	TopSeverity string `json:"top_severity"`
}

type TrendPoint struct {
	Date            string `json:"date"`
	Scans           int64  `json:"scans"`
	Vulnerabilities int64  `json:"vulnerabilities"`
}

type DashboardSummary struct {
	TotalScans           int64             `json:"total_scans"`
	TotalVulnerabilities int64             `json:"total_vulnerabilities"`
	Critical             int64             `json:"critical"`
	High                 int64             `json:"high"`
	Medium               int64             `json:"medium"`
	Low                  int64             `json:"low"`
	PackagesFixed        int64             `json:"packages_fixed"`
	SeverityDistribution SeverityBreakdown `json:"severity_distribution"`
	RecentScans          []ScanListItem    `json:"recent_scans"`
	TopRiskPackages      []PackageRiskItem `json:"top_risk_packages"`
	RiskTrend            []TrendPoint      `json:"risk_trend"`
	RepoSummaries        []RepoSummary     `json:"repo_summaries"`
}

type ScanListItem struct {
	JobID              string            `json:"job_id"`
	Repo               string            `json:"repo"`
	Status             string            `json:"status"`
	CommitHash         string            `json:"commit_hash"`
	CreatedAt          string            `json:"created_at"`
	TopSeverity        string            `json:"top_severity"`
	VulnerabilityCount int64             `json:"vulnerability_count"`
	SeverityCounts     SeverityBreakdown `json:"severity_counts"`
}

type ScanDetails struct {
	JobID              string            `json:"job_id"`
	Repo               string            `json:"repo"`
	Status             string            `json:"status"`
	CommitHash         string            `json:"commit_hash"`
	CreatedAt          string            `json:"created_at"`
	VulnerabilityCount int64             `json:"vulnerability_count"`
	SeverityCounts     SeverityBreakdown `json:"severity_counts"`
	Ecosystems         []FacetCount      `json:"ecosystems"`
	FixStates          []FacetCount      `json:"fix_states"`
	TopPackages        []PackageRiskItem `json:"top_packages"`
}

type ScanListResponse struct {
	Items      []ScanListItem `json:"items"`
	Pagination Pagination     `json:"pagination"`
}

type VulnerabilityRecord struct {
	ID                 string   `json:"id"`
	JobID              string   `json:"job_id"`
	Package            string   `json:"package"`
	Version            string   `json:"version"`
	Severity           string   `json:"severity"`
	NormalizedSeverity string   `json:"normalized_severity"`
	Summary            string   `json:"summary"`
	Urls               []string `json:"urls"`
	FixVersion         []string `json:"fix_version"`
	FixState           string   `json:"fix_state"`
	Risk               float64  `json:"risk"`
	Namespace          string   `json:"namespace"`
	MatchType          string   `json:"match_type"`
	VersionConstraint  string   `json:"version_constraint"`
	DataSource         string   `json:"data_source"`
	Source             string   `json:"source"`
	CWEIDs             []string `json:"cwe_ids"`
	Ecosystem          string   `json:"ecosystem"`
	CreatedAt          string   `json:"created_at"`
}

type VulnerabilityFacets struct {
	Ecosystems []FacetCount `json:"ecosystems"`
	FixStates  []FacetCount `json:"fix_states"`
}

type VulnerabilityListResponse struct {
	Items      []VulnerabilityRecord `json:"items"`
	Pagination Pagination            `json:"pagination"`
	Facets     VulnerabilityFacets   `json:"facets"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type CompareBucket struct {
	Count int64                 `json:"count"`
	Items []VulnerabilityRecord `json:"items"`
}

type ScanCompareResponse struct {
	BaseScan      *ScanDetails      `json:"base_scan"`
	TargetScan    *ScanDetails      `json:"target_scan"`
	New           CompareBucket     `json:"new"`
	Fixed         CompareBucket     `json:"fixed"`
	Persisting    CompareBucket     `json:"persisting"`
	NewSeverity   SeverityBreakdown `json:"new_severity"`
	FixedSeverity SeverityBreakdown `json:"fixed_severity"`
}
