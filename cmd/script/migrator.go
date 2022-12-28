package main

import (
	"log"
	"sampleProject/cmd/client"
	"sampleProject/config"
	"sampleProject/internal/pkg/sql/gorm"
)

func main() {
	path := "config.json"

	var err error

	// init Config
	config.Conf, err = config.Load(path)
	if err != nil {
		log.Println(err)
		return
	}

	// init db
	err = client.NewGormClient(config.Conf)
	if err != nil {
		log.Panicln(err)
	}

	// migrate
	db := client.GormClient
	err = db.AutoMigrate(&gorm.TaskDAO{})
	if err != nil {
		log.Panicln(err)
	}
}
