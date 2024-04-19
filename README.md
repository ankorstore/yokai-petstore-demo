# Yokai Petstore Demo

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go version](https://img.shields.io/badge/Go-1.22-blue)](https://go.dev/)

> HTTP application demo based on the [Yokai](https://github.com/ankorstore/yokai) Go framework.

<!-- TOC -->
* [Overview](#overview)
  * [Layout](#layout)
  * [Makefile](#makefile)
* [Usage](#usage)
<!-- TOC -->

## Overview

This demo provides a REST API example, with:

- a [Yokai](https://github.com/ankorstore/yokai) application container (with [Air](https://github.com/cosmtrek/air)):
  - with the [HTTP server](https://ankorstore.github.io/yokai/modules/fxhttpserver/) module
  - with a [fxdatabase](internal/module/fxdatabase) module, to provide `database/sql` + [gomigrate](https://github.com/golang-migrate/migrate)
  - with a [fxsqlc](internal/module/fxsqlc) module,  to provide [sqlc](https://github.com/sqlc-dev/sqlc)
- a [MySQL](https://www.mysql.com/) container for storage
- a [Jaeger](https://www.jaegertracing.io/) container for tracing

### Layout

This demo is following the [recommended project layout](https://go.dev/doc/modules/layout#server-project):

- `cmd/`: entry points
- `configs/`: configuration files
- `db/`
  - `migrations/`: SQL migrations files for gomigrate
  - `queries/`: SQL queries files for sqlc stubs generation
  - `sqlc/`: sqlc generated stubs
- `internal/`:
  - `handler/`: HTTP handlers
  - `module/`: 
    - `fxdatabase`: database/sql + gomigrate module
    - `fxsqlc`: sqlc module
  - `bootstrap.go`: bootstrap
  - `register.go`: dependencies registration
  - `router.go`: routing registration
- `sqlc.yaml`: sqlc configuration

### Makefile

This template provides a [Makefile](Makefile):

```
make up              # start the docker compose stack
make down            # stop the docker compose stack
make logs            # stream the docker compose stack logs
make fresh           # refresh the docker compose stack
make sqlc            # regenerate sqlc stubs
make migrate-create  # create migration
make migrate-up      # apply all migrations 
make migrate-down    # revert all migrations 
make test            # run tests
make lint            # run linter
```

## Usage

Start the application with:

```shell
make fresh
```

Then apply database migrations with:

```shell
make migrate-up
```

Once ready, the application will be available on:
- [http://localhost:8080](http://localhost:8080) for the application HTTP server
- [http://localhost:8081](http://localhost:8081) for the application core dashboard
- [http://localhost:16686](http://localhost:16686) for the Jaeger dashboard

Available endpoints on [http://localhost:8080](http://localhost:8080):

| Route                   | Description     |
|-------------------------|-----------------|
| `[GET] /owners`         | List all owners |
| `[POST] /owners`        | Create an owner |
| `[GET] /owners/:id`     | Get an owner    |
| `[DELETE] /gophers/:id` | Delete an owner |