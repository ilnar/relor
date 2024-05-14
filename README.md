> ⚠️ 
> This is a prototype and can't be used in production, feel free to contribute or subscribe

# Relentless Orchestrator

Relentless Orchestrator helps to automate business-critical processes that require a high completion rate and are subject to strict SLAs.

The orchestrator simplifies the use of cycles in the execution graph which allows to model complex error handling logic.

*For example, certain operation might fail and be retried three times and then a human intervention might be requested.
The latter then might be re-request after a timeout.
This can go further and trigger incident creation if the failure hasn't been addressed after several intervention requests.*

## Features

1. Explicit, persistent workflow execution graph supporting cycles
2. Operation timeouts and retries
3. Idempotency

## Build and Run

### Fast Build

Fast build will not regenerate any files and use already generated go code in `gen/` folder.

```sh
make build
```

### Full Build

Full build will regenerate all files in `gen/` folder. This is needed when changes to APIs or DB Schema have been made.

This will require installing `protoc` and `sqlc`. 

```sh
make
```

### Local Run

Local run requires `docker` and `golang-migrate`.

```sh 
make startpg
./bin/server
```

Example worker then can be started too:

```sh
./bin/example_worker

```