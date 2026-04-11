package scanner

import (
	"context"
	"log"

	"github.com/fullstacksam23/GitSecure/internal/db"
	"github.com/fullstacksam23/GitSecure/internal/engines/grype"
	"github.com/fullstacksam23/GitSecure/internal/engines/osv"
	"github.com/fullstacksam23/GitSecure/internal/engines/sbom"
	"github.com/fullstacksam23/GitSecure/internal/models"
	"golang.org/x/sync/errgroup"
)

func RunFullScan(ctx context.Context, repo, jobID, githubToken string) error {

	log.Println(repo, jobID)

	pkgs, sbom, err := sbom.GetDependencies(repo, githubToken)
	if err != nil {
		return err
	}
	log.Println("Dependencies Extracted: ")

	g, ctx := errgroup.WithContext(ctx)

	//variables used in the goroutines
	var (
		advisories   map[string]models.OSVAdvisory
		canonicalMap map[string]string
		grypeResp    grype.GrypeResponse
	)

	//Run OSV Pipeline
	g.Go(func() error {
		advs, err := osv.FetchAllOSVAdvisories(pkgs)

		if err != nil {
			return err
		}
		log.Println("OSV API queried...")

		graph := osv.BuildVulnGraph(advs)

		cMap := graph.CanonicalMap()
		canonicalMap = cMap

		advs = osv.CanonicalizeAdvisories(advs, cMap)
		advisories = advs

		log.Println("Setting advisory id with right priority...")
		return nil
	})

	//Run Grype Pipeline
	g.Go(func() error {
		raw, err := grype.GrypeScan(sbom)
		if err != nil {
			return err
		}

		gResp, err := grype.ParseGrype(raw)
		if err != nil {
			return err
		}
		grypeResp = gResp

		log.Println("Grype response generated...")

		return nil
	})

	//Wait for the 2 parallel go routines to finish and log if any errors
	if err := g.Wait(); err != nil {
		log.Println("Pipeline Failed: ", err)
		return err
	}
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

	return nil
}

// func RunBatchScan(ctx context.Context, repo, batchID, githubToken string) error {
// 	log.Println(repo, jobID)

// 	pkgs, sbom, err := sbom.GetDependencies(repo, githubToken)
// 	if err != nil {
// 		return err
// 	}
// 	log.Println("Dependencies Extracted: ")
// }
