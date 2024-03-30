# Workflow Orchestrator

A simple cloud-agnostic workflow orchestrator. It accepts a directed graph and executes it by sending jobs to a worker while traversing the graph.

## Features

1. Explicit workflow graph
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