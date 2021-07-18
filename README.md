# JumpCloud Coding Exercise

## Plan

- External Boundary is HTTP.
- Handlers should call an appropriate interface method that engages lower level module.
- Each of the different functionalities should have a package to handle it.
  - Identifier Package - Feature
    - `get`: Increments internal identifier number and return its value.
    - `current`: Returns the internal identifier number.
    - `start`: Sets up the package by registering as up with supervisor.
    - `stop`: Waits until all work is done, then registers as down with Supervisor.
  - Registrar Package - Feature
    - `put`: Takes 3 arguments: `id`, `now() + 5s`, `hash`
      - Stores the argument values in the Repository.
    - `get`: Takes 1 argument: `id`.
      - If `now()` is greater than or equal to stored timestamp for that `id` it returns.
      - Else rejects.
    - `start`: Sets up the package by registering as up with supervisor.
    - `stop`: Waits until all work is done, then registers as down with Supervisor.
  - Hash Package - Feature
    - `hash`: Takes 1 argument: `password`
      - Returns the hash and updates the hashtime average.
    - `getAvg`: Returns the hashtime average.
    - `start`: Sets up the package by registering as up with supervisor.
    - `stop`: Waits until all work is done, then registers as down with Supervisor.
  - Hashtime Package - Feature
    - `put`: Takes 1 argument: `time`
      - Adds the new `time` to a total stored in seconds.
    - `get`: Returns the current sum of times.
    - `start`: Sets up the package by registering as up with supervisor.
    - `stop`: Waits until all work is done, then registers as down with Supervisor.
  - Stats Package - Feature
    - `get`: Retrieves values from `Identifier.current` and `Hashtime.get` and returns the current Identifier value and `hashtime / identifier`.
    - `start`: Sets up the package by registering as up with supervisor.
    - `stop`: Waits until all work is done, then registers as down with Supervisor.
  - Supervisor Package - Service
    - `start`: launches all services
    - `register`: services register their up status with the Supervisor.
    - `unregister`: services register their down status with the Supervisor.
    - `stop`: requests all services to finish up current work and update their up/down status.
  - HTTP package - Service
    - `register`: Takes one arguments: `path`, `handler`
      - Registers a handler with a path.
    - `start`: Takes 2 arguments: `address`, `port`
      - Starts listening at specified address:port.
      - Registers as up with supervisor.
    - `stop`: Stops listening, waits until all work is done, then registers as down with Supervisor.
  - Repository Package - Services
    - `get`: Query by ID
    - `put`: Upsert by ID
    - `start`: Sets up the database and registers services as up with Supervisor
    - `stop`: Waits until all work is done, then registers as down with Supervisor.

## Optimistic Plan

- Include CI/CD to build a container.
- Include automated unit testing.
- Include good logging. (Ideally to a file in a directory mounted to a container, possibly include a config to send to a remote logging system.)
