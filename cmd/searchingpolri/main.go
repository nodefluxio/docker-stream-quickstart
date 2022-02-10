package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
	"github.com/urfave/cli"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	httpserver "gitlab.com/nodefluxio/vanilla-dashboard/cmd/searchingpolri/http"
	searchingpolri "gitlab.com/nodefluxio/vanilla-dashboard/cmd/searchingpolri/http"
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
		Value:  "6014",
		Usage:  "app http port",
		EnvVar: "PORT",
	},
	cli.StringFlag{
		Name:   "polri-base-url",
		Value:  "https://api.polri.go.id",
		Usage:  "base url porli api",
		EnvVar: "POLRI_BASE_URL",
	},
	cli.StringFlag{
		Name:   "polri-username",
		Value:  "admin",
		Usage:  "authorization basic username",
		EnvVar: "POLRI_USERNAME",
	},
	cli.StringFlag{
		Name:   "polri-password",
		Value:  "admin",
		Usage:  "authorization basic password",
		EnvVar: "POLRI_PASSWORD",
	},
	cli.StringFlag{
		Name:   "polri-plate-path",
		Value:  "/korlantas/getbynopol",
		Usage:  "path for api search by plate number / nopol",
		EnvVar: "POLRI_PLATE_PATH",
	},
	cli.StringFlag{
		Name:   "polri-nik-path",
		Value:  "/dukcapil/getbynik2",
		Usage:  "path for api search by nik",
		EnvVar: "POLRI_NIK_PATH",
	},
	cli.StringFlag{
		Name:   "seagate-base-url",
		Value:  "http://localhost:3003",
		Usage:  "base url seagate api",
		EnvVar: "SEAGATE_BASE_URL",
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
		Name:   "max-size-img, msi",
		Value:  "2097152",
		Usage:  "maximal size image face search in bytes, ex: 2097152 -> 2MB",
		EnvVar: "MAX_SIZE_IMG",
	},
}

func action(c *cli.Context) {
	options := searchingpolri.Options{
		LogLevel:             c.String("log"),
		AppPort:              c.String("port"),
		PolriBaseURL:         c.String("polri-base-url"),
		PolriBasicUsername:   c.String("polri-username"),
		PolriBasicPassword:   c.String("polri-password"),
		PolriSearchPlatePath: c.String("polri-plate-path"),
		PolriSearchNikPath:   c.String("polri-nik-path"),
		SeagateBaseURL:       c.String("seagate-base-url"),
		FremisnURL:           c.String("fremisn-url"),
		FremisKeyspace:       c.String("fremisn-keyspace"),
		MaxImageSize:         c.String("max-size-img"),
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
	app.Name = "Searching Polri"
	app.Usage = "service for handling all searching for polri"
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
