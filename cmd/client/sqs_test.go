package client

import (
	"github.com/stretchr/testify/assert"
	"sampleProject/config"
	"testing"
)

func TestNewSQSClient(t *testing.T) {
	conf, _ := config.Load("../../config.json")
	err := NewSQSClient(conf)
	assert.NoError(t, err)
}
