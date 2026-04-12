package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fullstacksam23/GitSecure/internal/db"
	"github.com/fullstacksam23/GitSecure/internal/models"
	"github.com/fullstacksam23/GitSecure/internal/redis"
	"github.com/google/uuid"
)

var q = redis.NewRedisQueue("localhost:6379", "scan_queue")

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
			BatchID:    "",
			Repo:       repoFullPath,
			Status:     "queued",
			CommitHash: hash,
			JobType:    "single",
		}

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

func BatchScan(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Language  string `json:"language"`
		RepoCount int    `json:"repo_count"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Language == "" || req.RepoCount <= 0 {
		http.Error(w, "both language and repo count should be specified for batch request", http.StatusBadRequest)
		return
	}
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		http.Error(w, "Github token not set", http.StatusInternalServerError)
		return
	}

	repos, err := GetRepos(req.Language, githubToken, req.RepoCount)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error getting batch repos from github: ", http.StatusInternalServerError)
		return
	}

	//create batch job
	batchJob := models.BatchJob{
		BatchID:        uuid.NewString(),
		Language:       req.Language,
		Status:         "queued",
		RepoCount:      req.RepoCount,
		CompletedRepos: 0,
		TotalRepos:     len(repos),
	}

	err = db.CreateBatchJob(batchJob)
	if err != nil {
		http.Error(w, "Error Inserting batch job into supabase", http.StatusInternalServerError)
		return
	}

	go ProcessBatchRepos(context.Background(), repos, batchJob.BatchID)

	// return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	json.NewEncoder(w).Encode(batchJob)
}

func ProcessBatchRepos(ctx context.Context, repos []Repo, batchID string) error {

	ecosystemRepos := make([]models.EcosystemRepo, len(repos))
	for i, curr := range repos {
		ecosystemRepos[i] = models.EcosystemRepo{
			BatchID:  batchID,
			RepoName: curr.FullName,
			Stars:    curr.Stars,
			RepoRank: i + 1,
		}
	}
	ids, err := db.CreateEcosystemRepos(ecosystemRepos)
	if err != nil {
		return err
	}
	log.Printf("IDs retrieved: %v", ids)

	for i, repoID := range ids {

		repoName := repos[i].FullName
		split := strings.Split(repoName, "/")
		owner := split[0]
		repo := split[1]
		hash, err := getCurrentCommitHash(owner, repo) //try to get this to 1 api call
		if err != nil {
			log.Printf("Error getting current commit hash %v", err)
			return err
		}

		job := models.ScanJob{
			JobID:      uuid.NewString(),
			BatchID:    batchID,
			Repo:       repoName,
			Status:     "queued",
			CommitHash: hash,
			RepoID:     int(repoID),
			JobType:    "batch",
		}

		err = q.Enqueue(ctx, job)
		if err != nil {
			log.Printf("Failed to enqueue batch job id: %s, err: %v", job.JobID, err)
			return err
		}

		err = db.InsertJob(job)
		if err != nil {
			log.Printf("Error Inserting job (batch) into supabase: %v", err)
			return err
		}
	}
	return nil
}
