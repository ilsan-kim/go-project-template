package main

import (
	"log"
	"sampleProject/cmd/client"
	"sampleProject/config"
	"sampleProject/internal/common"
	"sampleProject/internal/pkg/queue"
	"sampleProject/internal/pkg/sql/gorm"
	"sampleProject/internal/pkg/sql/sqlc"
	"sampleProject/internal/service/task_consumer"
	"sampleProject/internal/worker"
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

	log.Printf("conf : %s\n", config.Conf)

	// init db
	err = client.NewSQLClient(config.Conf)
	if err != nil {
		log.Panicln(err)
	}
	err = client.NewGormClient(config.Conf)
	if err != nil {
		log.Panicln(err)
	}

	// init sqs
	err = client.NewSQSClient(config.Conf)
	if err != nil {
		log.Panicln(err)
	}

	var taskRepo task_consumer.ConsumerRepository
	switch config.Conf.Db.DBClient {
	case "gorm":
		taskRepo = gorm.NewTask()
	case "sqlc":
		taskRepo = sqlc.NewTask()
	default:
		log.Println("db client not supported(or not selected)")
		return
	}
	//taskRepo = sqlc.NewTask()
	taskRepo = gorm.NewTask()

	var taskQueue task_consumer.ConsumerQueue
	taskQueue, err = queue.NewQueue()
	if err != nil {
		log.Panicln(err)
	}

	taskConsumerService := task_consumer.NewConsumerService(taskRepo, taskQueue)
	services := &worker.Service{Task: taskConsumerService}
	worker := worker.New(services)

	shutdown := func() {
		worker.ShutdownGraceFully()
	}
	wait := common.RegisterSignal(shutdown)

	worker.Run(2)

	<-wait
	log.Println("Worker has been gracefully shutdown")
}
