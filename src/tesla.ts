import axios, { AxiosInstance } from "axios"
import * as dotenv from "dotenv"
import * as fs from "fs"
import { logInfo, logError } from "./logger.js"

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

	async getProfile() {
		return this.get("/api/1/users/me")
	}

	async getVehicles() {
		return this.get("/api/1/vehicles")
	}

	async getVehicle(id) {
		return this.get(`/api/1/vehicles/${id}`)
	}
}

class TeslaVehicle extends TeslaVehicleBase {
	async wakeUp() {
		return this.post("/wake_up")
	}

	async getLatestData() {
		return this.getOrNull("/latest_vehicle_data")
	}
}
