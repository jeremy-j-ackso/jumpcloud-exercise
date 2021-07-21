// Repository Package is where we would typically hook in a database or something.
// In this case we're just using a Map as a stand-in.
package hashrepository

import (
  "errors"
  "sync"
  "time"
)

const PackageName string = "hashrepository"

var (
  database = make(map[int]Record)
  active int = 0
  activeMutex sync.Mutex
)

// Stand-in for a table schema.
type Record struct {
  Hash string
  Hashtime time.Time
}

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

// Stand-in for a select function offered by a DB library.
func query(id int, ch chan<- Record) {
  output, present := database[id]
  if present {
    ch <- output
  } else {
    close(ch)
  }
  return
}

// Stand-in for an upsert function offered by a DB library.
func upsert(id int, data Record) {
  activate()
  database[id] = data
  deactivate()
}

// Sets up the repository and database access.
func Start() {
}

// Tears down the repository and database access.
func Stop(unregister func(string)) {
  for {
    if active == 0 {
      break
    }
  }
  unregister(PackageName)
}

// Returns a value from the repository by ID.
func Get(id int) (Record, error) {
  ch := make(chan Record)
  go query(id, ch)
  record, ok := <-ch
  if ok {
    return record, nil
  } else {
    return Record{"", time.Unix(0, 0)}, errors.New("ID %d does not exist")
  }
}

// Upserts a value into the repository by ID.
func Put(id int, hash string, hashtime time.Time) (int, error) {
  if active > 0 {
    // TODO: Fire and forget upserts should be refactored away if/when an actual DB implementation is brought in..
    go upsert(id, Record{hash, hashtime})
    return 1, nil
  } else {
    return 0, errors.New("Shutting down")
  }
}
