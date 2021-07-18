// Starts a supervisor to start services and features, and monitor for safe shutdown conditions.
package supervisor

var (
  register := make(map[string]Package)
)

type Package struct {
  Start func(),
  Stop func()
}

func Start(pkg) {
  register[pkg].Start()
}

func StartAll() {
  for pkg := range register {
    go Start(pkg)
  }
}

func Stop(pkg) {
  register[pkg].Stop()
}

func StopAll() {
  for pkg := range register {
    go Stop(pkg)
  }
}

func Register(pkg, stop) {
  register[pkg] := Package{
    Start: start,
    Stop: stop
  }
}

func Unregister(pkg) {
  delete(register, pkg)
}
