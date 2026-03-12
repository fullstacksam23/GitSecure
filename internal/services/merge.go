package services

import "github.com/fullstacksam23/GitSecure/internal/models"

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

		vulns[i].Source = "grype+osv"
	}

	return vulns
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
