package clihandler

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

// New return instnce of cli application
func New() *cli.App {
	cliApp := cli.NewApp()
	if err := cliApp.Run(os.Args); err != nil {
		log.Fatalf("Failed to start application : %v", err)
	}

	return cliApp
}
