package worker

import (
	"log"
	"os"
	"runtime/pprof"
	"sampleProject/internal/common"
	"sampleProject/internal/service/task_consumer"
	"time"
)

type Service struct {
	Task *task_consumer.ConsumerService
}

type Consumer struct {
	service *Service

	IsRunning bool
}

func New(service *Service) *Consumer {
	consumer := new(Consumer)
	consumer.service = service

	return consumer
}

func (c *Consumer) Run(numWorkers int) {
	c.IsRunning = true
	workers := make(chan bool, numWorkers)
	done := make(chan bool, numWorkers)

	for c := 0; c < numWorkers; c++ {
		done <- true
	}

	for c.IsRunning {
		workers <- true
		<-done
		go func() {
			err := c.service.Task.DeleteTask()
			<-workers
			done <- true
			if err != nil {
				log.Println(err)
			}
		}()
	}
	return
}

func (c *Consumer) ShutdownGraceFully() {
	c.IsRunning = false
	for {
		cnt := len(common.GracefulShutdownJob)
		if cnt == 0 {
			break
		}
		log.Println("remained job : ", cnt)
		f, _ := os.OpenFile("jobs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		pprof.Lookup("goroutine").WriteTo(f, 1)
		f.Close()
		time.Sleep(time.Second * 5)
	}

	return
}
