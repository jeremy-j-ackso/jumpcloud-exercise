// Jumpcloud-Exercise implements an HTTP API for assigned routes and features.
package main

import (
	"jumpcloudexercise/app"
	"log"
)

func main() {
	log.Println("App starting")
	app.Start()
}

//TODO: Add a SIGTERM handler to stop the app.
