// Increments identifier.
package identifier

import (
  "sync"
)

var (
  active int = 0
  activeMutex sync.Mutex
  id int = 0
  mux sync.Mutex
  packageName := "identifier"
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
func Get() {
  activate()
  mux.Lock()
  id++
  output := id
  mux.Unlock()
  deactivate()
  return output
}

// Returns the current identifier number.
func Current() {
  activate()
  mux.Lock()
  output := id
  mux.Unlock()
  deactivate()
  return output
}

// Registers as up with the Supervisor package.
func Start(register func()) {
  register(packageName)
}

// Register as down with the Supervisor package when all work complete.
func Stop(unregister func()) {
  for {
    if !active {
      break
    }
  }
  unregister(packageName)
}
