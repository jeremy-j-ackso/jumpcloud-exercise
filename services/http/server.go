// HTTP Server implementation
package httpserver

// TODO: It looks like there's probably a better way to deal with Up/Down stuff
// that utilizes Context and ConnState. Should reseearch and refactor. May be
// able to eliminate or smooth out the Supervisor package by doing so.

import (
  "context"
  "net/http"
  "sync"
)

const PackageName string = "httpserver"

var (
  active int = 0
  activeMutex sync.Mutex
  srv http.Server
  mux *http.ServeMux = http.NewServeMux()
)

// init {
//   mux = http.NewServeMux()
// }

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
    activate()
    cha := make(chan int)
    handler(w, r, cha)
    <-cha
    deactivate()
  })
}

// Sets up the repository and database access.
func Start(address string, port string) {
  srv = http.Server{
    Addr: address + ":" + port,
    Handler: mux,
  }
  srv.ListenAndServe()
}

// Tears down the repository and database access.
func Stop(unregister func(pkg string)) {
  srv.Shutdown(context.Background())
  for {
    if active == 0 {
      break
    }
  }
  unregister(PackageName)
}

