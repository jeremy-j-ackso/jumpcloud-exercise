// Increments identifier.
package identifier

import (
  "log"
  "sync"
)

const PackageName string = "identifier"

var (
  active int = 0
  activeMutex sync.Mutex
  id int = 0
  mux sync.Mutex
)

// Used for indicating up/down status to Supervisor.
func activate() {
  activeMutex.Lock()
  active++
  activeMutex.Unlock()
}

// Used for indicating up/down status to Supervisor.
func deactivate() {
  activeMutex.Lock()
  active--
  activeMutex.Unlock()
}

// Returns the next identifier number.
func Get() int {
  activate()
  mux.Lock()
  id++
  output := id
  mux.Unlock()
  deactivate()
  return output
}

// Returns the current identifier number.
func Current() int {
  activate()
  mux.Lock()
  output := id
  mux.Unlock()
  deactivate()
  return output
}

// Registers as up with the Supervisor package.
func Start() {
  log.Printf("%s Package started", PackageName)
}

// Register as down with the Supervisor package when all work complete.
func Stop(unregister func(string)) {
  log.Printf("%s Package stopped", PackageName)
  for {
    if active == 0 {
      break
    }
  }
  unregister(PackageName)
}
