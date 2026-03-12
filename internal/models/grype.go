package models

type GrypeResponse struct {
	Matches []Match `json:"matches"`
}

type Match struct {
	Vulnerability GrypeVulnerability `json:"vulnerability"`
	Artifact      Artifact           `json:"artifact"`
}

type GrypeVulnerability struct {
	ID          string   `json:"id"`
	Severity    string   `json:"severity"`
	Description string   `json:"description"`
	Urls        []string `json:"urls"`

	CVSS []struct {
		Metrics struct {
			BaseScore float64 `json:"baseScore"`
		} `json:"metrics"`
	} `json:"cvss"`

	Fix struct {
		Versions []string `json:"versions"`
	} `json:"fix"`
}

type Artifact struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Purl    string `json:"purl"`
}
