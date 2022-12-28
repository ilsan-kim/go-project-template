package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var Conf *Config

type Config struct {
	Db struct {
		DbPort     string `json:"db_port"`
		DbURL      string `json:"db_url"`
		DbPassword string `json:"db_password"`
		DbUsername string `json:"db_username"`
		DbSchema   string `json:"db_schema"`
		DBClient   string `json:"db_client"`
	} `json:"db"`
	RedisCache struct {
		RedisPort string `json:"redis_port"`
		RedisURL  string `json:"redis_url"`
		RedisDb   string `json:"redis_db"`
	} `json:"redis"`
	S3 struct {
		S3URL        string `json:"s3_url"`
		S3BucketName string `json:"s3_bucket_name"`
	} `json:"s3"`
	SqsQueue struct {
		QueueURL string `json:"queue_url"`
	} `json:"sqs_queue"`
}

func Load(path string) (conf *Config, err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		return
	}

	return
}
