// HTTP Server implementation
package httpserver

// TODO: It looks like there's probably a better way to deal with Up/Down stuff
// that utilizes Context and ConnState. Should reseearch and refactor. May be
// able to eliminate or smooth out the Supervisor package by doing so.

import (
  "net/http"
  "sync"
)

const packageName := "httpserver"

var (
  active int = 0
  activeMutex sync.Mutex
  srv http.Server
  mux http.ServeMux
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
func Register(path string, handler func()) {
  mux.HandleFunc(path, func(w http.ResponseWriter, r http.Request) {
    activate()
    ch := make(chan int)
    handler(w, r, ch)
    <-ch
    deactivate()
  })
}

// Sets up the repository and database access.
func Start(address string, port string, register func()) {
  srv = http.Server{
    Addr: address + ":" + port,
    Handler: mux
  }
  srv.ListenAndServe()
  register(packageName)
}

// Tears down the repository and database access.
func Stop(unregister func()) {
  srv.Shutdown()
  for {
    if !active {
      break
    }
  }
  unregister(packageName)
}

