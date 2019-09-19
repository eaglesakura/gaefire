package gaefire

import "log"

func logInfo(message string) {
	log.Printf("gaefire.info: " + message)
}

func logDebug(message string) {
	log.Printf("gaefire.debug: " + message)
}

func logError(message string) {
	log.Printf("gaefire.error: " + message)
}
