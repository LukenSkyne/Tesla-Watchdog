package tesla

type Vehicle struct {
	Id   string
	base string
	cmd  string
	data string
	c    *Client
}

func (c *Client) UseVehicle(id string) *Vehicle {
	return &Vehicle{
		id,
		"/vehicles/" + id,
		"/vehicles/" + id + "/command",
		"/vehicles/" + id + "/data_request",
		c,
	}
}

func (c *Client) UseMainVehicle() *Vehicle {
	return c.UseVehicle(c.config.MainVehicle)
}
