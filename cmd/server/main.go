package main

import (
	"context"
	"log"
	"sampleProject/cmd/client"
	"sampleProject/config"
	"sampleProject/internal/common"
	"sampleProject/internal/pkg/queue"
	"sampleProject/internal/pkg/redis"
	"sampleProject/internal/pkg/sql/gorm"
	"sampleProject/internal/pkg/sql/sqlc"
	"sampleProject/internal/server"
	"sampleProject/internal/service/task"
	"time"
)

func main() {
	path := "config.json"

	var err error

	common.GracefulShutdownJob = make(chan bool, 1000)

	// init Config
	config.Conf, err = config.Load(path)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("conf : %v", config.Conf)

	// init gorm
	err = client.NewGormClient(config.Conf)
	if err != nil {
		log.Panicln(err)
	}
	// init db
	err = client.NewSQLClient(config.Conf)
	if err != nil {
		log.Panicln(err)
	}
	// init redis
	err = client.NewRedisClient(config.Conf)
	if err != nil {
		log.Panicln(err)
	}
	// init sqs
	err = client.NewSQSClient(config.Conf)
	if err != nil {
		log.Panicln(err)
	}

	// task Repo -> sql 구현체로 초기화..
	var taskRepo task.TaskRepository
	switch config.Conf.Db.DBClient {
	case "gorm":
		taskRepo = gorm.NewTask()
	case "sqlc":
		taskRepo = sqlc.NewTask()
	default:
		log.Println("db client not supported(or not selected)")
		return
	}
	taskRepo = sqlc.NewTask()
	//taskRepo = gorm.NewTask()

	// task cache -> redis 구현체로 초기화..
	var taskCache task.TaskCache
	taskCache = redis.NewTask()

	var taskQueue task.TaskQueue
	taskQueue, err = queue.NewQueue()
	if err != nil {
		log.Panicln(err)
	}

	// 서비스에 등록..
	taskService := task.NewTaskService(taskRepo, taskCache, taskQueue)
	services := &server.Service{
		Task: taskService,
	}

	srv := server.New(services)
	shutdown := func() {
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		srv.ShutdownGraceFully(ctx)
	}
	wait := common.RegisterSignal(shutdown)

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("listen: %v\n", err)
	}

	<-wait
	log.Println("Server has been gracefully shutdown")
}
