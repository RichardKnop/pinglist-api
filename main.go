package main

import (
	"log"
	"os"

	"github.com/RichardKnop/pinglist-api/commands"
	"github.com/codegangsta/cli"
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
			Action: func(c *cli.Context) {
				if err := commands.Migrate(); err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:  "loaddata",
			Usage: "load data from fixture",
			Action: func(c *cli.Context) {
				if err := commands.LoadData(c.Args()); err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:  "createaccount",
			Usage: "create new account",
			Action: func(c *cli.Context) {
				if err := commands.CreateAccount(); err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:  "createsuperuser",
			Usage: "create new superuser",
			Action: func(c *cli.Context) {
				if err := commands.CreateSuperuser(); err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:  "runscheduler",
			Usage: "run scheduler",
			Action: func(c *cli.Context) {
				if err := commands.RunScheduler(); err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:  "runserver",
			Usage: "run web server",
			Action: func(c *cli.Context) {
				if err := commands.RunServer(); err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:  "runall",
			Usage: "run both scheduler and web server",
			Action: func(c *cli.Context) {
				if err := commands.RunAll(); err != nil {
					log.Fatal(err)
				}
			},
		},
	}

	// Run the CLI app
	cliApp.Run(os.Args)
}
