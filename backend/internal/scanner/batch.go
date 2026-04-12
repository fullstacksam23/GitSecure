package scanner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// nonProjectPatterns are name/description signals that indicate a repo is
// a curated list, tutorial, or resource collection rather than a real project.
var nonProjectPatterns = []string{
	"awesome", "list", "book", "tutorial", "course",
	"interview", "cheatsheet", "resource", "guide", "roadmap",
	"algorithm", "leetcode", "learn", "study", "collection",
	"free-", "public-api", "howto", "example", "sample",
	"exercism", "kata", "challenge", "practice",
	"beginner", "snippet", "for-beginners", "seconds-of",
	"workshop", "curriculum", "lesson", "education", "bootcamp",
	"starter", "template", "boilerplate", "scaffold", "demo",
	"reference", "handbook", "wiki", "faq", "notes",
}
var nonProjectTopics = map[string]bool{
	"tutorial": true, "education": true, "learning": true,
	"course": true, "beginner": true, "awesome": true,
	"awesome-list": true, "list": true, "resources": true,
	"workshop": true, "curriculum": true, "snippets": true,
}

// dependencyFiles maps a language to the files that indicate a real project
// with trackable dependencies a security scanner can actually act on.
var dependencyFiles = map[string][]string{
	"python":     {"requirements.txt", "setup.py", "pyproject.toml", "Pipfile"},
	"javascript": {"package.json"},
	"typescript": {"package.json"},
	"go":         {"go.mod"},
	"java":       {"pom.xml", "build.gradle"},
	"ruby":       {"Gemfile"},
	"rust":       {"Cargo.toml"},
	"php":        {"composer.json"},
	"kotlin":     {"build.gradle", "pom.xml"},
	"swift":      {"Package.swift", "Podfile"},
}

// GetRepos returns repositories that are real, dependency-having projects
// suitable for security scanning — not curated lists or learning resources.
//
// Strategy:
//  1. Fetch a larger pool (up to GitHub's per_page max of 100) so we have
//     candidates to filter down from.
//  2. Reject repos whose name or description match known non-project patterns.
//  3. Verify each surviving candidate actually contains a dependency file for
//     the language by hitting the GitHub contents API.
//  4. Return the first `count` repos that pass all checks.
func GetRepos(language, githubToken string, count int) ([]Repo, error) {
	const fetchMultiplier = 5 // fetch 5x to have enough repos left after filtering
	fetchCount := min(count*fetchMultiplier, 100)

	candidates, err := searchRepos(language, githubToken, fetchCount)
	if err != nil {
		return nil, err
	}

	depFiles := dependencyFiles[strings.ToLower(language)]

	var results []Repo
	for _, repo := range candidates {
		if isNonProject(repo) {
			continue
		}
		if len(depFiles) > 0 {
			hasDeps, err := hasDependencyFile(repo, depFiles, githubToken)
			if err != nil || !hasDeps {
				continue
			}
		}
		results = append(results, repo)
		if len(results) == count {
			break
		}
	}

	return results, nil
}

func isNonProject(repo Repo) bool {
	haystack := strings.ToLower(repo.Name + " " + repo.Description)
	for _, pattern := range nonProjectPatterns {
		if strings.Contains(haystack, pattern) {
			return true
		}
	}
	// Topics are explicit author-set tags — much higher signal than text matching
	for _, topic := range repo.Topics {
		if nonProjectTopics[strings.ToLower(topic)] {
			return true
		}
	}
	return false
}

// hasDependencyFile checks whether any of the known dependency filenames exist
// at the root of the repo's default branch via the GitHub contents API.
func hasDependencyFile(repo Repo, files []string, token string) (bool, error) {

	for _, file := range files {
		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/contents/%s?ref=%s",
			repo.FullName,
			file,
			repo.DefaultBranch,
		)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return false, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")

		resp, err := Client.Do(req)
		if err != nil {
			return false, err
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return true, nil
		}
	}

	return false, nil
}

// searchRepos calls the GitHub search API and returns raw results.
func searchRepos(language, token string, count int) ([]Repo, error) {
	url := fmt.Sprintf(
		"https://api.github.com/search/repositories?q=language:%s+stars:>1000+fork:false&sort=stars&order=desc&per_page=%d",
		language,
		count,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var result BatchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Items, nil
}
