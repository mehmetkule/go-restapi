package main

import (
	"flag"
	"github.com/mehmetkule/go-restapi/logger"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
)

var configPath = flag.String("config", "./config.yaml", "path to config file")

const appName = "cervice"

func main() {
	flag.Parse()
	var conf *Config
	var err error
	if *configPath != "" {
		if conf, err = NewConfig(*configPath); err != nil {
			logger.Logger().Error("Failed to read configuration file", zap.Error(err))

		}

		//start pprof
		go func() {
			logrus.Println(http.ListenAndServe("localhost:6060", nil))
		}()

		logger.Logger().Info("Read config file")
		conf.AppName = appName
		app := NewApp(conf)
		if err = app.Initialize(); err != nil {
			logger.Logger().Fatal("Failed to initialize the application", zap.Error(err))
		}

		logger.Logger().Info("Running server....")
		if err = app.Run(); err != nil {
			logger.Logger().Error("Failed to start application", zap.Error(err))
		}
		logger.Logger().Warn("Application stopped", zap.String("", conf.AppName))
	}
}
