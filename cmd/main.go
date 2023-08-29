package main

import (
	"fmt"
	"http-service/cmd/log"
	"os"
)

func main() {
	configPrefix := os.Getenv("config_prefix")
	if configPrefix == "" {
		fmt.Println("please set env: config_prefix")
		os.Exit(1)
	}
	conf, err := ParseConfig(configPrefix)
	if err != nil {
		panic(err)
	}

	logger, err := log.NewDefaultLogger(conf.Log.Level, conf.Log.AppName, conf.Log.EnvName)
	if err != nil {
		panic(err)
	}

	gormDB, err := NewDB(conf, logger)
	if err != nil {
		panic(err)
	}

	s, err := NewServer(conf, gormDB, logger)
	if err != nil {
		panic(err)
	}

	if err := s.Run(); err != nil {
		panic(err)
	}

	defer ClonseDB(gormDB)
}
