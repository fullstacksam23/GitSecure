package grype

type GrypeResponse struct {
	Matches []Match `json:"matches"`
}

type Match struct {
	Vulnerability GrypeVulnerability `json:"vulnerability"`
	Artifact      Artifact           `json:"artifact"`
	MatchDetails  []MatchDetail      `json:"matchDetails"`
}

type GrypeVulnerability struct {
	ID          string   `json:"id"`
	Severity    string   `json:"severity"`
	Description string   `json:"description"`
	Urls        []string `json:"urls"`

	DataSource string  `json:"dataSource"`
	Namespace  string  `json:"namespace"`
	Risk       float64 `json:"risk"`

	CVSS []struct {
		Metrics struct {
			BaseScore float64 `json:"baseScore"`
		} `json:"metrics"`
	} `json:"cvss"`

	Fix struct {
		Versions []string `json:"versions"`
		State    string   `json:"state"`
	} `json:"fix"`
}

type Artifact struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Purl    string `json:"purl"`
}

type MatchDetail struct {
	Type string `json:"type"`

	Found struct {
		VersionConstraint string `json:"versionConstraint"`
	} `json:"found"`
}
