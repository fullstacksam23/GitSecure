package models

type OSVAdvisory struct {
	ID               string   `json:"id"`
	Summary          string   `json:"summary"`
	Details          string   `json:"details"`
	Published        string   `json:"published"`
	Aliases          []string `json:"aliases"`
	DatabaseSpecific struct {
		Severity       string   `json:"severity"`
		NVDPublishedAt *string  `json:"nvd_published_at"`
		CWEIDs         []string `json:"cwe_ids"`
		GithubReviewed bool     `json:"github_reviewed"`
	} `json:"database_specific"`

	References []struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"references"`

	Affected []struct {
		Package struct {
			Name      string `json:"name"`
			Ecosystem string `json:"ecosystem"`
			Purl      string `json:"purl"`
		} `json:"package"`

		Ranges []struct {
			Type   string `json:"type"`
			Events []struct {
				Introduced string `json:"introduced,omitempty"`
				Fixed      string `json:"fixed,omitempty"`
			} `json:"events"`
		} `json:"ranges"`

		DatabaseSpecific struct {
			Source string `json:"source"`
		} `json:"database_specific"`
	} `json:"affected"`

	Severity []struct {
		Type  string `json:"type"`
		Score string `json:"score"`
	} `json:"severity"`

	SchemaVersion string `json:"schema_version"`
}

type Purl struct {
	Purl string `json:"purl"`
}

type PackageQuery struct {
	Package Purl `json:"package"`
}

type OSVRequest struct {
	Queries []PackageQuery `json:"queries"`
}

type Vulnerability struct {
	Id       string `json:"id"`
	Modified string `json:"modified"`
}

type Result struct {
	Vulns []Vulnerability `json:"vulns"`
}

type OSVResponse struct {
	Results []Result `json:"results"`
}
