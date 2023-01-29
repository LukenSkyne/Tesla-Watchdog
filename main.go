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

	client = tesla.NewClient(log)
	car = client.UseMainVehicle()

	log.Infow("Started")

	//tmp := car.GetVehicleState()
	//
	//if tmp != nil {
	//	log.Infow("TEMP", "tmp", tmp)
	//}

	for {
		select {
		case <-sc:
			log.Infow("Stopped")
			return
		default:
		}

		tick()
		time.Sleep(1 * time.Second)
	}
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
	log.Debugw("Tick", "elapsed", elapsed)

	info := car.GetInfo()

	if info == nil {
		log.Error("gathering info failed")
		return
	}

	sleeping = info.Response.State != "online"

	if sleeping {
		log.Debugw("car is sleeping", "state", info.Response.State)
		return
	}

	car.WakeUp() // keep the car online
	driveState := car.GetDriveState()

	if driveState == nil {
		log.Error("drive state unavailable")
		return
	}

	log.Infow("DriveState",
		"Shift", driveState.Response.ShiftState,
		"Speed", driveState.Response.Speed,
	)

	vehicleState := car.GetVehicleState()

	if vehicleState == nil {
		log.Error("vehicle state unavailable")
		return
	}

	log.Infow("VehicleState",
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
		log.Error("failed to lock doors")
		return
	}

	if !lockDoorsResult.Response.Result {
		log.Debugw("unable to lock doors", "result", lockDoorsResult)
		return
	}

	log.Info("Locking Doors")
	nextLock = false
}
