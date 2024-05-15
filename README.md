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

## Overview

### Key Concepts

**Workflow** models a real-world process as a series of operations. It combines an execution graph, the current position in this graph and some metadata.

**Operation** is the minimal meaningful building block of the workflow. Operations are expected to be atomic and idempotent. The operation usually updates the state in various systems the workflow relies upon. Operations define workflow behaviour and can be perceived as a function declaration.

**Job** represents an attempt to perform an operation with given inputs. Jobs are not reusable and are short-lived. A job can be compared with a function call.

**Worker** pulls jobs and performs them. Worker implements operations. Following the function example, the worker holds the definition of a function, while the operation declares that function and job triggers a call to that function.

**Label** is the output of a job. Each operation is offered a set of labels to choose from. It is up to the implementation of the operation to choose a single resulting label for a given job. The resulting label will control the progression of the workflow.

**Execution Graph** or Graph is a directed graph that has operations in its nodes and labels on its edges. When the operation is completed the resulting label is used to choose the next operation. The workflow is completed when there are no operations left, i.e., the execution reaches the leaf node in the graph.