package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

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

var osvUrl = "https://api.osv.dev/v1/querybatch"

func OSVScan(pkgs []Package) (OSVResponse, error) {
	var osvResp OSVResponse

	payload := createPayload(pkgs)

	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return osvResp, err
	}
	resp, err := httpClient.Post(osvUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return osvResp, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&osvResp)
	if err != nil {
		return osvResp, err
	}
	return osvResp, nil
}

func createPayload(pkgs []Package) OSVRequest {
	queries := make([]PackageQuery, len(pkgs))

	for i, pkg := range pkgs {
		queries[i] = PackageQuery{
			Package: Purl{
				Purl: pkg.ReferenceLocator,
			},
		}
	}
	return OSVRequest{Queries: queries}
}
