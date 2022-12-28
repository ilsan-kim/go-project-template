package client

import (
	"github.com/stretchr/testify/assert"
	"sampleProject/config"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestNewSQLClient(t *testing.T) {
	conf, _ := config.Load("../../config.json")
	err := NewSQLClient(conf)
	assert.NoError(t, err)
}

func TestNewGormClient(t *testing.T) {
	conf, _ := config.Load("../../config.json")
	err := NewGormClient(conf)
	assert.NoError(t, err)
}
