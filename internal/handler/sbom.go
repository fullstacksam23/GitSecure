package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fullstacksam23/GitSecure/internal/services"
)

func ExtractDependencies(w http.ResponseWriter, r *http.Request) {
	repoName := r.URL.Query().Get("repo")
	if repoName == "" {
		http.Error(w, "missing repo param", 400)
		return
	}

	// repoName = "appsecco/dvna"
	sbomURL := "https://api.github.com/repos/" + repoName + "/dependency-graph/sbom"

	pkgs, err := services.ExtractDependencies(sbomURL)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pkgs)
}
