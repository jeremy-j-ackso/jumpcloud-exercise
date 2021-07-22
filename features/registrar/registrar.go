// Registrar package stores calculated hashes and the time they may be released.
package registrar

import (
  "jumpcloudexercise/services/repository"
  "log"
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
  repository.Put(id, hash, hashtime)
  deactivate()
}

func Start() {
  log.Printf("%s Package started", PackageName)
}

func Stop(unregister func(string)) {
  log.Printf("Stop %s package", PackageName)
  for {
    if active == 0 {
      break
    }
  }
  unregister(PackageName)
}
