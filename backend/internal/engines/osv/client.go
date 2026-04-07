package osv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fullstacksam23/GitSecure/internal/core"
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

var CHUNK_SIZE = 100

func FetchAllOSVAdvisories(pkgs []core.Package) (map[string]models.OSVAdvisory, error) {

	advisories := make(map[string]models.OSVAdvisory)
	seen := map[string]bool{}

	var wg sync.WaitGroup
	var mu sync.Mutex

	sem := make(chan struct{}, 7)

	resp, err := batchProcess(pkgs, CHUNK_SIZE)
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

func MergeOSVData(vulns []models.UnifiedVuln, advisories map[string]models.OSVAdvisory, canonical map[string]string) []models.UnifiedVuln {

	for i := range vulns {
		id := vulns[i].ID

		// Normalize advisory ID
		if canonicalID, ok := canonical[id]; ok {
			id = canonicalID
		}

		adv, ok := advisories[id]
		if !ok {
			continue
		}

		// Fill summary if missing
		if vulns[i].Summary == "" {
			vulns[i].Summary = adv.Summary
		}

		// Add OSV references
		var osvUrls []string
		for _, ref := range adv.References {
			osvUrls = append(osvUrls, ref.URL)
		}

		vulns[i].Urls = addUniqueUrls(vulns[i].Urls, osvUrls)

		if shouldUseDatabaseSeverity(vulns[i].Severity, adv.DatabaseSpecific.Severity) {
			vulns[i].Severity = adv.DatabaseSpecific.Severity
		}
		vulns[i].CWEIDs = adv.DatabaseSpecific.CWEIDs

		if len(adv.Affected) > 0 {
			vulns[i].Ecosystem = adv.Affected[0].Package.Ecosystem
		}

		if len(vulns[i].FixVersion) == 0 {
			for _, aff := range adv.Affected {
				for _, r := range aff.Ranges {
					for _, e := range r.Events {
						if e.Fixed != "" {
							vulns[i].FixVersion = append(vulns[i].FixVersion, e.Fixed)
						}
					}
				}
			}
		}

		vulns[i].Source = "grype+osv"
	}
	//remove duplicate vulns before returning
	return deduplicateVulns(vulns)
}

func addUniqueUrls(existing []string, newUrls []string) []string {

	seen := make(map[string]struct{})

	for _, u := range existing {
		seen[u] = struct{}{}
	}

	for _, u := range newUrls {
		if _, ok := seen[u]; !ok {
			existing = append(existing, u)
			seen[u] = struct{}{}
		}
	}

	return existing
}

func deduplicateVulns(vulns []models.UnifiedVuln) []models.UnifiedVuln {
	seen := make(map[string]models.UnifiedVuln)

	for _, v := range vulns {
		key := v.ID + "|" + v.Package + "|" + v.Version

		existing, ok := seen[key]
		if !ok {
			seen[key] = v
			continue
		}

		// Merge URLs
		existing.Urls = addUniqueUrls(existing.Urls, v.Urls)

		// Merge Fix Versions
		existing.FixVersion = append(existing.FixVersion, v.FixVersion...)

		seen[key] = existing
	}

	var result []models.UnifiedVuln
	for _, v := range seen {
		result = append(result, v)
	}

	return result
}

func shouldUseDatabaseSeverity(current, candidate string) bool {
	candidate = strings.TrimSpace(candidate)
	if candidate == "" {
		return false
	}

	current = strings.TrimSpace(current)
	if current == "" {
		return true
	}

	return strings.HasPrefix(strings.ToUpper(current), "CVSS:")
}

func getOSVAdvisory(id string) (models.OSVAdvisory, error) {
	var advisory models.OSVAdvisory
	if id == "" || len(id) == 0 {
		return advisory, nil
	}
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

type Result struct {
	Resp *models.OSVResponse
	Err  error
}

func processChunk(chunk []core.Package, ch chan Result) {

	osvUrl := "https://api.osv.dev/v1/querybatch"
	var osvResp models.OSVResponse

	payload := createPayload(chunk)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		ch <- Result{Resp: nil, Err: err}
		return
	}
	resp, err := osvHTTPClient.Post(osvUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		ch <- Result{Resp: nil, Err: err}
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- Result{Resp: nil, Err: fmt.Errorf("Batch api query failed statuscode: %d", resp.StatusCode)}
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&osvResp)
	if err != nil {
		ch <- Result{Resp: nil, Err: err}
		return
	}
	ch <- Result{Resp: &osvResp, Err: nil}
}

func batchProcess(pkgs []core.Package, chunkSize int) (models.OSVResponse, error) {
	numChunks := int(math.Ceil(float64(len(pkgs)) / float64(chunkSize)))
	var wg sync.WaitGroup

	ch := make(chan Result, numChunks)

	sem := make(chan struct{}, 10) // limit to 10 concurrent workers

	for i := 0; i < len(pkgs); i += chunkSize {
		j := i + chunkSize
		if j > len(pkgs) {
			j = len(pkgs)
		}

		chunk := pkgs[i:j]
		wg.Add(1)

		sem <- struct{}{} // acquire slot

		go func(chunk []core.Package) {
			defer wg.Done()
			defer func() { <-sem }() // release slot

			processChunk(chunk, ch)
		}(chunk)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	combinedResult := models.OSVResponse{}

	var firstErr error

	for res := range ch {
		if res.Err != nil && firstErr == nil {
			firstErr = res.Err
			continue
		}
		if res.Resp != nil {
			combinedResult.Results = append(combinedResult.Results, res.Resp.Results...)
		}
	}
	if firstErr != nil {
		return combinedResult, firstErr
	}
	return combinedResult, nil
}

func createPayload(pkgs []core.Package) models.OSVRequest {
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
