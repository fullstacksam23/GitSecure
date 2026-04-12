package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/fullstacksam23/GitSecure/internal/models"
)

var httpClient = &http.Client{}

func CreateBatchJob(batchJob models.BatchJob) error {
	if Client == nil {
		return errors.New("Client not initialized")
	}
	_, _, err := Client.From("ecosystem_batches").Insert(batchJob, false, "", "", "").Execute()
	if err != nil {
		return err
	}
	return nil
}

func IncrementBatchProgress(batchID string) error {
	url := os.Getenv("SUPABASE_URL") + "/rest/v1/rpc/increment_batch_progress"

	body := map[string]interface{}{
		"p_batch_id": batchID,
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("apikey", os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("rpc failed: %s", resp.Status)
	}

	return nil
}

func MarkBatchRunning(batchID string) error {
	if Client == nil {
		return errors.New("client not initialized")
	}

	updates := map[string]interface{}{
		"status": "running",
	}

	_, _, err := Client.From("ecosystem_batches").
		Update(updates, "", "").
		Eq("batch_id", batchID).
		Eq("status", "queued").
		Execute()

	return err
}

// CreateEcosystemRepo inserts multiple repos and returns a list of their auto-generated IDs
func CreateEcosystemRepos(repos []models.EcosystemRepo) ([]int64, error) {
	if Client == nil {
		return nil, errors.New("client not initialized")
	}

	// "representation" tells supabase to return newly created rows
	data, _, err := Client.From("ecosystem_repos").Insert(repos, false, "", "representation", "").Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to insert repos: %w", err)
	}

	var insertedRepos []models.EcosystemRepo
	if err := json.Unmarshal(data, &insertedRepos); err != nil {
		return nil, fmt.Errorf("failed to unmarshal returned repos: %w", err)
	}

	if len(insertedRepos) == 0 {
		return nil, errors.New("repos inserted but no data returned")
	}

	// Create a slice to hold all the newly generated BIGSERIAL IDs
	ids := make([]int64, 0, len(insertedRepos))
	for _, repo := range insertedRepos {
		ids = append(ids, repo.ID)
	}

	return ids, nil
}
