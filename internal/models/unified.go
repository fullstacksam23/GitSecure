package models

type UnifiedVuln struct {
	ID       string   `json:"id"`
	Package  string   `json:"package"`
	Version  string   `json:"version"`
	Severity string   `json:"severity"`
	Summary  string   `json:"summary"`
	Urls     []string `json:"urls"`

	FixVersion []string `json:"fix_version"`
	FixState   string   `json:"fix_state"`

	Risk      float64 `json:"risk"`
	Namespace string  `json:"namespace"`

	MatchType  string `json:"match_type"`
	Constraint string `json:"constraint"`

	DataSource string `json:"data_source"`
	Source     string `json:"source"`

	CWEIDs    []string `json:"cwe_ids"`
	Ecosystem string   `json:"ecosystem"`
}
