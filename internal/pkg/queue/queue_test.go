package queue

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sampleProject/cmd/client"
	"sampleProject/config"
	"testing"
)

func TestSQS(t *testing.T) {
	conf, _ := config.Load("../../../config.json")
	client.NewSQSClient(conf)
	q, _ := NewQueue()

	var receiptHandle string

	t.Run("send message", func(t *testing.T) {
		obj := "new task"
		err := q.SendMessage(obj)
		assert.NoError(t, err)
	})

	t.Run("receive message", func(t *testing.T) {
		out, err := q.ReceiveMessage()
		fmt.Println(out)
		assert.NoError(t, err)
		assert.Equal(t, len(out.Messages), 1)

		receiptHandle = *out.Messages[0].ReceiptHandle
	})

	t.Run("delete message", func(t *testing.T) {
		fmt.Println(receiptHandle)
		err := q.DeleteMessage(receiptHandle)
		assert.NoError(t, err)
	})
}
