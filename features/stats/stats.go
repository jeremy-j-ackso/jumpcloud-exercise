// Stats package returns relevant statistics.
package stats

import (
  "jumpcloudexercise/feature/hash"
  "jumpcloudexercise/feature/identifier"
  "sync"
)

const PackageName string = "stats"

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

func Stop(unregister func(string)) {
  for {
    if active == 0 {
      break
    }
  }
  unregister(PackageName)
}
