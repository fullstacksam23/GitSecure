package db

import (
	"context"
	"errors"
	"log"

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

func InsertVulns(ctx context.Context, vulns interface{}) error {
	if Client == nil {
		return errors.New("client not initialized")
	}
	_, _, err := Client.
		From("vulnerabilities").
		Insert(vulns, false, "", "", "").
		Execute()

	return err
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
