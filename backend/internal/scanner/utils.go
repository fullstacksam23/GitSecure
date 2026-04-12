package scanner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var Client = &http.Client{}

func getCurrentCommitHash(owner, repo string) (string, error) {

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/HEAD", owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	token := os.Getenv("GITHUB_TOKEN")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := Client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch commit: %s", resp.Status)
	}

	var commit struct {
		SHA string `json:"sha"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&commit); err != nil {
		return "", err
	}

	return commit.SHA, nil
}

func GetRepoDetails(owner, repo, token string) (Repo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Repo{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := Client.Do(req)
	if err != nil {
		return Repo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Repo{}, fmt.Errorf("failed to fetch repo: %s", resp.Status)
	}

	var r Repo
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return Repo{}, err
	}

	return r, nil
}
