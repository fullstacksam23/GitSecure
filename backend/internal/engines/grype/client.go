package grype

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fullstacksam23/GitSecure/internal/models"
	"github.com/google/uuid"
)

func GrypeScan(sbom []byte) ([]byte, error) {
	tmpDir := os.TempDir()
	sbomPath := filepath.Join(tmpDir, "sbom-"+uuid.New().String()+".json")

	err := os.WriteFile(sbomPath, sbom, 0644)
	if err != nil {
		return nil, err
	}
	defer os.Remove(sbomPath)

	cmd := exec.Command("grype", "sbom:"+sbomPath, "-o", "json")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("grype failed: %w\n%s", err, string(output))
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
