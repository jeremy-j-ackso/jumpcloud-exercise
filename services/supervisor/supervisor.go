// Starts a supervisor to start services and features, and monitor for safe shutdown conditions.
package supervisor

var (
  register = make(map[string]Package)
)

type Package struct {
  Start func()
  Stop func(func(string))
}

func Register(pkg string, start func(), stop func(func(string))) {
  register[pkg] = Package{start, stop}
}

func Unregister(pkg string) {
  delete(register, pkg)
}

func Start(pkg string) {
  register[pkg].Start()
}

func StartAll() {
  for pkg := range register {
    go Start(pkg)
  }
}

func Stop(pkg string) {
  register[pkg].Stop(Unregister)
}

func StopAll() {
  for pkg := range register {
    go Stop(pkg)
  }
}
