package main

import (
	"database/sql"
	"fmt"
	"http-service/cmd/log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBLog struct {
	log log.Logger
}

func (d DBLog) Printf(format string, data ...interface{}) {
	s := fmt.Sprintf(format, data...)
	d.log.Info(s)
}

func NewDB(conf *Config, log log.Logger) (*gorm.DB, error) {
	db, err := sql.Open("mysql", conf.Database.DatabaseUrl)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Second * time.Duration(conf.Database.ConnMaxLifetime))
	db.SetMaxOpenConns(conf.Database.MaxOpenConns)
	db.SetMaxIdleConns(conf.Database.MaxIdleConns)
	gormDB, err := gorm.Open(
		mysql.New(mysql.Config{Conn: db}),
		&gorm.Config{
			Logger: logger.New(
				DBLog{log: log},
				logger.Config{
					LogLevel: logger.Info,
				},
			),
		})
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

func ClonseDB(gorm *gorm.DB) error {
	db, err := gorm.DB()
	if err != nil {
		return err
	}
	db.Close()
	return nil
}
