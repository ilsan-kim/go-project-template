package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigLoad(t *testing.T) {
	t.Run("load config", func(t *testing.T) {
		conf, err := Load("../config.json")
		assert.NoError(t, err)

		dbport := "3306"
		assert.Equal(t, dbport, conf.Db.DbPort)
	})
}
