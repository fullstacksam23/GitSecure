package db

import (
	"errors"

	"github.com/fullstacksam23/GitSecure/internal/models"
)

func CreateBatchJob(batchJob models.BatchJob) error {
	if Client == nil {
		return errors.New("Client not initialized")
	}
	_, _, err := Client.From("ecosystem_scans").Insert(batchJob, false, "", "", "").Execute()
	if err != nil {
		return err
	}
	return nil
}
