// HTTP Server implementation
package httpserver

// TODO: It looks like there's probably a better way to deal with Up/Down stuff
// that utilizes Context and ConnState. Should reseearch and refactor. May be
// able to eliminate or smooth out the Supervisor package by doing so.

import (
  "context"
  "log"
  "net/http"
  "sync"
  "time"
)

const PackageName string = "httpserver"

var (
  active int = 0
  activeMutex sync.Mutex
  srv http.Server
  mux *http.ServeMux = http.NewServeMux()
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

// Registers a URL path to a handler.
func Register(path string, handler func(w http.ResponseWriter, r *http.Request, ch chan int)) {
  mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
    t0 := time.Now()
    activate()
    cha := make(chan int)
    go handler(w, r, cha)
    <-cha
    deactivate()
    t1 := time.Now()
    log.Printf("%s took $vms", path, t1.Sub(t0).Microseconds())
  })
}

// Sets up the repository and database access.
func Start(address string, port string) {
  srv = http.Server{
    Addr: address + ":" + port,
    Handler: mux,
  }
  log.Println("HTTP Server Started")
  srv.ListenAndServe()
}

// Tears down the repository and database access.
func Stop() {
  go srv.Shutdown(context.Background())
  log.Println("HTTP Server Stopped")
  for {
    if active == 0 {
      break
    }
  }
}

