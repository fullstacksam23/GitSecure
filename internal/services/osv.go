package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/fullstacksam23/GitSecure/internal/models"
)

var osvHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
	},
}

func FetchAllOSVAdvisories(pkgs []Package) (map[string]models.OSVAdvisory, error) {

	advisories := make(map[string]models.OSVAdvisory)
	seen := map[string]bool{}

	var wg sync.WaitGroup
	var mu sync.Mutex

	sem := make(chan struct{}, 10)

	resp, err := batchProcess(pkgs)
	if err != nil {
		return nil, err
	}

	for _, result := range resp.Results {
		for _, vuln := range result.Vulns {

			if seen[vuln.Id] {
				continue
			}
			seen[vuln.Id] = true

			wg.Add(1)

			go func(id string) {
				defer wg.Done()

				sem <- struct{}{}
				defer func() { <-sem }()

				adv, err := getOSVAdvisory(id)
				if err != nil {
					fmt.Println("OSV advisory fetch failed:", id, err)
					return
				}

				mu.Lock()
				advisories[id] = adv
				mu.Unlock()

			}(vuln.Id)
		}
	}

	wg.Wait()

	return advisories, nil
}

func getOSVAdvisory(id string) (models.OSVAdvisory, error) {
	var advisory models.OSVAdvisory
	url := "https://api.osv.dev/v1/vulns/" + id

	resp, err := osvHTTPClient.Get(url)
	if err != nil {
		return advisory, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return advisory, fmt.Errorf("OSV API returned status %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&advisory)
	return advisory, err
}

func batchProcess(pkgs []Package) (models.OSVResponse, error) {
	osvUrl := "https://api.osv.dev/v1/querybatch"
	var osvResp models.OSVResponse

	payload := createPayload(pkgs)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return osvResp, err
	}
	resp, err := osvHTTPClient.Post(osvUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return osvResp, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return osvResp, fmt.Errorf("OSV batch query failed: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&osvResp)
	if err != nil {
		return osvResp, err
	}
	return osvResp, nil
}

func createPayload(pkgs []Package) models.OSVRequest {
	queries := make([]models.PackageQuery, len(pkgs))

	for i, pkg := range pkgs {
		queries[i] = models.PackageQuery{
			Package: models.Purl{
				Purl: pkg.ReferenceLocator,
			},
		}
	}
	return models.OSVRequest{Queries: queries}
}
