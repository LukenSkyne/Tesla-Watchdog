package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"tesla-watchdog/internal/discord"
	"tesla-watchdog/internal/watchdog"
	"tesla-watchdog/pkg/tesla"
)

var (
	log           *zap.SugaredLogger
	discordClient *discord.Discord
	watchDog      *watchdog.WatchDog
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	log = logger.Sugar()

	teslaClient := tesla.NewClient(log)

	if teslaClient.FirstTimeSetup() {
		return
	}

	if err := godotenv.Load(); err == nil {
		discordClient = discord.NewDiscord(
			log,
			os.Getenv("BOT_TOKEN"),
			os.Getenv("CHANNEL_ID"),
		)

		var ok bool
		log, ok = discordClient.Start()

		if ok {
			defer discordClient.Stop()
		} else {
			discordClient = nil
		}
	}

	watchDog = watchdog.NewWatchDog(log, discordClient, teslaClient)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	log.Infow("started")

	go watchDog.Run()

	<-stop
	log.Infow("gracefully shutting down...")
	watchDog.Stop()
}
