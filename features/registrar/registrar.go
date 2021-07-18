// Registrar package stores calculated hashes and the time they may be released.
package registrar

import (
  "jumpcloud-exercises/services/hashrepository"
  "sync"
)

var (
  active int = 0
  activeMutex sync.Mutex
  packageName := "registrar"
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

func Get(id) (struct, error) {
  activate()
  output, error := hashrepository.Get(id)
  deactivate()
  return output, error
}

func Put(id int, hash string, hashtime int64) (int, error) {
  activate()
  ok, error := hashrepository.Put(id, hash, hashtime)
  deactivate()
  return ok, error
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
