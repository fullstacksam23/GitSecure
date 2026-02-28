package services

import (
	"encoding/json"
	"io"

	"github.com/fullstacksam23/GitSecure/internal/models"
)

type Package struct {
	Name             string
	Version          string
	ReferenceType    string
	ReferenceLocator string
}

func parseSbom(r io.Reader) (models.SPDXDocument, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return models.SPDXDocument{}, err
	}

	// github sbom formatting returns a SBOMResponse
	var wrapper models.SBOMResponse
	if err := json.Unmarshal(data, &wrapper); err == nil && len(wrapper.SBOM.Packages) > 0 {
		return wrapper.SBOM, nil
	}

	// manual syft sbom formatting returns a SPDXDocument
	var doc models.SPDXDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return models.SPDXDocument{}, err
	}

	return doc, nil
}

func ExtractDependencies(r io.Reader) ([]Package, error) {
	doc, err := parseSbom(r)
	if err != nil {
		return nil, err
	}

	packages := []Package{}

	for _, p := range doc.Packages {
		var refType, refLoc string
		if len(p.ExternalRefs) > 0 {
			refType = p.ExternalRefs[0].ReferenceType
			refLoc = p.ExternalRefs[0].ReferenceLocator
		}

		packages = append(packages, Package{
			Name:             p.Name,
			Version:          p.VersionInfo,
			ReferenceType:    refType,
			ReferenceLocator: refLoc,
		})
	}

	return packages, nil
}
