package worker

import (
	"context"
	"log"

	"github.com/fullstacksam23/GitSecure/internal/queue"
	"github.com/fullstacksam23/GitSecure/internal/services"
)

func StartWorker(ctx context.Context) {

	q := queue.NewRedisQueue("localhost:6379", "scan_queue")

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
			log.Println("Processing job:", job.JobID)

			err = services.RunFullScan(job.Repo)
			if err != nil {
				log.Println("Scan failed:", err)
				continue
			}

			log.Println("Scan complete:", job.JobID)
		}
	}
}
