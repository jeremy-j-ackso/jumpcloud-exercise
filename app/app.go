// App Package brings all the services and features together to provide a single point of entry for main.
package app

import (
  "encoding/json"
  "jumpcloudexercise/features/hash"
  "jumpcloudexercise/features/identifier"
  "jumpcloudexercise/features/registrar"
  "jumpcloudexercise/features/stats"
  "jumpcloudexercise/services/httpserver"
  "jumpcloudexercise/services/repository"
  "jumpcloudexercise/services/supervisor"
  "net/http"
  "strconv"
  "strings"
  "time"
)

type packageSupervisors struct {
  PackageName string
  Start func()
  Stop func(func(string))
}

type endpointHandler struct {
  Path string
  Handler func(http.ResponseWriter, *http.Request, chan int)
}

var (
  toRegisterWithSupervisor = [...]packageSupervisors{
    packageSupervisors{hash.PackageName, hash.Start, hash.Stop},
    packageSupervisors{identifier.PackageName, identifier.Start, identifier.Stop},
    packageSupervisors{registrar.PackageName, registrar.Start, registrar.Stop},
    packageSupervisors{stats.PackageName, stats.Start, stats.Stop},
    packageSupervisors{repository.PackageName, repository.Start, repository.Stop},
  }

  hashPath = endpointHandler{"/hash", hashHandler}
  hashedPath = endpointHandler{"/hash/*", hashedHandler}
  shutdownPath = endpointHandler{"/shutdown", shutdownHandler}
  statsPath = endpointHandler{"/stats", statsHandler}

  toRegisterWithHttp = []endpointHandler{
    hashPath,
    hashedPath,
    shutdownPath,
    statsPath,
  }
)

func hashHandler(w http.ResponseWriter, r *http.Request, ch chan int) {
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
  ch<- 1
}

func hashedHandler(w http.ResponseWriter, r *http.Request, ch chan int) {
  switch r.Method {
  case "GET":
    path := r.URL.Path
    pathSplit := strings.Split(path, "/")

    id, _ := strconv.Atoi(pathSplit[len(pathSplit)])
    record, err := registrar.Get(id)

    if (err == nil) && record.Hashtime.Before(time.Now()) {
      w.Write([]byte(record.Hash))
    } else {
      w.WriteHeader(404)
    }
  default:
    w.WriteHeader(405) // Method not allowed.
  }
  ch<- 1
}

func shutdownHandler (w http.ResponseWriter, r *http.Request, ch chan int) {
  // TODO: Better define which method should be used to issue a shutdown request.
  w.WriteHeader(204)
  Stop()
  ch<- 1
}

func statsHandler (w http.ResponseWriter, r *http.Request, ch chan int) {
  switch r.Method {
  case "GET":
    response := make(map[string]int64)
    total, average := stats.Get()

    response["total"] = int64(total)
    response["average"] = average

    jsonResponse, err := json.Marshal(response)

    if err == nil {
      w.Header().Set("Content-Type", "application/json")
      w.Write(jsonResponse)
    } else {
      w.WriteHeader(500)
    }
  default:
    w.WriteHeader(405)
  }
  ch<- 1
}

func Start() {
  for pkg := range toRegisterWithSupervisor {
    supervisor.Register(
      toRegisterWithSupervisor[pkg].PackageName,
      toRegisterWithSupervisor[pkg].Start,
      toRegisterWithSupervisor[pkg].Stop,
    )
  }
  supervisor.StartAll()

  for endpoint := range toRegisterWithHttp {
    httpserver.Register(toRegisterWithHttp[endpoint].Path, toRegisterWithHttp[endpoint].Handler)
  }

  httpserver.Start("localhost", "8080")
}

func Stop() {
  httpserver.Stop()
  supervisor.StopAll()
}
