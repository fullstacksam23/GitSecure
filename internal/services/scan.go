package services

import (
	"log"
)

func RunFullScan(repo string) error {

	pkgs, err := getDependencies(repo)
	if err != nil {
		return err
	}
	log.Println("Dependencies Extracted: ")
	log.Println(pkgs)
	//TODO: Store the pkgs in db (maybe supabase) and update job status
	//TODO: Also create handlers for db related functionality for user to get current status and job results
	return nil
}
