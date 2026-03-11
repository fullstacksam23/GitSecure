package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fullstacksam23/GitSecure/internal/models"
)

type Package struct {
	Name             string
	Version          string
	ReferenceType    string
	ReferenceLocator string
}

func getDependencies(repoName string) ([]Package, []byte, error) {
	var pkgs []Package
	log.Println("trying to fetch sbom using github api...")
	if repoName == "" {
		return pkgs, nil, errors.New("Repo name null/empty")
	}
	// repoName = "appsecco/dvna" for testing
	sbomURL := "https://api.github.com/repos/" + repoName + "/dependency-graph/sbom"
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Get(sbomURL)

	if err != nil {
		return pkgs, nil, err
	}

	defer resp.Body.Close()
	var body []byte
	//sbom not available in this case
	if resp.StatusCode == 404 {
		log.Println("SBOM not available... parsing manually")
		pkgs, err = ExtractDependenciesManual(repoName)
		if err != nil {
			return nil, nil, err
		}
	} else {

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, err
		}

		pkgs, err = ExtractDependencies(bytes.NewReader(body))
		if err != nil {
			return nil, nil, err
		}
	}
	return pkgs, body, nil

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
			//set refType and refLoc to the first values as default
			refType = p.ExternalRefs[0].ReferenceType
			refLoc = p.ExternalRefs[0].ReferenceLocator
			// specifically look for the purl references
			for i := 0; i < len(p.ExternalRefs); i++ {
				if p.ExternalRefs[i].ReferenceType == "purl" {
					refType = p.ExternalRefs[i].ReferenceType
					refLoc = p.ExternalRefs[i].ReferenceLocator
					break
				}
			}
		}
		decodedRefLoc, err := url.QueryUnescape(refLoc)
		if err != nil {
			return nil, err
		}
		if decodedRefLoc == "" || strings.HasPrefix(decodedRefLoc, "pkg:github") {
			continue
		}
		packages = append(packages, Package{
			Name:             p.Name,
			Version:          p.VersionInfo,
			ReferenceType:    refType,
			ReferenceLocator: decodedRefLoc,
		})
	}

	return packages, nil
}
