package main

import (
	"github.com/urfave/cli"
	"os"
	"fmt"
	"github.com/ntfs32/gvm/src/gvm/action"
)

var (
	Revision = "1.0.0"
)


func init()  {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "gvm Version: %s\n", c.App.Version)
	}
}

func main()  {
	app := cli.NewApp()
	app.Name = "gvm"
	app.Usage = "Golang Version Manager"
	app.Description = "Simple bash script to manage multiple active golang versions"
	app.Version = Revision

	app.Commands = []cli.Command{
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "Install version `[version]`",
			Action:  action.Install,
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list Version had installed local",
			Action:  action.ListLocal,
		},
		{
			Name:    "list-remote",
			Aliases: []string{"lr"},
			Usage:   "list remote golang release version",
			Action:  action.ListRemote,
		},
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.Writer, "There is no %q here.\n", command)
	}
	app.Run(os.Args)
}