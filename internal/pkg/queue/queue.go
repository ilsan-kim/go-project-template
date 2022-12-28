package queue

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"sampleProject/cmd/client"
)

type SQSQueue struct {
	sqs *sqs.SQS
	url string
}

func NewQueue() (*SQSQueue, error) {
	return &SQSQueue{
		sqs: client.SQSClient,
		url: client.QueueUrl,
	}, nil
}

func (q *SQSQueue) SendMessage(obj string) error {
	params := &sqs.SendMessageInput{
		QueueUrl:    aws.String(q.url),
		MessageBody: aws.String(obj),
	}

	_, err := q.sqs.SendMessage(params)
	if err != nil {
		return err
	}
	return nil
}

func (q *SQSQueue) ReceiveMessage() (*sqs.ReceiveMessageOutput, error) {
	params := &sqs.ReceiveMessageInput{
		MaxNumberOfMessages: aws.Int64(1),
		QueueUrl:            aws.String(q.url),
	}
	out, err := q.sqs.ReceiveMessage(params)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (q *SQSQueue) DeleteMessage(receiptHandle string) error {
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(q.url),
		ReceiptHandle: aws.String(receiptHandle),
	}
	_, err := q.sqs.DeleteMessage(params)
	if err != nil {
		return err
	}
	return nil
}
