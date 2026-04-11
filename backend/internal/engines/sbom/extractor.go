package sbom

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fullstacksam23/GitSecure/internal/core"
	"github.com/fullstacksam23/GitSecure/internal/models"
)

func GetDependencies(repoName, githubToken string) ([]core.Package, []byte, error) {
	var pkgs []core.Package
	log.Println("trying to fetch sbom using github api...")
	if repoName == "" {
		return pkgs, nil, errors.New("Repo name null/empty")
	}
	// repoName = "appsecco/dvna" for testing
	sbomURL := "https://api.github.com/repos/" + repoName + "/dependency-graph/sbom"
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, _ := http.NewRequest("GET", sbomURL, nil)

	req.Header.Set("Authorization", "Bearer "+githubToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		return pkgs, nil, err
	}

	defer resp.Body.Close()

	//sbom not available in this case
	if resp.StatusCode == 404 {
		log.Println("SBOM not available... parsing manually")
		pkgs, sbom, err := ExtractDependenciesManual(repoName)
		if err != nil {
			return nil, nil, err
		}

		return pkgs, sbom, nil
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("GitHub error:", resp.StatusCode, string(body))
		return nil, nil, fmt.Errorf("github api failed")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	cleanSbom, err := parseSbom(body)
	if err != nil {
		return nil, nil, err
	}

	pkgs, err = ExtractDependencies(cleanSbom)
	if err != nil {
		return nil, nil, err
	}

	return pkgs, cleanSbom, nil

}

func parseSbom(data []byte) ([]byte, error) {

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	// github api gives an extra sbom field
	if sbomRaw, ok := raw["sbom"]; ok {
		sbomBytes, err := json.Marshal(sbomRaw)
		if err != nil {
			return nil, err
		}
		return sbomBytes, nil //return with the extra {sbom: {}} removed
	}

	// Already SPDX returned by syft validate structure
	if _, ok := raw["spdxVersion"]; ok {
		return data, nil
	}

	return nil, errors.New("unsupported SBOM format")
}

func ExtractDependencies(sbom []byte) ([]core.Package, error) {

	var doc models.SPDXDocument
	err := json.Unmarshal(sbom, &doc)
	if err != nil {
		return nil, err
	}

	packages := []core.Package{}
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
		packages = append(packages, core.Package{
			Name:             p.Name,
			Version:          p.VersionInfo,
			ReferenceType:    refType,
			ReferenceLocator: decodedRefLoc,
		})
	}

	return packages, nil
}
