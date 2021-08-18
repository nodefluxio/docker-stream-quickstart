package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
	"github.com/urfave/cli"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	httpserver "gitlab.com/nodefluxio/vanilla-dashboard/cmd/vanend/http"
	vanend "gitlab.com/nodefluxio/vanilla-dashboard/cmd/vanend/http"
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
		Value:  "80",
		Usage:  "app http port",
		EnvVar: "PORT",
	},
	cli.StringFlag{
		Name:   "node-env",
		Value:  "production",
		Usage:  "running mode, development or production",
		EnvVar: "NODE_ENV",
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
		Value:  "postgres",
		Usage:  "postgreSQL database name, ex: postgres",
		EnvVar: "DB_NAME",
	},
	cli.StringFlag{
		Name:   "fremisn-url, frnurl",
		Value:  "http://localhost:4005/v1/face",
		Usage:  "fremisn service url, ex: http://localhost:4005/v1/face",
		EnvVar: "FREMISN_URL",
	},
	cli.StringFlag{
		Name:   "fremisn-keyspace, frnks",
		Value:  "some-keyspace",
		Usage:  "fremisn user keyspace, ex: some-keyspace",
		EnvVar: "FREMISN_KEYSPACE",
	},
	cli.StringFlag{
		Name:   "visionaire-host, vshst",
		Value:  "localhost:4004",
		Usage:  "visionaire docker stream host, ex: localhost:4004, 127.0.0.1:4004",
		EnvVar: "VISIONAIRE_HOST",
	},
	cli.StringFlag{
		Name:   "cron-partition, crp",
		Value:  "0 0 * * *",
		Usage:  "database cronjob partition, ex: 0 0 * * *",
		EnvVar: "CRONJOB_PARTITION_SPEC",
	},
	cli.StringFlag{
		Name:   "website-host",
		Value:  "localhost",
		Usage:  "front end website host",
		EnvVar: "REACT_APP_HOST",
	},
	cli.StringFlag{
		Name:   "use-ces",
		Value:  "false",
		Usage:  "activate CES (Centralize Enrollment System)",
		EnvVar: "USER_CES",
	},
	cli.StringFlag{
		Name:   "agent-url",
		Value:  "http://localhost:6013",
		Usage:  "CES agent service url, ex: http://localhost:6013",
		EnvVar: "AGENT_URL",
	},
	cli.StringFlag{
		Name:   "max-size-img-enrollment, msie",
		Value:  "2097152",
		Usage:  "maximal size image face enrollment in bytes, ex: 2097152 -> 2MB",
		EnvVar: "MAX_SIZE_IMG_ENROLLMENT",
	},
}

func action(c *cli.Context) {
	options := vanend.Options{
		LogLevel:               c.String("log"),
		FEAppHost:              c.String("website-host"),
		AppPort:                c.String("port"),
		NodeEnv:                c.String("node-env"),
		DatabaseHost:           c.String("db-host"),
		DatabaseUsername:       c.String("db-user"),
		DatabasePassword:       c.String("db-password"),
		DatabaseName:           c.String("db-name"),
		FremisnURL:             c.String("fremisn-url"),
		FremisKeyspace:         c.String("fremisn-keyspace"),
		VisionaireHost:         c.String("visionaire-host"),
		CronjobPartitionSpec:   c.String("cron-partition"),
		UseCES:                 c.String("use-ces"),
		AgentURL:               c.String("agent-url"),
		MaxSizeImageEnrollment: c.String("max-size-img-enrollment"),
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
	app.Name = "vanend"
	app.Usage = "vanilla dashboard for visionaire v4"
	app.Version = "0.13.1"
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
