// Hash package calculates hashes and stores the avergae time of hash calculating.
package hash

import (
  "crypto/sha512"
  "encoding/base64"
  "time"
)

var (
  active int = 0
  activeMutex sync.Mutex
  hashes int = 0
  hashtime int64 = 0
  hashMux sync.Mutex
  hashtimeMux sync.Mutex
  packageName := "hash"
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

func calculateHash(toHash string) string {
  hashed := sha512.Sum512([]byte(toHash))
  asBase64 := base64.StdEncoding.EncodeToSting(hashed)
  return asBase64
}

func calculateNewHashAvg(newDuration int64, currentHashes int) {
  return (hashtime + newDuration) / currentHashes
}

func Hash(password string) string {
  activate()

  t0 := time.Now()
  output := calculateHash(password)
  t1 := time.Now()

  hashMux.Lock()
  hashes++

  hashtimeMux.Lock()
  hashtime = calculateNewHashAvg(t1.Sub(t0), hashes)
  hashtimeMux.Unlock()

  hashMux.Unlock()

  deactivate()
  return output
}

func GetAvg() int64 {
  activate()

  hashtimeMux.Lock()
  output := hashtime
  hashtimeMux.Unlock()

  deactivate()
  return output
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
