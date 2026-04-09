package db

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/fullstacksam23/GitSecure/internal/models"
	"github.com/supabase-community/postgrest-go"
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

	job.Status = NormalizeJobStatus(job.Status)

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

	if status, ok := updates["status"].(string); ok {
		updates["status"] = NormalizeJobStatus(status)
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

func CheckExistingJob(repoFullPath string) (bool, *models.ScanJob, error) {
	if Client == nil {
		return false, nil, errors.New("client not initialized")
	}

	data, _, err := Client.
		From("scan_jobs").
		Select("*", "", false).
		Eq("repo", repoFullPath).
		In("status", []string{"complete", "completed"}).
		Order("created_at", &postgrest.OrderOpts{
			Ascending: false, // DESC
		}).
		Limit(1, "").
		Execute()

	if err != nil {
		return false, nil, err
	}
	var result []models.ScanJob

	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Println("json unmarshalling error", err)
		return false, nil, err
	}

	if len(result) == 0 {
		log.Println("data not cached")
		return false, nil, nil
	}

	result[0].Status = NormalizeJobStatus(result[0].Status)

	return true, &result[0], nil
}
