package client

import (
	"github.com/stretchr/testify/assert"
	"sampleProject/config"
	"testing"
)

func TestNewRedisClient(t *testing.T) {
	conf, _ := config.Load("../../config.json")
	err := NewRedisClient(conf)
	assert.NoError(t, err)
}
