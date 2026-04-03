package scanner

import (
	"context"
	"log"
	"strconv"

	"github.com/fullstacksam23/GitSecure/internal/db"
	"github.com/fullstacksam23/GitSecure/internal/engines/grype"
	"github.com/fullstacksam23/GitSecure/internal/engines/osv"
	"github.com/fullstacksam23/GitSecure/internal/engines/sbom"
)

func RunFullScan(ctx context.Context, repo, jobID string) error {

	pkgs, sbom, err := sbom.GetDependencies(repo)
	if err != nil {
		return err
	}
	log.Println("Dependencies Extracted: ")

	advisories, err := osv.FetchAllOSVAdvisories(pkgs)
	// log.Println(advisories)
	if err != nil {
		return err
	}
	log.Println("OSV API queried...")

	graph := osv.BuildVulnGraph(advisories)

	canonicalMap := graph.CanonicalMap()

	advisories = osv.CanonicalizeAdvisories(advisories, canonicalMap)

	log.Println("Setting advisory id with right priority...")

	// Run grype
	raw, err := grype.GrypeScan(sbom)
	if err != nil {
		log.Println("Grype error:", string(raw))
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

	for _, v := range vulns {
		log.Println(v.ID, v.Package, v.Version, v.Severity, v.Source, v.JobID)
		log.Println("URLS count:" + strconv.Itoa(len(v.Urls)))
		log.Println("Fix count: " + strconv.Itoa(len(v.FixVersion)))
	}

	//TODO: Also create handlers for db related functionality for user to get current status and job results
	return nil
}
