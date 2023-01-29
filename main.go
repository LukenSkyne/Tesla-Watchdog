package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"tesla-watchdog/internal/discord"
	"tesla-watchdog/pkg/tesla"
	"time"
)

var (
	log    *zap.SugaredLogger
	client *tesla.Client
	car    *tesla.Vehicle
)

func init() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	log = logger.Sugar()

	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}
}

func main() {
	discordLogger := discord.NewDiscord(log)

	var ok bool
	log, ok = discordLogger.Start()

	if ok {
		defer discordLogger.Stop()
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	shutdown := make(chan bool)

	log.Infow("Started")

	client = tesla.NewClient(log)
	car = client.UseMainVehicle()

	go func() {
		for {
			select {
			case <-shutdown:
				return
			default:
			}

			tick()
			time.Sleep(1 * time.Second)
		}
	}()

	<-stop
	log.Infow("Gracefully shutting down...")
	shutdown <- true
}

var (
	lastTick = time.Now()
	sleeping = false
	nextLock = false
)

func tick() {
	elapsed := time.Now().Sub(lastTick).Seconds()

	if elapsed < 10 || (sleeping && elapsed < 60) {
		return
	}

	lastTick = time.Now()
	//log.Debugw("Tick", "elapsed", elapsed)

	info := car.GetInfo()

	if info == nil {
		//log.Error("gathering info failed")
		return
	}

	sleeping = info.Response.State != "online"

	if sleeping {
		//log.Debugw("car is sleeping", "state", info.Response.State)
		return
	}

	//car.WakeUp() // keep the car online
	driveState := car.GetDriveState()

	if driveState == nil {
		//log.Error("drive state unavailable")
		return
	}

	log.Infow("DriveState",
		"Shift", driveState.Response.ShiftState,
		"Speed", driveState.Response.Speed,
	)

	vehicleState := car.GetVehicleState()

	if vehicleState == nil {
		//log.Error("vehicle state unavailable")
		return
	}

	log.Debugw("VehicleState",
		"IsUserPresent", vehicleState.Response.IsUserPresent,
		"DisplayState", vehicleState.Response.CenterDisplayState,
		"Locked", vehicleState.Response.Locked,
		"FrontLeft", vehicleState.Response.DoorDriverFront,
		"FrontRight", vehicleState.Response.DoorPassengerFront,
		"BackLeft", vehicleState.Response.DoorDriverRear,
		"BackRight", vehicleState.Response.DoorPassengerRear,
		"FrontTrunk", vehicleState.Response.DoorFrontTrunk,
		"RearTrunk", vehicleState.Response.DoorRearTrunk,
	)

	if vehicleState.Response.Locked {
		nextLock = false
		log.Debug("already locked")
		return
	}

	shouldLock := !vehicleState.Response.IsUserPresent &&
		vehicleState.Response.DoorDriverFront == 0 &&
		vehicleState.Response.DoorDriverRear == 0 &&
		vehicleState.Response.DoorPassengerFront == 0 &&
		vehicleState.Response.DoorPassengerRear == 0 &&
		vehicleState.Response.DoorFrontTrunk == 0 &&
		vehicleState.Response.DoorRearTrunk == 0

	if !shouldLock {
		nextLock = false
		log.Debug("requirements to lock not met")
		return
	}

	if !nextLock {
		nextLock = true
		log.Debug("locking doors on next iteration")
		return
	}

	lockDoorsResult := car.LockDoors()

	if lockDoorsResult == nil {
		//log.Error("failed to lock doors")
		return
	}

	if !lockDoorsResult.Response.Result {
		log.Warnw("unable to lock doors", "result", lockDoorsResult)
		return
	}

	log.Info("Doors Locked")
	nextLock = false
}
