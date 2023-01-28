package main

import (
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"tesla-watchdog/pkg/tesla"
	"time"
)

var (
	log    *zap.SugaredLogger
	client *tesla.Client
	car    *tesla.Vehicle
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	log = logger.Sugar()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	lastTick := time.Now()

	client = tesla.NewClient(log)
	car = client.UseMainVehicle()

	log.Infow("Started")

	for {
		select {
		case <-sc:
			log.Infow("Stopped")
			return
		default:
		}

		if time.Now().Sub(lastTick).Seconds() < 30 {
			time.Sleep(1 * time.Second)
			continue
		}

		lastTick = time.Now()
		tick()
	}
}

func tick() {
	log.Debug("Tick")

	info := car.GetInfo()

	if info == nil {
		log.Error("gathering info failed")
		return
	}

	if info.Response.State == "asleep" {
		log.Debug("car is sleeping")
		return
	}

	log.Infow("VehicleInfo", "State", info.Response.State)

	//wakeUp := car.WakeUp()
	//
	//if wakeUp != nil {
	//	log.Infow("WakeUp", "Response", wakeUp.Response)
	//}
}
