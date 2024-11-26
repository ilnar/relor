[![codecov](https://codecov.io/gh/gemlab-dev/relor/graph/badge.svg?token=Y6XRNVMZ1N)](https://codecov.io/gh/gemlab-dev/relor)
> ⚠️ 
> This is a prototype and can't be used in production; feel free to contribute or subscribe

# Relentless Orchestrator

Relentless Orchestrator helps to automate long-running business-critical processes that require a high completion rate and are subject to strict SLAs.

- It safeguards long-running processes by relying on **static workflow definitions**. This helps deploy code frequently and with confidence.
- It offers robust error handling by supporting **cycles in workflows**. It allows modelling retries, escalations, and automatic and manual mitigation procedures.
- It executes workflow steps **at least once** in the order they appear in the workflow definition.

In addition, Relentless Orchestrator consists of isolated services that can be scaled separately to optimise resource utilisation. 



## Build and Run

### Fast Build

Fast build will not regenerate any files and use already generated go code in `gen/` folder.

```sh
make build
```

### Full Build

A full build will regenerate all files in the `gen/` folder. This is needed when changes to APIs or DB Schema have been made.

This will pull and build Docker images used for code generation.

```sh
make
```

### Local Run

Local run requires `docker`.

Single-node orchestrator:

```sh 
make startpg
./bin/server --config=dev/single-node.json
```

Example worker pulling jobs:

```sh
./bin/example_worker

```

Multi-node cluster:

```sh
make startpg
./bin/server --config=dev/cluster-node-1.json
./bin/server --config=dev/cluster-node-2.json
./bin/server --config=dev/cluster-node-3.json
```

## Overview

### Key Concepts

**Workflow** models a real-world process as a series of operations. It combines a graph definition, the current position in this graph and some metadata.

**Operation** is the minimal meaningful building block of the workflow. It is expected to be atomic and idempotent. Operations usually update the state in various systems on which the workflow relies. They define workflow behaviour and can be perceived as function declarations.

**Job** represents an attempt to perform an operation with given inputs. It is not reusable and is short-lived. A job can be compared with a function call.

**Worker** pulls jobs and performs them. The worker implements operations. Following the function example, the worker holds the definition of a function while the operation declares that function and the job triggers a call to that function.

**Label** is the output of a job. Each operation is offered a set of labels to choose from. The implementation of the operation selects a single resulting label for a given job. The resulting labels control the progression of the workflow.

**Graph Definition** is a directed graph with operations in its nodes and labels on its edges. When an operation is completed, the resulting label is used to choose the next operation. The workflow is completed when the execution reaches a leaf node in the graph.
