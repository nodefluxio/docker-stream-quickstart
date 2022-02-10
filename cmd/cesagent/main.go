package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
	"github.com/urfave/cli"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	cesagent "gitlab.com/nodefluxio/vanilla-dashboard/cmd/cesagent/http"
	httpserver "gitlab.com/nodefluxio/vanilla-dashboard/cmd/cesagent/http"
	validatorHelper "gitlab.com/nodefluxio/vanilla-dashboard/internal/delivery/http/middleware/validator/helper"
	"gopkg.in/go-playground/validator.v9"
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name:   "log, l",
		Value:  "info",
		Usage:  "logging level, useful for debugging session. available (warning, info, debug)",
		EnvVar: "LOG_LEVEL",
	},
	cli.StringFlag{
		Name:   "port, p",
		Value:  "6013",
		Usage:  "app http port",
		EnvVar: "PORT",
	},
	cli.StringFlag{
		Name:   "db-host, dbhost",
		Value:  "127.0.0.1:5432",
		Usage:  "postgreSQL host address, ex: 127.0.0.1:5432",
		EnvVar: "DB_HOST",
	},
	cli.StringFlag{
		Name:   "db-user, dbusr",
		Value:  "postgres",
		Usage:  "postgreSQL database username, ex: postgres",
		EnvVar: "DB_USERNAME",
	},
	cli.StringFlag{
		Name:   "db-password, dbpass",
		Value:  "test",
		Usage:  "postgreSQL database password, ex: password",
		EnvVar: "DB_PASSWORD",
	},
	cli.StringFlag{
		Name:   "db-name, dbnm",
		Value:  "ces_agent",
		Usage:  "postgreSQL database name, ex: postgres",
		EnvVar: "DB_NAME",
	},
	cli.StringFlag{
		Name:   "db-max-idle, dbmi",
		Value:  "4",
		Usage:  "postgreSQL database max idle connection",
		EnvVar: "DB_MAX_IDLE_CONNECTION",
	},
	cli.StringFlag{
		Name:   "db-max-open, dbmc",
		Value:  "8",
		Usage:  "postgreSQL database max open connection",
		EnvVar: "DB_MAX_OPEN_CONNECTION",
	},
	cli.StringFlag{
		Name:   "db-max-lifetime, dbmt",
		Value:  "5",
		Usage:  "postgreSQL database max connection lifetime in minute",
		EnvVar: "DB_CONNECTION_MAX_LIFETIME_IN_MINUTE",
	},
	cli.StringFlag{
		Name:   "sync-period, sp",
		Value:  "0 */1 * * *",
		Usage:  "cornjob format for set syncronize , ex: 0 */1 * * *",
		EnvVar: "SYNC_PERIOD",
	},
	cli.StringFlag{
		Name:   "agent-name, a",
		Usage:  "agent name",
		EnvVar: "AGENT_NAME",
	},
	cli.StringFlag{
		Name:   "coor-url, cu",
		Value:  "http://localhost:6012",
		Usage:  "CES coordinator REST API URL",
		EnvVar: "COORDINATOR_URL",
	},
	cli.StringFlag{
		Name:   "enrollment-vanilla-url, ecu",
		Value:  "http://localhost:6010",
		Usage:  "vanilla dashboard face enrollment REST API URL",
		EnvVar: "ENROLLMENT_VANILLA_URL",
	},
	cli.StringFlag{
		Name:   "total-event-sync, tes",
		Value:  "10",
		Usage:  "total event enrollment to sync",
		EnvVar: "TOTAL_EVENT_SYNC",
	},
}

func action(c *cli.Context) {
	options := cesagent.Options{
		LogLevel:                    c.String("log"),
		AppPort:                     c.String("port"),
		DatabaseHost:                c.String("db-host"),
		DatabaseUsername:            c.String("db-user"),
		DatabasePassword:            c.String("db-password"),
		DatabaseName:                c.String("db-name"),
		DatabaseMaxIdleConn:         c.String("db-max-idle"),
		DatabaseMaxOpenConn:         c.String("db-max-open"),
		DatabaseMaxLifetimeInMinute: c.String("db-max-lifetime"),
		SyncPeriod:                  c.String("sync-period"),
		AgentName:                   c.String("agent-name"),
		CoordinatorURL:              c.String("coor-url"),
		EnrollmentVanillaURL:        c.String("enrollment-vanilla-url"),
		TotalEventSync:              c.String("total-event-sync"),
	}
	logutil.Init(options.LogLevel)
	validate := validator.New()
	i18n := validatorHelper.RequiredErrorMessage(validate)
	err := validate.Struct(options)
	if err != nil {
		logutil.LogObj.SetFatalLog(gin.H{
			"errors": validatorHelper.ErrorMessageTranslator(err, i18n),
		}, "App cannot start, the arguments is not complete or invalid")
	}
	logutil.LogObj.SetInfoLog(map[string]interface{}{
		"app_port":  c.String("port"),
		"log_level": c.String("log"),
	}, "Application start successfully!")

	httpserver.Start(&options)
}

func main() {
	gotenv.Load()
	app := cli.NewApp()
	app.Name = "CES agent"
	app.Usage = "Centralize Enrollment System Agent"
	app.Version = "0.1.0"
	app.Flags = flags
	app.Action = action

	err := app.Run(os.Args)
	if err != nil {
		logutil.Init("info")
		logutil.LogObj.SetPanicLog(map[string]interface{}{
			"error": err,
		}, "Application failed to start")
	}
}
