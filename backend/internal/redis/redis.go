package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fullstacksam23/GitSecure/internal/models"
	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	Client    *redis.Client
	QueueName string
}

// Pass configuration to the constructor
func NewRedisQueue(redisAddr, queueName string) *RedisQueue {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	return &RedisQueue{
		Client:    rdb,
		QueueName: queueName,
	}
}

func (q *RedisQueue) Enqueue(ctx context.Context, job models.ScanJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return q.Client.RPush(ctx, q.QueueName, data).Err()
}

func (q *RedisQueue) Dequeue(ctx context.Context) (*models.ScanJob, error) {

	res, err := q.Client.BLPop(ctx, 5*time.Second, q.QueueName).Result()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var job models.ScanJob
	err = json.Unmarshal([]byte(res[1]), &job)

	return &job, err
}
