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
