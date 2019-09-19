package gaefire

import "log"

func logInfo(message string) {
	log.Println("gaefire.info: " + message)
}

func logDebug(message string) {
	log.Println("gaefire.debug: " + message)
}

func logError(message string) {
	log.Println("gaefire.error: " + message)
}
