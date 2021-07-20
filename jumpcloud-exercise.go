// Jumpcloud-Exercise implements an HTTP API for assigned routes and features.
package jumpcloudexercises

import (
  "jumpcloudexercises/app/app"
)

func main() {
  app.Start()
}

//TODO: Add a SIGTERM handler to stop the app.
