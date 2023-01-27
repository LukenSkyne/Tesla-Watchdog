package tesla

type TokenRefreshResponse struct {
	IdToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
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
