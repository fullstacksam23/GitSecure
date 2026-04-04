package db

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/fullstacksam23/GitSecure/internal/models"
	"github.com/supabase-community/supabase-go"
)

var Client *supabase.Client

func InitSupabase(url, key string) error {

	client, err := supabase.NewClient(url, key, nil)
	if err != nil {
		return err
	}

	Client = client
	log.Println("Supabase client initialized...")
	return nil
}

func InsertVulns(ctx context.Context, vulns []models.UnifiedVuln) error {
	if Client == nil {
		return errors.New("Client not initialized")
	}
	if len(vulns) == 0 {
		log.Println("vulns [] length is zero... No vulns found")
		return nil
	}
	_, _, err := Client.From("vulnerabilities").Insert(vulns, false, "", "", "").Execute()
	if err != nil {
		return err
	}

	return nil
}

func InsertJob(job models.ScanJob) error {
	if Client == nil {
		return errors.New("client not initialized")
	}

	_, _, err := Client.From("scan_jobs").Insert(job, false, "", "", "").Execute()
	if err != nil {
		return err
	}
	return nil
}
func UpdateJobStatus(jobID string, updates map[string]interface{}) error {
	if Client == nil {
		return errors.New("client not initialized")
	}

	_, _, err := Client.
		From("scan_jobs").
		Update(updates, "", "").
		Eq("job_id", jobID).
		Execute()
	if err != nil {
		return err
	}
	return nil
}

func CheckExistingJob(repo, commitHash string) (*models.ScanJob, []models.UnifiedVuln, error) {
	if Client == nil {
		return nil, nil, errors.New("client not initialized")
	}

	// 1️⃣ Get job
	data, _, err := Client.
		From("scan_jobs").
		Select("*", "", false).
		Eq("repo", repo).
		Eq("commit_hash", commitHash).
		Limit(1, "").
		Execute()

	if err != nil {
		return nil, nil, err
	}

	var jobs []models.ScanJob
	if err := json.Unmarshal(data, &jobs); err != nil {
		return nil, nil, err
	}

	if len(jobs) == 0 {
		// no cache
		return nil, nil, nil
	}

	job := jobs[0]

	// 2️⃣ Get vulnerabilities using job_id
	vulnData, _, err := Client.
		From("vulnerabilities").
		Select("*", "", false).
		Eq("job_id", job.JobID).
		Execute()

	if err != nil {
		return &job, nil, err
	}

	var vulns []models.UnifiedVuln
	if err := json.Unmarshal(vulnData, &vulns); err != nil {
		return &job, nil, err
	}

	return &job, vulns, nil
}
