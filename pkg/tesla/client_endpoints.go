package tesla

func (c *Client) GetUser() *UserInfo {
	return Get[UserInfo](c, "/users/me")
}

func (c *Client) GetVehicleList() *VehicleList {
	return Get[VehicleList](c, "/vehicles")
}

func (c *Client) GetVehicleInfo(id string) *VehicleInfoWrapper {
	return Get[VehicleInfoWrapper](c, "/vehicles/"+id)
}
