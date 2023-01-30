package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"tesla-watchdog/internal/discord"
	"tesla-watchdog/internal/watchdog"
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

	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}

	discordClient = discord.NewDiscord(
		log,
		os.Getenv("BOT_TOKEN"),
		os.Getenv("CHANNEL_ID"),
	)

	var ok bool
	log, ok = discordClient.Start()

	if ok {
		defer discordClient.Stop()
	}

	watchDog = watchdog.NewWatchDog(log, discordClient)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	log.Infow("started")

	go watchDog.Run()

	<-stop
	log.Infow("gracefully shutting down...")
	watchDog.Stop()
}
