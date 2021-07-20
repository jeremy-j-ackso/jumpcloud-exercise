// App Package brings all the services and features together to provide a single point of entry for main.
package app

import (
  "encoding/json"
  "http"
  "jumpcloudexercises/features/hash"
  "jumpcloudexercises/features/identifier"
  "jumpcloudexercises/features/registrar"
  "jumpcloudexercises/features/stats"
  "jumpcloudexercises/services/httpserver"
  "jumpcloudexercises/services/repository"
  "jumpcloudexercises/services/supervisor"
  "strconv"
  "strings"
  "time"
)

type packageSupervisors struct {
  PackageName string,
  Start func(),
  Stop func()
}

type endpointHandler struct {
  Path string,
  Handler func()
}

const (
  toRegisterWithSupervisor = []packageSupervisors{
    packageSupervisors{hash.PackageName, hash.Start, hash.Stop},
    packageSupervisors{identifier.PackageName, identifier.Start, identifier.Stop},
    packageSupervisors{registrar.PackageName, registrar.Start, registrar.Stop},
    packageSupervisors{stats.PackageName, stats.Start, stats.Stop},
    packageSupervisors{repository.PackageName, repository.Start, repository.Stop}
  }

  hashPath := endpointHandler{'/hash', hashHandler}
  hashedPath := endpointHandler{'/hash/*', hashedHandler}
  shutdownPath := endpointHandler{'/shutdown', shutdownHandler}
  statsPath := endpointHandler{'/stats', statsHandler}

  toRegisterWithHttp = []endpointHandler{
    hashPath,
    hashedPath,
    shutdownPath,
    statsPath
  }
)

func hashHandler(w http.ResponseWriter, r *http.Request, ch chan) {
  switch r.Method {
  case "POST":
    password := r.FormValue("password")
    if password == "" { // TODO: Work with PM to define password requirements to match off of.
      w.WriteHeader(400)
    } else {
      id := identifier.Get()
      w.Write([]byte(string(id)))
      hashed := hash.Hash(password)
      now := time.Now()
      registrar.Put(id, hashed, now.Add(time.Second * 5))
    }
  default:
    w.WriteHeader(405) // Method not allowed.
  }
  ch<-
}

func hashedHandler(w http.ResponseWriter, r *http.Request, ch chan) {
  switch r.Method {
  case "GET":
    path := r.URL.Path
    pathSplit := strings.Split(path, "/")

    id, _ := strconv.Atoi(pathSplit[len(pathSplit)])
    record, err := registrar.Get(id)

    if (err == nil) && (record.hashtime >= time.Now()) {
      w.Write([]byte(record.hash))
    } else {
      w.WriteHeader(404)
    }
  default:
    w.WriteHeader(405) // Method not allowed.
  }
  ch<-
}

func shutdownHandler (w http.ResponseWriter, r *http.Request, ch chan) {
  // TODO: Better define which method should be used to issue a shutdown request.
  w.WriteHeader(204)
  Stop()
  ch<-
}

func statsHandler (w http.ResponseWriter, r *http.Request, ch chan) {
  switch r.Method {
  case "GET":
    response := make[string]int
    total, average := stats.Get()

    response["total"] = total
    response["average"] = average

    jsonResponse, err := json.Marshall(response)

    if err == nil {
      w.Header.Set("Content-Type", "application/json")
      w.Write(jsonResponse)
    } else {
      w.WriteHeader(500)
    }
  default:
    w.WriteHeader(405)
  }
  ch<-
}

func Start() {
  for pkg := range toRegisterWithSupervisor {
    supervisor.Register(pkg.PackageName, pkg.Start, pkg.Stop)
  }
  supervisor.StartAll()

  for endpoint := range toRegisterWithHttp {
    httpserver.Register(endpoint.Path, endpoint.Handler)
  }

  httpserver.Start("localhost", "8080")
}

func Stop() {
  httpserver.Stop()
  supervisor.StopAll()
}
