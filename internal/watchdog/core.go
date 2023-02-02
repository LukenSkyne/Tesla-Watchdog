package watchdog

import (
	"go.uber.org/zap"
	"strings"
	"tesla-watchdog/internal/discord"
	"tesla-watchdog/pkg/tesla"
	"time"
)

type WatchDog struct {
	stop    chan bool
	log     *zap.SugaredLogger
	discord *discord.Discord
	client  *tesla.Client
	car     *tesla.Vehicle
	state   *State
}

type State struct {
	lastTick  time.Time
	lastState string
	wasAsleep bool
	doLock    bool
}

func NewWatchDog(log *zap.SugaredLogger, discord *discord.Discord, client *tesla.Client) *WatchDog {
	return &WatchDog{
		stop:    make(chan bool),
		log:     log,
		discord: discord,
		client:  client,
		car:     client.UseMainVehicle(),
		state: &State{
			lastTick:  time.Now(),
			lastState: "",
			wasAsleep: false,
			doLock:    false,
		},
	}
}

func (w *WatchDog) Run() {
	for {
		select {
		case <-w.stop:
			return
		default:
		}

		w.tick()
		time.Sleep(1 * time.Second)
	}
}

func (w *WatchDog) Stop() {
	w.stop <- true
}

// main WatchDog routine
func (w *WatchDog) tick() {
	elapsed := time.Now().Sub(w.state.lastTick).Seconds()

	if elapsed < 10 || (w.state.wasAsleep && elapsed < 30) {
		return // too early for another tick
	}

	w.state.lastTick = time.Now()

	info, err := w.car.GetInfo()

	if !validate(w, info, err, "GetInfo") {
		return // failed to get info
	}

	sleeping := info.Response.State != "online"

	if w.state.lastState != info.Response.State {
		w.state.lastState = info.Response.State
		w.log.Infof("car is now %v", info.Response.State)
	}

	if w.state.wasAsleep != sleeping {
		w.state.wasAsleep = sleeping

		if w.discord != nil {
			go w.discord.UpdateStatus(!sleeping)
		}
	}

	if sleeping {
		return // car is in sleep mode
	}

	latestData, err := w.car.GetLatestData()

	if !validate(w, latestData, err, "GetLatestData") {
		return // failed to get car data
	}

	driveState := latestData.Response.Legacy.DriveState
	vehicleState := latestData.Response.Legacy.VehicleState

	idle := driveState.ShiftState == nil

	if !idle {
		//w.log.Debugw("car is not idle",
		//	"ShiftState", driveState.ShiftState,
		//	"Speed", driveState.Speed,
		//)
		return // car is not idle
	}

	if vehicleState.Locked {
		if w.state.doLock {
			w.log.Info("doors already got locked")
		}

		w.state.doLock = false
		return // already locked
	}

	canLock := !vehicleState.IsUserPresent &&
		vehicleState.DoorDriverFront == 0 &&
		vehicleState.DoorDriverRear == 0 &&
		vehicleState.DoorPassengerFront == 0 &&
		vehicleState.DoorPassengerRear == 0 &&
		vehicleState.DoorFrontTrunk == 0 &&
		vehicleState.DoorRearTrunk == 0

	w.log.Debugw("VehicleState",
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

	if !canLock {
		w.state.doLock = false
		return // requirements to lock not met
	}

	if !w.state.doLock {
		w.state.doLock = true
		w.log.Info("locking doors on next iteration")
		return // locking doors on next iteration
	}

	lockDoorsResult, err := w.car.LockDoors()

	if !validate(w, lockDoorsResult, err, "LockDoors") {
		return // failed to lock doors
	}

	if !lockDoorsResult.Response.Result {
		w.log.Warnw("LockDoors | unable to lock doors", "result", lockDoorsResult)
		return // doors are unable to be locked (unknown reason)
	}

	w.log.Info("doors locked")
	w.state.doLock = false
}

func validate[T any](w *WatchDog, r *tesla.Wrapper[T], err error, name string) bool {
	if err != nil {
		w.log.Errorf("%v | %v", name, err)
		return false
	}

	if r.Response == nil {
		if !isRegularTimeout(r.Error) {
			if isMotherShipTimeout(r.Error) {
				w.log.Debugf("%v | %v", name, r.Error)
			} else {
				w.log.Warnf("%v | %v", name, r.Error)
			}
		}

		return false
	}

	return true
}

func isRegularTimeout(msg string) bool {
	return len(msg) < 30 && strings.Contains(msg, "timeout")
}

func isMotherShipTimeout(msg string) bool {
	return strings.Contains(msg, "operation_timedout")
}
