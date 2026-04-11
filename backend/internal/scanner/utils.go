package scanner

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func getCurrentCommitHash(owner, repo string) (string, error) {
	repoFullName := owner + "/" + repo

	// Step 1: Get repo info to get default branch
	resp, err := http.Get("https://api.github.com/repos/" + repoFullName)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("unable to get repo info to check default branch - check private repo???")
	}

	var r Repo
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}

	// Step 2: Get latest commit from default branch
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/commits/%s",
		owner,
		repo,
		r.DefaultBranch,
	)

	resp, err = http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch commit: %s", resp.Status)
	}

	var commit CommitResponse
	if err := json.NewDecoder(resp.Body).Decode(&commit); err != nil {
		return "", err
	}

	return commit.SHA, nil
}
