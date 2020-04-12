package main

import (
	"github.com/houstonj1/go-postgres/config"
	"github.com/houstonj1/go-postgres/pq"
	"go.uber.org/zap"
)

func main() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic("error creating logger")
	}
	logger := zapLogger.Sugar()

	config := config.NewConfig()
	logger.Debugf("%s", config.Print())

	logger.Info("-------------------------------------")
	logger.Info("-------------  lib/pq  --------------")
	logger.Info("-------------------------------------")
	pq.Pq(logger)
}
