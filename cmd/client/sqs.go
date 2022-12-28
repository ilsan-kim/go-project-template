package client

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"sampleProject/config"
)

var SQSClient *sqs.SQS
var QueueUrl string

func NewSQSClient(conf *config.Config) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	input := &sqs.GetQueueAttributesInput{
		AttributeNames: []*string{aws.String("CreatedTimestamp")},
		QueueUrl:       aws.String(conf.SqsQueue.QueueURL),
	}

	_, err := svc.GetQueueAttributes(input)
	if err != nil {
		return errors.New("can't set sqs Client (fail on ping-pong)")
	}

	SQSClient = svc
	QueueUrl = conf.SqsQueue.QueueURL
	return nil
}
