package tesla

func (c *Client) GetUser() (*UserInfo, error) {
	return Get[UserInfo](c, "/users/me")
}

func (c *Client) GetProductList() (*ProductList, error) {
	return Get[ProductList](c, "/products")
}

// only available through fleet api: https://developer.tesla.com/docs/fleet-api/endpoints/vehicle-endpoints#list
func (c *Client) GetVehicleList() (*VehicleList, error) {
	return Get[VehicleList](c, "/vehicles")
}

func (c *Client) GetVehicleInfo(id string) (*Wrapper[VehicleInfo], error) {
	return Get[Wrapper[VehicleInfo]](c, "/vehicles/"+id)
}
