// Jumpcloud-Exercise implements an HTTP API for assigned routes and features.
package jumpcloud-exercises

import (
  "jumpcloud-exercises/app/app"
)

func main() {
  app.Start()
}

//TODO: Add a SIGTERM handler to stop the app.
