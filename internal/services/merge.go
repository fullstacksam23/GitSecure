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

		if vulns[i].Severity == "" && len(adv.Severity) > 0 {
			vulns[i].Severity = adv.Severity[0].Score
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

func DeduplicateVulns(vulns []models.UnifiedVuln) []models.UnifiedVuln {
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
