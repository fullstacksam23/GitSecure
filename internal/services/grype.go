package services

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

	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, err
	}

	return output, nil
}
