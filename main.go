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
	car := client.UseMainVehicle()

	info := car.GetInfo()

	if info != nil {
		log.Infow("VehicleInfo", "State", info.Response.State)
	}

	wakeUp := car.WakeUp()

	if wakeUp != nil {
		log.Infow("WakeUp", "Response", wakeUp.Response)
	}
}
