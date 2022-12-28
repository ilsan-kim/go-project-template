package client

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sampleProject/config"
	"strconv"
)

var RedisClient *redis.Client

func NewRedisClient(conf *config.Config) error {
	dbNum, _ := strconv.Atoi(conf.RedisCache.RedisDb)

	connectAddr := fmt.Sprintf("%s:%s", conf.RedisCache.RedisURL, conf.RedisCache.RedisPort)

	c := redis.NewClient(&redis.Options{
		Addr: connectAddr,
		DB:   dbNum,
	})

	if _, err := c.Ping(context.Background()).Result(); err != nil {
		return err
	}

	RedisClient = c
	return nil
}
