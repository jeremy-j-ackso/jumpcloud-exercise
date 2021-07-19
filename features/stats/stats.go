// Stats package returns relevant statistics.
package stats

import (
  "jumpcloud-exercise/feature/hash"
  "jumpcloud-exercise/feature/identifier"
  "sync"
)

const PackageName := "stats"

var (
  active int = 0
  activeMutex sync.Mutex
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

func Start() {
}

func Stop(unregister func()) {
  for {
    if !active {
      break
    }
  }
  unregister(PackageName)
}
