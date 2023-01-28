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

type VehicleState struct {
	ApiVersion               int    `json:"api_version"`
	AutoparkStateV2          string `json:"autopark_state_v2"` // "unavailable" | "standby"
	AutoparkStyle            string `json:"autopark_style"`    // "standard"
	CalendarSupported        bool   `json:"calendar_supported"`
	CarVersion               string `json:"car_version"` // "2022.44.30.1 9bf26f084642"
	CenterDisplayState       int    `json:"center_display_state"`
	DashcamClipSaveAvailable bool   `json:"dashcam_clip_save_available"`
	DashcamState             string `json:"dashcam_state"` // "Recording"
	DoorDriverFront          int    `json:"df"`
	DoorDriverRear           int    `json:"dr"`
	DoorPassengerFront       int    `json:"pf"`
	DoorPassengerRear        int    `json:"pr"`
	DoorFrontTrunk           int    `json:"ft"`
	DoorRearTrunk            int    `json:"rt"`
	WindowDriverFront        int    `json:"fd_window"`
	WindowPassengerFront     int    `json:"fp_window"`
	WindowDriverRear         int    `json:"rd_window"`
	WindowPassengerRear      int    `json:"rp_window"`
	FeatureBitmask           string `json:"feature_bitmask"` // "dffbff,0"
	IsUserPresent            bool   `json:"is_user_present"`
	Locked                   bool   `json:"locked"`
	MediaInfo                struct {
		AudioVolume          float64 `json:"audio_volume"`
		AudioVolumeIncrement float64 `json:"audio_volume_increment"`
		AudioVolumeMax       float64 `json:"audio_volume_max"`
		MediaPlaybackStatus  string  `json:"media_playback_status"` // "Stopped"
		NowPlayingAlbum      string  `json:"now_playing_album"`
		NowPlayingArtist     string  `json:"now_playing_artist"`
		NowPlayingDuration   float64 `json:"now_playing_duration"`
		NowPlayingElapsed    float64 `json:"now_playing_elapsed"`
		NowPlayingSource     string  `json:"now_playing_source"`
		NowPlayingStation    string  `json:"now_playing_station"`
		NowPlayingTitle      string  `json:"now_playing_title"`
	} `json:"media_info"`
	MediaState struct {
		RemoteControlEnabled bool `json:"remote_control_enabled"`
	} `json:"media_state"`
	NotificationsSupported  bool    `json:"notifications_supported"`
	Odometer                float64 `json:"odometer"`
	ParsedCalendarSupported bool    `json:"parsed_calendar_supported"`
	RemoteStart             bool    `json:"remote_start"`
	RemoteStartEnabled      bool    `json:"remote_start_enabled"`
	RemoteStartSupported    bool    `json:"remote_start_supported"`
	SantaMode               int     `json:"santa_mode"`
	SentryMode              bool    `json:"sentry_mode"`
	SentryModeAvailable     bool    `json:"sentry_mode_available"`
	ServiceMode             bool    `json:"service_mode"`
	ServiceModePlus         bool    `json:"service_mode_plus"`
	SoftwareUpdate          struct {
		DownloadPerc        int    `json:"download_perc"`
		ExpectedDurationSec int    `json:"expected_duration_sec"`
		InstallPerc         int    `json:"install_perc"`
		Status              string `json:"status"`  // ""
		Version             string `json:"version"` // ""
	} `json:"software_update"`
	SpeedLimitMode struct {
		Active          bool    `json:"active"`
		CurrentLimitMph float64 `json:"current_limit_mph"`
		MaxLimitMph     float64 `json:"max_limit_mph"`
		MinLimitMph     float64 `json:"min_limit_mph"`
		PinCodeSet      bool    `json:"pin_code_set"`
	} `json:"speed_limit_mode"`
	Timestamp                  int     `json:"timestamp"`
	TpmsPressureFl             float64 `json:"tpms_pressure_fl"`
	TpmsPressureFr             float64 `json:"tpms_pressure_fr"`
	TpmsPressureRl             float64 `json:"tpms_pressure_rl"`
	TpmsPressureRr             float64 `json:"tpms_pressure_rr"`
	TpmsRcpFrontValue          float64 `json:"tpms_rcp_front_value"`
	TpmsRcpRearValue           float64 `json:"tpms_rcp_rear_value"`
	TpmsSoftWarningFl          bool    `json:"tpms_soft_warning_fl"`
	TpmsSoftWarningFr          bool    `json:"tpms_soft_warning_fr"`
	TpmsSoftWarningRl          bool    `json:"tpms_soft_warning_rl"`
	TpmsSoftWarningRr          bool    `json:"tpms_soft_warning_rr"`
	TpmsHardWarningFl          bool    `json:"tpms_hard_warning_fl"`
	TpmsHardWarningFr          bool    `json:"tpms_hard_warning_fr"`
	TpmsHardWarningRl          bool    `json:"tpms_hard_warning_rl"`
	TpmsHardWarningRr          bool    `json:"tpms_hard_warning_rr"`
	TpmsLastSeenPressureTimeFl int     `json:"tpms_last_seen_pressure_time_fl"`
	TpmsLastSeenPressureTimeFr int     `json:"tpms_last_seen_pressure_time_fr"`
	TpmsLastSeenPressureTimeRl int     `json:"tpms_last_seen_pressure_time_rl"`
	TpmsLastSeenPressureTimeRr int     `json:"tpms_last_seen_pressure_time_rr"`
	ValetMode                  bool    `json:"valet_mode"`
	ValetPinNeeded             bool    `json:"valet_pin_needed"`
	VehicleName                string  `json:"vehicle_name"`
	WebcamAvailable            bool    `json:"webcam_available"`
}

type VehicleStateWrapper struct {
	Response VehicleState `json:"response"`
}
