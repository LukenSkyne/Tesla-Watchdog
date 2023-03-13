package tesla

/* basics */

func (v *Vehicle) GetInfo() (*Wrapper[VehicleInfo], error) {
	return Get[Wrapper[VehicleInfo]](v.c, v.base)
}

func (v *Vehicle) WakeUp() (*Wrapper[WakeUpInfo], error) {
	return Post[Wrapper[WakeUpInfo]](v.c, v.base+"/wake_up", nil)
}

func (v *Vehicle) GetData() (*Wrapper[VehicleData], error) {
	return Get[Wrapper[VehicleData]](v.c, v.base+"/vehicle_data")
}

/* data requests */

// TODO: add data requests

func (v *Vehicle) GetDriveState() (*Wrapper[DriveState], error) {
	return Get[Wrapper[DriveState]](v.c, v.data+"/drive_state")
}

func (v *Vehicle) GetVehicleState() (*Wrapper[VehicleState], error) {
	return Get[Wrapper[VehicleState]](v.c, v.data+"/vehicle_state")
}

/* commands */

// TODO: add commands

func (v *Vehicle) LockDoors() (*Wrapper[CommandResponse], error) {
	return Post[Wrapper[CommandResponse]](v.c, v.cmd+"/door_lock", nil)
}

/*
func (v *Vehicle) UnlockDoors() (*Wrapper[CommandResponse], error) {
	return Post[Wrapper[CommandResponse]](v.c, v.cmd+"/door_unlock", nil)
}
*/

/*
func (v *Vehicle) SetTemps() (*Wrapper[CommandResponse], error) {
	jsonData, _ := json.Marshal(map[string]interface{}{
		"driver_temp":    21.0,
		"passenger_temp": 21.0,
	})

	return Post[Wrapper[CommandResponse]](v.c, v.cmd+"/set_temps", bytes.NewBuffer(jsonData))
}
*/
