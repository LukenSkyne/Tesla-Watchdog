interface TokenRefreshResponse {
	access_token: string;
	refresh_token: string;
	id_token: string;
	expires_in: number;
	token_type: string;
}

interface JWTContent {
	iss: string;
	aud: string[];
	azp: string;
	sub: string;
	scp: string[];
	amr: string[];
	exp: number;
	iat: number;
	auth_time: number;
}

interface Profile {
	email: string;
	full_name: string;
	profile_image_url: string;
}

type State = "online" | "asleep"

interface VehicleInfo {
	id: number;
	vehicle_id: number;
	vin: string;
	display_name: string;
	option_codes?: string;
	color: null; // unknown
	access_type: string;
	tokens: string[];
	state: State;
	in_service: boolean;
	id_s: string;
	calendar_enabled: boolean;
	api_version: number;
	backseat_token: null; // unknown
	backseat_token_updated_at: null; // unknown
}

interface VehicleWakeInfo {
	id: number;
	user_id: number;
	vehicle_id: number;
	vin: string;
	display_name: string;
	option_codes?: string;
	color: null; // unknown
	tokens: string[];
	state: State;
	in_service: boolean;
	id_s: string;
	calendar_enabled: boolean;
	api_version: number;
	backseat_token: null; // unknown
	backseat_token_updated_at: null; // unknown
}
