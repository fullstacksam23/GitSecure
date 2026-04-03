package grype

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fullstacksam23/GitSecure/internal/models"
)

func GrypeScan(sbom []byte) ([]byte, error) {

	var wrapper models.SBOMResponse
	err := json.Unmarshal(sbom, &wrapper)
	if err != nil {
		return nil, err
	}

	actualSBOM, err := json.Marshal(wrapper.SBOM)
	if err != nil {
		return nil, err
	}

	tmpDir := os.TempDir()
	sbomPath := filepath.Join(tmpDir, "sbom.json")

	err = os.WriteFile(sbomPath, actualSBOM, 0644)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("grype", "sbom:"+sbomPath, "-o", "json")

	output, err := cmd.CombinedOutput() //returns both the std output and err
	if err != nil {
		return output, err
	}

	return output, nil
}

func ParseGrype(output []byte) (GrypeResponse, error) {

	var resp GrypeResponse

	err := json.Unmarshal(output, &resp)
	return resp, err
}

func NormalizeGrype(grype GrypeResponse, canonical map[string]string, jobID string) []models.UnifiedVuln {

	var vulns []models.UnifiedVuln

	for _, match := range grype.Matches {

		id := match.Vulnerability.ID
		if canonicalID, ok := canonical[id]; ok && canonicalID != "" {
			id = canonicalID
		}

		// Extract match details safely
		matchType, constraint := pickBestMatch(match.MatchDetails)

		v := models.UnifiedVuln{
			ID:       id,
			JobID:    jobID,
			Package:  match.Artifact.Name,
			Version:  match.Artifact.Version,
			Severity: match.Vulnerability.Severity,
			Summary:  match.Vulnerability.Description,
			Urls:     match.Vulnerability.Urls,

			FixVersion: match.Vulnerability.Fix.Versions,
			FixState:   match.Vulnerability.Fix.State,

			Risk:      match.Vulnerability.Risk,
			Namespace: match.Vulnerability.Namespace,

			MatchType:  matchType,
			Constraint: constraint,

			DataSource: match.Vulnerability.DataSource,
			Source:     "grype",
		}
		vulns = append(vulns, v)
	}

	return vulns
}

func pickBestMatch(details []MatchDetail) (string, string) {
	bestType := ""
	constraint := ""

	for _, d := range details {
		if d.Type == "exact-direct-match" {
			return d.Type, d.Found.VersionConstraint
		}
		if d.Type == "exact-indirect-match" {
			bestType = d.Type
			constraint = d.Found.VersionConstraint
		}
	}

	if bestType != "" {
		return bestType, constraint
	}

	if len(details) > 0 {
		return details[0].Type, details[0].Found.VersionConstraint
	}

	return "", ""
}
