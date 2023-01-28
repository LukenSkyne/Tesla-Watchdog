package tesla

type TokenRefreshResponse struct {
	IdToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type UserInfo struct {
	Email           string `json:"email"`
	FullName        string `json:"full_name"`
	ProfileImageUrl string `json:"profile_image_url"`
}

type VehicleList struct {
	Response []VehicleInfo `json:"response"`
	Count    int           `json:"count"`
}

type VehicleInfo struct {
	Id                     int      `json:"id"`
	IdString               string   `json:"id_s"`
	VehicleId              int      `json:"vehicle_id"`
	Vin                    string   `json:"vin"`
	DisplayName            string   `json:"display_name"`
	OptionCodes            string   `json:"option_codes"`
	Color                  *string  `json:"color"` // unknown
	AccessType             string   `json:"access_type"`
	Tokens                 []string `json:"tokens"`
	State                  string   `json:"state"` // "online" | "asleep"
	InService              bool     `json:"in_service"`
	CalendarEnabled        bool     `json:"calendar_enabled"`
	ApiVersion             int      `json:"api_version"`
	BackseatToken          *string  `json:"backseat_token"`            // unknown
	BackseatTokenUpdatedAt *string  `json:"backseat_token_updated_at"` // unknown
}

type VehicleInfoWrapper struct {
	Response VehicleInfo `json:"response"`
}

type WakeUpInfo struct {
	VehicleInfo
	UserId int `json:"user_id"`
}

type WakeUpInfoWrapper struct {
	Response WakeUpInfo `json:"response"`
}

type CommandResponse struct {
	Reason string `json:"reason"`
	Result bool   `json:"result"`
}

type DriveState struct {
	ActiveRouteLatitude            float64 `json:"active_route_latitude"`
	ActiveRouteLongitude           float64 `json:"active_route_longitude"`
	ActiveRouteTrafficMinutesDelay float64 `json:"active_route_traffic_minutes_delay"`
	GpsAsOf                        int     `json:"gps_as_of"`
	Heading                        int     `json:"heading"`
	Latitude                       float64 `json:"latitude"`
	Longitude                      float64 `json:"longitude"`
	NativeLatitude                 float64 `json:"native_latitude"`
	NativeLongitude                float64 `json:"native_longitude"`
	NativeLocationSupported        int     `json:"native_location_supported"`
	NativeType                     string  `json:"native_type"` // "wgs"
	Power                          int     `json:"power"`
	ShiftState                     *int    `json:"shift_state"`
	Speed                          *int    `json:"speed"`
	Timestamp                      int     `json:"timestamp"`
}

type DriveStateWrapper struct {
	Response DriveState `json:"response"`
}
