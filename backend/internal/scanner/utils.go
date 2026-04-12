package scanner

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getCurrentCommitHash(owner, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/HEAD", owner, repo)

	resp, err := http.Get(url)
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
