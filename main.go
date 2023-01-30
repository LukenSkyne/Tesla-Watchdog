package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strings"
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
	lastTick     = time.Now()
	lastSleeping = false
	nextLock     = false
)

func checkTimeout(msg string) bool {
	return len(msg) < 30 && strings.Contains(msg, "timeout")
}

func tick() {
	elapsed := time.Now().Sub(lastTick).Seconds()

	if elapsed < 10 || (lastSleeping && elapsed < 30) {
		return
	}

	//log.Debugw("Tick", "elapsed", elapsed)
	lastTick = time.Now()
	info, err := car.GetInfo()

	if err != nil {
		log.Errorf("GetInfo | %v", err)
		return
	}

	if info.Response == nil {
		if !checkTimeout(info.Error) {
			log.Warnf("GetInfo | %v", info.Error)
		}
		return
	}

	sleeping := info.Response.State != "online"

	if lastSleeping != sleeping {
		lastSleeping = sleeping
		log.Infof("car is now %v", info.Response.State)
	}

	if sleeping {
		return
	}

	latestData, err := car.GetLatestData()

	if err != nil {
		log.Errorf("GetLatestData | %v", err)
		return
	}

	if latestData.Response == nil {
		if !checkTimeout(latestData.Error) {
			log.Warnf("GetLatestData | %v", latestData.Error)
		}
		return
	}

	driveState := latestData.Response.Legacy.DriveState
	vehicleState := latestData.Response.Legacy.VehicleState

	if driveState.ShiftState != nil {
		log.Debugw("car is not idle",
			"ShiftState", driveState.ShiftState,
			"Speed", driveState.Speed,
		)
		return
	}

	if vehicleState.Locked {
		nextLock = false
		return
	}

	shouldLock := !vehicleState.IsUserPresent &&
		vehicleState.DoorDriverFront == 0 &&
		vehicleState.DoorDriverRear == 0 &&
		vehicleState.DoorPassengerFront == 0 &&
		vehicleState.DoorPassengerRear == 0 &&
		vehicleState.DoorFrontTrunk == 0 &&
		vehicleState.DoorRearTrunk == 0

	log.Debugw("VehicleState",
		"IsUserPresent", vehicleState.IsUserPresent,
		"DisplayState", vehicleState.CenterDisplayState,
		"Locked", vehicleState.Locked,
		"FrontLeft", vehicleState.DoorDriverFront,
		"FrontRight", vehicleState.DoorPassengerFront,
		"BackLeft", vehicleState.DoorDriverRear,
		"BackRight", vehicleState.DoorPassengerRear,
		"FrontTrunk", vehicleState.DoorFrontTrunk,
		"RearTrunk", vehicleState.DoorRearTrunk,
	)

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

	lockDoorsResult, err := car.LockDoors()

	if err != nil {
		log.Errorf("LockDoors | %v", err)
		return
	}

	if lockDoorsResult.Response == nil {
		if !checkTimeout(lockDoorsResult.Error) {
			log.Warnf("LockDoors | %v", lockDoorsResult.Error)
		}
		return
	}

	if !lockDoorsResult.Response.Result {
		log.Warnw("LockDoors | unable to lock doors", "result", lockDoorsResult)
		return
	}

	log.Info("doors locked")
	nextLock = false
}
