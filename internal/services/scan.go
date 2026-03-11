package services

import (
	"log"
)

func RunFullScan(repo string) error {

	pkgs, sbom, err := getDependencies(repo)
	log.Println(string(sbom))
	if err != nil {
		return err
	}
	log.Println("Dependencies Extracted: ")

	resp, err := OSVScan(pkgs)
	if err != nil {
		return err
	}
	log.Println("OSV api queried...")
	log.Println(resp)

	output, err := GrypeScan(sbom)
	if err != nil {
		log.Println("Grype error:", string(output))
		return err
	}
	log.Println("Grype scan completed...")
	log.Println(string(output))
	//TODO: Store the pkgs in db (maybe supabase) and update job status
	//TODO: Also create handlers for db related functionality for user to get current status and job results
	return nil
}
