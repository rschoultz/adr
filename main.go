package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := cli.NewApp()
	app.Name = "adr"
	app.Usage = "Work with Architecture Decision Records (ADRs)"
	app.Version = "0.3.2"

	setCommands(app)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
