package db

import (
	"context"
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

func InsertJob(job interface{}) error {
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
