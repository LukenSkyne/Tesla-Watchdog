import { TeslaClient } from "./tesla.js"

const client = new TeslaClient()
const car = client.initMainVehicle()

console.log("getProfile", await client.getProfile())
console.log("getVehicles", await client.getVehicles())
console.log("wakeUp", await car.wakeUp())
//console.log("getLatestData", await m3.getLatestData())
