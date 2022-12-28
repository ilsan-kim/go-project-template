package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"sampleProject/cmd/client"
	"sampleProject/internal/domain"
	"strconv"
	"time"
)

type Task struct {
	client     *redis.Client
	expiration time.Duration
}

func NewTask() *Task {
	return &Task{
		client:     client.RedisClient,
		expiration: 10 * time.Minute,
	}
}

func (t *Task) Set(ctx context.Context, data *domain.Task) error {
	key := strconv.Itoa(int(data.ID))
	bytes, _ := json.Marshal(*data)
	err := t.client.Set(ctx, key, bytes, t.expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) Get(ctx context.Context, id int) (*domain.Task, error) {
	res := &domain.Task{}

	key := strconv.Itoa(id)
	task, err := t.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(task, res); err != nil {
		return nil, err
	}

	return res, nil
}
