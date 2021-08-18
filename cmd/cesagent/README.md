# Centralize Enrollment System (CES) - Agent

## Description

CES agent is part of Centralize Enrollment System, this service will handle sync proses at enrollment system from coordinator and vanilla dashboard.

## Prerequisites

- Golang >= version 1.14
- PosgreeSQL >= version 12
- vanend = version latest ( internal service )
- CES Coordinator = version latest ( internal service )

## How To Run

### Development

1. install all depedency

```bash
$ go mod download
```

2. run with go

```bash
$ go run cmd/cesagent/main.go
```

### Docker

1. install docker
2. run this command

```bash
$ docker run -it -d --net host --restart unless-stopped --name ces-agent registry.gitlab.com/nodefluxio/vanilla-dashboard:latest ./cesagent --agent-name nodeflux
```

## Config Options

```bash
NAME:
   CES agent - Centralize Enrollment System Agent

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log value, -l value                        logging level, useful for debugging session. available (warning, info, debug) (default: "info") [$LOG_LEVEL]
   --port value, -p value                       app http port (default: "6013") [$PORT]
   --db-host value, --dbhost value              postgreSQL host address, ex: 127.0.0.1:5432 (default: "127.0.0.1:5432") [$DB_HOST]
   --db-user value, --dbusr value               postgreSQL database username, ex: postgres (default: "postgres") [$DB_USERNAME]
   --db-password value, --dbpass value          postgreSQL database password, ex: password (default: "test") [$DB_PASSWORD]
   --db-name value, --dbnm value                postgreSQL database name, ex: postgres (default: "ces_agent") [$DB_NAME]
   --sync-period value, --sp value              cornjob format for set syncronize , ex: 0 */1 * * * (default: "0 */1 * * *") [$SYNC_PERIOD]
   --agent-name value, -a value                 agent name [$AGENT_NAME]
   --coor-url value, --cu value                 CES coordinator REST API URL (default: "http://localhost:6012") [$COORDINATOR_URL]
   --enrollment-vanilla-url value, --ecu value  vanilla dashboard face enrollment REST API URL (default: "http://localhost:6010") [$ENROLLMENT_VANILLA_URL]
   --total-event-sync value, --tes value        total event enrollment to sync (default: "10") [$TOTAL_EVENT_SYNC]
   --help, -h                                   show help
   --version, -v                                print the version

```
