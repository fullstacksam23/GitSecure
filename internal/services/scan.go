package services

import (
	"log"
)

func RunFullScan(repo string) error {

	pkgs, sbom, err := getDependencies(repo)
	if err != nil {
		return err
	}
	log.Println("Dependencies Extracted: ")

	advisories, err := FetchAllOSVAdvisories(pkgs)
	if err != nil {
		return err
	}
	log.Println("OSV API queried...")

	graph := BuildVulnGraph(advisories)

	canonicalMap := graph.CanonicalMap()

	// Run grype
	raw, err := GrypeScan(sbom)
	if err != nil {
		log.Println("Grype error:", string(raw))
		return err
	}

	// Parse grype JSON
	grypeResp, err := ParseGrype(raw)
	if err != nil {
		return err
	}
	log.Println("Grype response generated...")
	// Normalize IDs
	vulns := NormalizeGrype(grypeResp, canonicalMap)

	for _, v := range vulns {
		log.Println(v)
	}

	//TODO: Store the pkgs in db (maybe supabase) and update job status
	//TODO: Also create handlers for db related functionality for user to get current status and job results
	return nil
}
