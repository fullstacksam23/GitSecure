package services

import (
	"context"
	"log"
	"strconv"

	"github.com/fullstacksam23/GitSecure/internal/db"
)

func RunFullScan(ctx context.Context, repo, jobID string) error {

	pkgs, sbom, err := getDependencies(repo)
	if err != nil {
		return err
	}
	log.Println("Dependencies Extracted: ")

	advisories, err := FetchAllOSVAdvisories(pkgs)
	// log.Println(advisories)
	if err != nil {
		return err
	}
	log.Println("OSV API queried...")

	graph := BuildVulnGraph(advisories)

	canonicalMap := graph.CanonicalMap()

	advisories = CanonicalizeAdvisories(advisories, canonicalMap)

	log.Println("Setting advisory id with right priority...")

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
	vulns := NormalizeGrype(grypeResp, canonicalMap, jobID)

	//Enrich grype with OSV data
	vulns = DeduplicateVulns(MergeOSVData(vulns, advisories, canonicalMap))

	log.Println("Vulns list generated")

	err = db.InsertVulns(ctx, vulns)
	if err != nil {
		log.Println("Error Inserting vulns into supabase")
		return err
	}
	log.Println("Supabase Updated")

	for _, v := range vulns {
		log.Println(v.ID, v.Package, v.Version, v.Severity, v.Source, v.JobID)
		log.Println("URLS count:" + strconv.Itoa(len(v.Urls)))
		log.Println("Fix count: " + strconv.Itoa(len(v.FixVersion)))
	}

	//TODO: Also create handlers for db related functionality for user to get current status and job results
	return nil
}
