package scanner

import (
	"context"
	"log"

	"github.com/fullstacksam23/GitSecure/internal/db"
	"github.com/fullstacksam23/GitSecure/internal/engines/grype"
	"github.com/fullstacksam23/GitSecure/internal/engines/osv"
	"github.com/fullstacksam23/GitSecure/internal/engines/sbom"
)

func RunFullScan(ctx context.Context, repo, jobID string) error {

	log.Println(repo, jobID)

	pkgs, sbom, err := sbom.GetDependencies(repo)
	if err != nil {
		return err
	}
	log.Println("Dependencies Extracted: ")

	advisories, err := osv.FetchAllOSVAdvisories(pkgs)

	if err != nil {
		return err
	}
	log.Println("OSV API queried...")

	graph := osv.BuildVulnGraph(advisories)

	canonicalMap := graph.CanonicalMap()

	advisories = osv.CanonicalizeAdvisories(advisories, canonicalMap)

	log.Println("Setting advisory id with right priority...")

	raw, err := grype.GrypeScan(sbom)
	if err != nil {
		return err
	}

	// Parse grype JSON
	grypeResp, err := grype.ParseGrype(raw)
	if err != nil {
		return err
	}
	log.Println("Grype response generated...")
	// Normalize IDs
	vulns := grype.NormalizeGrype(grypeResp, canonicalMap, jobID)

	//Enrich grype with OSV data
	vulns = osv.MergeOSVData(vulns, advisories, canonicalMap)

	log.Println("Vulns list generated")

	err = db.InsertVulns(ctx, vulns)
	if err != nil {
		log.Println("Error Inserting vulns into supabase")
		return err
	}
	log.Println("Supabase Updated")

	//TODO: Also create handlers for db related functionality for user to get current status and job results
	return nil
}
