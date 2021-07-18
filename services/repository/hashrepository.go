// Repository Package is where we would typically hook in a database or something.
// In this case we're just using a Map as a stand-in.
package hashrepository

import (
  "errors"
  "sync"
)

const packageName := "hashrepository"

var (
  database := make(map[int]Record)
  active int = 0
  activeMutex sync.Mutex
)

// Stand-in for a table schema.
type Record struct {
  Hash string
  Hashtime int64
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
func select(id, ch chan<- Record) {
  output, present := database[id]
  if present {
    ch <- output
  } else {
    close(ch)
  }
  return
}

// Stand-in for an upsert function offered by a DB library.
func upsert(id, data) {
  activate()
  database[id] = data
  deactivate()
}

// Sets up the repository and database access.
func Start(register func()) {
  register(packageName)
}

// Tears down the repository and database access.
func Stop(unregister func()) {
  for {
    if !active {
      break
    }
  }
  unregister(packageName)
}

// Returns a value from the repository by ID.
func Get(id) (Record, error) {
  ch := make(chan Record)
  go select(id, ch)
  record, ok := <-ch
  if ok {
    return record, nil
  } else {
    return nil, errors.New("ID %d does not exist")
}

// Upserts a value into the repository by ID.
func Put(id int, hash string, hashtime int64) (int, error) {
  if active {
    // TODO: Fire and forget upserts should be refactored away if/when an actual DB implementation is brought in..
    go upsert(id, Record{hash, hashtime})
    return 1, nil
  } else {
    return nil, errors.New("Shutting down")
  }
}
