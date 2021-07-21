// Registrar package stores calculated hashes and the time they may be released.
package registrar

import (
  "jumpcloudexercise/services/repository"
  "sync"
  "time"
)

const PackageName string = "registrar"

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

func Get(id int) (repository.Record, error) {
  activate()
  output, error := repository.Get(id)
  deactivate()
  return output, error
}

func Put(id int, hash string, hashtime time.Time) {
  activate()
  repository.Put(id, hash, hashtime) // TODO: Should implement retry logic.
  deactivate()
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
