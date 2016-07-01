package main

import (
	"gopkg.in/urfave/cli.v2"
	"commands"
	"os"
)

func main ()  {
	app := &cli.App{
		Name: "Qbyco Backup",
		Version: "0.1.0 alpha",
		Usage: "Backup and be safe :)",
		EnableBashCompletion: true,
		Commands: []*cli.Command {},
	}

	app.Commands = append(app.Commands, commands.InitCommand())
	app.Commands = append(app.Commands, commands.BackupCommand())

	app.Run(os.Args)
}