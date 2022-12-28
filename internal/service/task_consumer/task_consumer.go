package task_consumer

import (
	"context"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	task2 "sampleProject/internal/service/task"
	"strconv"
	"strings"
	"time"
)

type ConsumerRepository interface {
	Delete(ctx context.Context, id int) error

	Begin(ctx context.Context) error
	Commit() error
	Rollback() error
}

type ConsumerQueue interface {
	ReceiveMessage() (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(receiptHandle string) error
}

type ConsumerService struct {
	repo  ConsumerRepository
	queue ConsumerQueue
}

func NewConsumerService(repo ConsumerRepository, queue ConsumerQueue) *ConsumerService {
	return &ConsumerService{
		repo:  repo,
		queue: queue,
	}
}

func (c *ConsumerService) DeleteTask() error {
	// first : receiveMessage
	receiveMessageOutput, err := c.queue.ReceiveMessage()
	if err != nil {
		return err
	}

	// hard delete process
	if len(receiveMessageOutput.Messages) == 0 {
		log.Println("no events to consume")
		time.Sleep(1 * time.Second)
		return nil
	}

	// retrieve job and target
	event := strings.Split(*receiveMessageOutput.Messages[0].Body, "|")
	receiptHandle := *receiveMessageOutput.Messages[0].ReceiptHandle
	job, task := event[0], event[1]

	// consume
	switch job {
	case "DELETE":
		ctx := context.WithValue(context.Background(), task2.Transactional{}, false)
		taskId, _ := strconv.Atoi(task)
		err = c.repo.Delete(ctx, taskId)
		if err != nil {
			return err
		}
		err = c.queue.DeleteMessage(receiptHandle)
		if err != nil {
			log.Println("successfully delete task from db but failed consume event")
			return err
		}

	default:
		log.Println("unknown job ", job)
		return nil
	}

	return nil
}
