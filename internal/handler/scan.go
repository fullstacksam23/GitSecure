package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fullstacksam23/GitSecure/internal/db"
	"github.com/fullstacksam23/GitSecure/internal/models"
	"github.com/fullstacksam23/GitSecure/internal/worker"
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

	q := worker.NewRedisQueue("localhost:6379", "scan_queue")

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
