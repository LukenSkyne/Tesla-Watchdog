import axios, { AxiosInstance } from "axios"
import * as dotenv from "dotenv"
import * as fs from "fs"
import { logInfo, logError } from "./logger"

const OWNER_API = "https://owner-api.teslamotors.com"
const AUTH_API = "https://auth.tesla.com"
const TOKEN_REFRESH_TIMER = 3600

function extractJWT(jwt): JWTContent {
	return JSON.parse(Buffer.from(jwt.split(".")[1], "base64").toString())
}

function checkRenewJWT(jwt) {
	const exp = jwt !== "" ? extractJWT(jwt).exp : 0

	return (exp - (new Date().valueOf() / 1000)) < TOKEN_REFRESH_TIMER
}

class TeslaClientBase {
	private accessToken: string
	private refreshToken: string
	protected readonly mainVehicleId: string
	api: AxiosInstance

	constructor() {
		dotenv.config()

		this.accessToken = process.env["ACCESS_TOKEN"]
		this.refreshToken = process.env["REFRESH_TOKEN"]
		this.mainVehicleId = process.env["MAIN_VEHICLE_ID"]

		if (this.accessToken === undefined || this.refreshToken === undefined || this.mainVehicleId === undefined) {
			logError("TeslaClient", "please configure your .env and start again")
			process.exit(1)
		}

		this.api = axios.create({
			baseURL: OWNER_API,
			headers: {
				"Content-Type": "application/json",
			},
		})
		this.api.interceptors.request.use(async (config) => {
			if (checkRenewJWT(this.accessToken)) {
				logInfo("TeslaClient::interceptor", "Token Expired")
				await this.refreshAccessToken()
			}

			config.headers["Authorization"] = `Bearer ${this.accessToken}`
			return config
		})
	}

	async refreshAccessToken() {
		try {
			const res: TokenRefreshResponse = (await axios.post("/oauth2/v3/token", {
				grant_type: "refresh_token",
				client_id: "ownerapi",
				refresh_token: this.refreshToken,
				scope: "openid email offline_access",
			}, {
				baseURL: AUTH_API,
			})).data

			this.accessToken = res.access_token
			this.refreshToken = res.refresh_token
			logInfo("TeslaClient::refreshAccessToken", "Success")

			try {
				const env = `ACCESS_TOKEN="${this.accessToken}"\nREFRESH_TOKEN="${this.refreshToken}"\nMAIN_VEHICLE_ID="${this.mainVehicleId}"\n`
				fs.writeFileSync(`${process.cwd()}/.env`, env)
			} catch (err) {
				logError("TeslaClient::refreshAccessToken", "Unable to write .env:", err)
			}
		} catch (err) {
			logError("TeslaClient::refreshAccessToken", "Unable to refresh token:", err.response?.status ?? err)
		}
	}

	protected async get(url: string) {
		return (await this.api.get(url)).data.response
	}

	protected async post(url: string, data?: object) {
		return (await this.api.post(url, data)).data.response
	}
}

class TeslaVehicleBase {
	private client: TeslaClientBase
	private baseURL: string
	id: string | number

	constructor(client: TeslaClientBase, id: string | number) {
		this.client = client
		this.baseURL = `/api/1/vehicles/${id}`
		this.id = id
	}

	protected async get(url: string) {
		return (await this.client.api.get(this.baseURL + url)).data.response
	}

	protected async getOrNull(url: string) {
		try	{
			return (await this.client.api.get(this.baseURL + url)).data.response
		} catch (err) {
			return null
		}
	}

	protected async post(url: string, data?: object) {
		return (await this.client.api.post(this.baseURL + url, data)).data.response
	}
}

export class TeslaClient extends TeslaClientBase {
	initVehicle(id) {
		return new TeslaVehicle(this, id)
	}

	initMainVehicle() {
		return new TeslaVehicle(this, this.mainVehicleId)
	}

	async getProfile(): Promise<Profile> {
		return this.get("/api/1/users/me")
	}

	async getVehicles(): Promise<VehicleInfo[]> {
		return this.get("/api/1/vehicles")
	}

	async getVehicle(id): Promise<VehicleInfo> {
		return this.get(`/api/1/vehicles/${id}`)
	}
}

class TeslaVehicle extends TeslaVehicleBase {

	async wakeUp(): Promise<VehicleWakeInfo> {
		return this.post("/wake_up")
	}

	/* Data Requests */

	async getLatestData() {
		return this.getOrNull("/latest_vehicle_data")
	}

	async getData() {
		return this.get("/vehicle_data")
	}

	async getState() {
		return this.get("/data_request/vehicle_state")
	}

	async getDriveState() {
		return this.get("/data_request/drive_state")
	}

	async getClimateState() {
		return this.get("/data_request/climate_state")
	}

	async getChargeState() {
		return this.get("/data_request/charge_state")
	}

	async getReleaseNotes() {
		return this.get("/release_notes")
	}

	/* Commands */

	async autoConditioningStart() {
		return this.post("/command/auto_conditioning_start")
	}

	async autoConditioningStop() {
		return this.post("/command/auto_conditioning_stop")
	}

	async honkHornTwice() {
		return this.post("/command/honk_horn")
	}

	async flashHeadlights() {
		return this.post("/command/flash_lights")
	}

	async remoteStartDrive() {
		return this.post("/command/remote_start_drive")
	}

	async unlockDoors() {
		return this.post("/command/door_unlock")
	}

	async lockDoors() {
		return this.post("/command/door_lock")
	}

	async openChargePort() {
		return this.post("/command/charge_port_door_open")
	}

	async closeChargePort() {
		return this.post("/command/charge_port_door_close")
	}
}
