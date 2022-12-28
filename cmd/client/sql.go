package client

import (
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sampleProject/config"

	_ "github.com/go-sql-driver/mysql"
)

var SQLClient *sql.DB
var GormClient *gorm.DB

func NewSQLClient(conf *config.Config) error {
	connectString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", conf.Db.DbUsername, conf.Db.DbPassword, conf.Db.DbURL, conf.Db.DbSchema)

	db, err := sql.Open("mysql", connectString)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return errors.New("can't set sql Client (fail on ping-pong)")
	}

	SQLClient = db
	return nil
}

func NewGormClient(conf *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", conf.Db.DbUsername, conf.Db.DbPassword, conf.Db.DbURL, conf.Db.DbPort, conf.Db.DbSchema)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return err
	}

	err = sqlDb.Ping()
	if err != nil {
		return errors.New("cant set gorm Client (fail on ping-pong")
	}

	GormClient = db
	return nil
}
