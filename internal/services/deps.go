package services

import (
	"encoding/json"
	"net/http"

	"github.com/fullstacksam23/GitSecure/internal/models"
)

type Package struct {
	Name             string
	Version          string
	ReferenceType    string
	ReferenceLocator string
}

func parseSbom(url string) (models.SBOMResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return models.SBOMResponse{}, err
	}
	defer resp.Body.Close()
	var sbomResp models.SBOMResponse

	err = json.NewDecoder(resp.Body).Decode(&sbomResp)

	return sbomResp, nil
}
func ExtractDependencies(url string) ([]Package, error) {
	sbom, err := parseSbom(url)
	if err != nil {
		return nil, err
	}
	packages := []Package{}
	for _, p := range sbom.SBOM.Packages {
		packages = append(packages, Package{
			Name:             p.Name,
			Version:          p.VersionInfo,
			ReferenceType:    p.ExternalRefs[0].ReferenceType,
			ReferenceLocator: p.ExternalRefs[0].ReferenceLocator,
		})
	}
	return packages, nil
}
