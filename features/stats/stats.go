// Stats package returns relevant statistics.
package stats

import (
  "jumpcloud-exercise/feature/hash"
  "jumpcloud-exercise/feature/identifier"
  "sync"
)

var (
  active int = 0
  activeMutex sync.Mutex
  packageName := "stats"
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

func Get() (total int, average int64) {
  activate()
  total := identifier.Current()
  average := hash.GetAvg()
  deactivate()
  return
}

func Start(register func()) {
  register(packageName)
}

func Stop(unregister func()) {
  for {
    if !active {
      break
    }
  }
  unregister(packageName)
}
