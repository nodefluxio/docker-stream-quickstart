# Centralize Enrollment System (CES) - Coordinator

## Description

CES coordinator is part of Centralize Enrollment System, this service will handle and save all event enrollment for serve to CES agent .

## Prerequisites

- Golang >= version 1.14
- PosgreeSQL >= version 12
- CES Agent = version latest ( internal service )

## How To Run

### Development

1. install all depedency

```bash
$ go mod download
```

2. run with go

```bash
$ go run cmd/cescoordinator/main.go
```

### Docker

1. install docker
2. run this command

```bash
$ docker run -it -d --net host --restart unless-stopped --name ces-coordinator registry.gitlab.com/nodefluxio/vanilla-dashboard:latest ./cescoordinator
```

## Config Options

```bash
NAME:
   CES coordinator - Centralize Enrollment System Coordinaotr

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log value, -l value                logging level, useful for debugging session. available (warning, info, debug) (default: "info") [$LOG_LEVEL]
   --port value, -p value               app http port (default: "6012") [$PORT]
   --node-env value                     running mode, development or production (default: "production") [$NODE_ENV]
   --db-host value, --dbhost value      postgreSQL host address, ex: 127.0.0.1:5432 (default: "127.0.0.1:5432") [$DB_HOST]
   --db-user value, --dbusr value       postgreSQL database username, ex: postgres (default: "postgres") [$DB_USERNAME]
   --db-password value, --dbpass value  postgreSQL database password, ex: password (default: "test") [$DB_PASSWORD]
   --db-name value, --dbnm value        postgreSQL database name, ex: postgres (default: "ces_coordinator") [$DB_NAME]
   --cron-partition value, --crp value  database cronjob partition, ex: 0 0 * * * (default: "0 0 * * *") [$CRONJOB_PARTITION_SPEC]
   --help, -h                           show help
   --version, -v                        print the version

```
