package tesla

/* basics */

func (v *Vehicle) GetInfo() *VehicleInfoWrapper {
	return Get[VehicleInfoWrapper](v.c, v.base)
}

func (v *Vehicle) WakeUp() *WakeUpInfoWrapper {
	return Post[WakeUpInfoWrapper](v.c, v.base+"/wake_up", nil)
}

/* data requests */

// TODO: add data requests

func (v *Vehicle) GetDriveState() *DriveStateWrapper {
	return Get[DriveStateWrapper](v.c, v.data+"/drive_state")
}

func (v *Vehicle) GetVehicleState() *VehicleStateWrapper {
	return Get[VehicleStateWrapper](v.c, v.data+"/vehicle_state")
}

/* commands */

// TODO: add commands

func (v *Vehicle) LockDoors() *CommandResponseWrapper {
	return Post[CommandResponseWrapper](v.c, v.cmd+"/door_lock", nil)
}

/*
func (v *Vehicle) UnlockDoors() *CommandResponse {
	return Post[CommandResponse](v.c, v.cmd+"/door_unlock", nil)
}
*/

/*
func (v *Vehicle) SetTemps() *CommandResponse {
	jsonData, _ := json.Marshal(map[string]interface{}{
		"driver_temp":    21.0,
		"passenger_temp": 21.0,
	})

	return Post[CommandResponse](v.c, v.cmd+"/set_temps", bytes.NewBuffer(jsonData))
}
*/
