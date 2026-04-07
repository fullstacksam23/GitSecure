package scanner

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/fullstacksam23/GitSecure/internal/db"
	"github.com/fullstacksam23/GitSecure/internal/models"
	"github.com/fullstacksam23/GitSecure/internal/redis"
	"github.com/google/uuid"
)

func StartScan(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Owner string `json:"owner"`
		Repo  string `json:"repo"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	// get the context from the request
	ctx := r.Context()

	owner := req.Owner
	repo := req.Repo

	log.Println(owner, repo)
	if owner == "" || repo == "" {
		http.Error(w, "both owner and repo should be specfied", http.StatusBadRequest)
		return
	}
	repoFullPath := fmt.Sprintf("%s/%s", owner, repo)

	exists, existingJob, err := db.CheckExistingJob(repoFullPath)
	if err != nil {
		http.Error(w, "Error checking if job cached or not", http.StatusInternalServerError)
		return
	}

	var job models.ScanJob
	hash, err := getCurrentCommitHash(owner, repo)
	if err != nil {
		log.Println("error getting current commit hash", err)
		http.Error(w, "Error while getting repo details - check private repo???", http.StatusInternalServerError)
		return
	}

	if exists && hash == existingJob.CommitHash {
		job = *existingJob
	} else {

		job = models.ScanJob{
			JobID:      uuid.New().String(),
			Repo:       repoFullPath,
			Status:     "queued",
			CommitHash: hash,
		}

		q := redis.NewRedisQueue("localhost:6379", "scan_queue")

		err = q.Enqueue(ctx, job)
		if err != nil {
			http.Error(w, "failed to enqueue job", http.StatusInternalServerError)
			return
		}

		err = db.InsertJob(job)
		if err != nil {
			http.Error(w, "Error updating supabase job status", http.StatusInternalServerError)
			return
		}
	}

	// return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	json.NewEncoder(w).Encode(job)
}

type Repo struct {
	DefaultBranch string `json:"default_branch"`
}

type CommitResponse struct {
	SHA string `json:"sha"`
}

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
