package main

import (
	"fmt"

	"github.com/gchaincl/go-etesync/api"
	"github.com/urfave/cli"
)

type App struct {
	cli    *cli.App
	client *api.Client
}

func NewApp() *App {
	app := &App{}

	app.cli = cli.NewApp()
	app.cli.Name = "etecli"
	app.cli.Usage = "ETESync cli tool"
	app.cli.Flags = []cli.Flag{
		cli.StringFlag{Name: "email"},
		cli.StringFlag{Name: "password"},
		cli.StringFlag{Name: "key", Usage: "Encryption key"},
	}
	app.cli.Before = app.Before

	app.cli.Commands = []cli.Command{
		app.JournalsCmd(),
	}

	return app
}

func (app *App) Before(ctx *cli.Context) error {
	client, err := api.NewClient(ctx.GlobalString("email"), ctx.GlobalString("password"))
	if err != nil {
		return err
	}

	app.client = client
	return nil
}

func (app *App) JournalsCmd() cli.Command {
	return cli.Command{
		Name:  "journals",
		Usage: "Display available journals",
		Action: func(ctx *cli.Context) error {
			js, err := app.client.Journals()
			if err != nil {
				return err
			}

			for _, j := range js {
				content, err := j.GetContent([]byte(ctx.GlobalString("key")))
				if err != nil {
					return err
				}
				fmt.Printf("<Journal uid:%s\n     content: %s>\n", j.UID, content)
			}

			return nil
		},
	}
}

func (app *App) Run() { app.cli.RunAndExitOnError() }

func main() {
	NewApp().Run()
}
