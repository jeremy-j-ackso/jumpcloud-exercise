// Registrar package stores calculated hashes and the time they may be released.
package registrar

import (
  "jumpcloudexercises/services/hashrepository"
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

func Get(id) (hashrepository.Record, error) {
  activate()
  output, error := hashrepository.Get(id)
  deactivate()
  return output, error
}

func Put(id int, hash string, hashtime time.Time) {
  activate()
  ok, error := hashrepository.Put(id, hash, hashtime) // TODO: Should implement retry logic.
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
