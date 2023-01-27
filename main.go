package main

import (
	"go.uber.org/zap"
	"tesla-watchdog/pkg/tesla"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	log := logger.Sugar()

	client := tesla.NewClient(log)
	client.DoSomething()
}
