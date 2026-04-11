package worker

import (
	"context"
	"log"

	"github.com/fullstacksam23/GitSecure/internal/db"
	"github.com/fullstacksam23/GitSecure/internal/redis"
	"github.com/fullstacksam23/GitSecure/internal/scanner"
)

func StartWorker(ctx context.Context, githubToken string) {

	q := redis.NewRedisQueue("localhost:6379", "scan_queue")

	log.Println("Worker started...")

	for {

		select {

		case <-ctx.Done():
			log.Println("Worker shutting down...")
			return

		default:

			job, err := q.Dequeue(ctx)
			if err != nil {
				if ctx.Err() != nil {
					log.Println("Worker stopped")
					return
				}

				log.Println("Queue error:", err)
				continue
			}
			if job == nil {
				continue
			}

			err = db.UpdateJobStatus(job.JobID, map[string]interface{}{
				"status": "running",
			})
			if err != nil {
				log.Println("Queue error:", err)
				continue
			}
			log.Println("Processing job:", job.JobID)

			err = scanner.RunFullScan(ctx, job.Repo, job.JobID, githubToken)
			if err != nil {
				log.Println("Scan failed:", err)
				updateErr := db.UpdateJobStatus(job.JobID, map[string]interface{}{
					"status": "failed",
				})
				if updateErr != nil {
					log.Println("Queue error:", updateErr)
				}
				continue
			}

			log.Println("Scan complete:", job.JobID)

			err = db.UpdateJobStatus(job.JobID, map[string]interface{}{
				"status": "completed",
			})
			if err != nil {
				log.Println("Queue error:", err)
				continue
			}
			//TODO: ERROR HANDLING AND PASS CONTEXT TO THE DB UPDATE AND INTSERT FUNCTIONS
		}
	}
}
