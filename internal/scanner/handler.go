package scanner

import (
	"encoding/json"
	"net/http"

	"github.com/fullstacksam23/GitSecure/internal/db"
	"github.com/fullstacksam23/GitSecure/internal/models"
	"github.com/fullstacksam23/GitSecure/internal/redis"
	"github.com/google/uuid"
)

func StartScan(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Repo string `json:"repo"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	// get the context from the request
	ctx := r.Context()

	json.NewDecoder(r.Body).Decode(&req)

	job := models.ScanJob{
		JobID:  uuid.New().String(),
		Repo:   req.Repo,
		Status: "queued",
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
	}

	// return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	json.NewEncoder(w).Encode(map[string]string{
		"job_id": job.JobID,
		"repo":   job.Repo,
		"status": "queued",
	})
}

// type Repo struct {
// 	defaultBranch string `json:"default_branch"`
// }

// func checkCached(owner, repo string) error {
// 	repoFullName := owner + "/" + repo
// 	resp, err := http.Get("https://api.github.com/repos/" + repoFullName)
// 	if err != nil {
// 		return err
// 	}

// 	var r Repo
// 	err = json.NewDecoder(resp.Body).Decode(&r)
// 	if err != nil {
// 		return err
// 	}

// 	//get the default branch
// 	defaultBranch := r.defaultBranch

// 	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", owner, repo, defaultBranch)
// 	resp, err = http.Get(url)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
