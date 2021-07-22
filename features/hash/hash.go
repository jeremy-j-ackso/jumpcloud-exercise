// Hash package calculates hashes and stores the avergae time of hash calculating.
package hash

import (
  "crypto/sha512"
  "encoding/base64"
  "log"
  "sync"
  "time"
)

const PackageName string = "hash"

var (
  active int = 0
  activeMutex sync.Mutex
  hashes int64 = 0
  hashtime int64 = 0
  hashtotal int64 = 0
  hashMux sync.Mutex
  hashtimeMux sync.Mutex
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

func calculateHash(toHash string, ch chan string) {
  hasher := sha512.New()
  hashed := hasher.Sum([]byte(toHash))
  asBase64 := base64.StdEncoding.EncodeToString(hashed)
  ch <- asBase64
  return
}

func calculateNewHashAvg(newDuration time.Duration, currentHashes int64) int64 {
  hashtotal += newDuration.Microseconds()
  return hashtotal / currentHashes
}

func Hash(password string) string {
  activate()

  t0 := time.Now()
  ch := make(chan string)
  go calculateHash(password, ch)
  output := <-ch
  t1 := time.Now()

  hashMux.Lock()
  hashes++

  hashtimeMux.Lock()
  hashtime = calculateNewHashAvg(t1.Sub(t0), hashes)
  hashtimeMux.Unlock()

  log.Printf("Hashing took %v", hashtime)

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

func Start() {
  log.Printf("%s Package Started", PackageName)
}

func Stop(unregister func(string)) {
  log.Printf("%s Package Stopped", PackageName)
  for {
    if active == 0 {
      break
    }
  }
  unregister(PackageName)
}
