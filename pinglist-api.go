package main

import (
	"log"
	"os"

	"github.com/RichardKnop/pinglist-api/cmd"
	"github.com/urfave/cli"
)

var (
	cliApp *cli.App
)

func init() {
	// Initialise a CLI app
	cliApp = cli.NewApp()
	cliApp.Name = "Pinglist"
	cliApp.Usage = "pinglist-api"
	cliApp.Author = "Richard Knop"
	cliApp.Email = "risoknop@gmail.com"
	cliApp.Version = "0.0.0"
}

func main() {
	// Set the CLI app commands
	cliApp.Commands = []cli.Command{
		{
			Name:  "migrate",
			Usage: "run migrations",
			Action: func(c *cli.Context) error {
				return cmd.Migrate()
			},
		},
		{
			Name:  "loaddata",
			Usage: "load data from fixture",
			Action: func(c *cli.Context) error {
				return cmd.LoadData(c.Args())
			},
		},
		{
			Name:  "createaccount",
			Usage: "create new account",
			Action: func(c *cli.Context) error {
				return cmd.CreateAccount()
			},
		},
		{
			Name:  "createsuperuser",
			Usage: "create new superuser",
			Action: func(c *cli.Context) error {
				return cmd.CreateSuperuser()
			},
		},
		{
			Name:  "runscheduler",
			Usage: "run scheduler",
			Action: func(c *cli.Context) error {
				return cmd.RunScheduler()
			},
		},
		{
			Name:  "runserver",
			Usage: "run web server",
			Action: func(c *cli.Context) error {
				return cmd.RunServer()
			},
		},
		{
			Name:  "runall",
			Usage: "run both scheduler and web server",
			Action: func(c *cli.Context) error {
				return cmd.RunAll()
			},
		},
	}

	// Run the CLI app
	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
