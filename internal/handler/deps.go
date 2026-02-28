package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fullstacksam23/GitSecure/internal/services"
)

func ExtractDependencies(w http.ResponseWriter, r *http.Request) {
	repoName := r.URL.Query().Get("repo")
	if repoName == "" {
		http.Error(w, "missing repo param", 400)
		return
	}

	// repoName = "appsecco/dvna" for testing
	sbomURL := "https://api.github.com/repos/" + repoName + "/dependency-graph/sbom"
	resp, err := http.Get(sbomURL)

	if err != nil {
		http.Error(w, "Error occurred while trying to get dependencies", 500)
		return
	}

	defer resp.Body.Close()

	var pkgs []services.Package

	//sbom not available in this case
	if resp.StatusCode == 404 {
		fmt.Println("SBOM not available... parsing manually")
		pkgs, err = services.ExtractDependenciesManual(repoName)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	} else {
		pkgs, err = services.ExtractDependencies(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pkgs)
}
