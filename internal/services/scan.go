package services

import (
	"errors"
	"log"
	"net/http"
	"time"
)

func getDependencies(repoName string) ([]Package, error) {
	var pkgs []Package
	log.Println("trying to fetch sbom using github api...")
	if repoName == "" {
		return pkgs, errors.New("Repo name null/empty")
	}
	// repoName = "appsecco/dvna" for testing
	sbomURL := "https://api.github.com/repos/" + repoName + "/dependency-graph/sbom"
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Get(sbomURL)

	if err != nil {
		return pkgs, err
	}

	defer resp.Body.Close()

	//sbom not available in this case
	if resp.StatusCode == 404 {
		log.Println("SBOM not available... parsing manually")
		pkgs, err = ExtractDependenciesManual(repoName)
		if err != nil {
			return nil, err
		}
	} else {
		pkgs, err = ExtractDependencies(resp.Body)
		if err != nil {
			return nil, err
		}
	}
	return pkgs, nil

}

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
