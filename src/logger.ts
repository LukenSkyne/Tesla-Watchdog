
function timestamp() {
	return `[${ new Date().toLocaleString().replace(", ", " @ ")}]`
}

export function logInfo(namespace, ...args) {
	console.info(timestamp(), namespace, "|", ...args)
}

export function logError(namespace, ...args) {
	console.error(timestamp(), namespace, "|", ...args)
}
