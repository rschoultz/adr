package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

func setCommands(app *cli.App) {
	app.Commands = []*cli.Command{
		{
			Name:    "new",
			Aliases: []string{"c"},
			Usage:   "Create a new ADR",
			Flags:   []cli.Flag{},
			Action: func(c *cli.Context) error {

				projectPath, _ := getProjectPathByCwd()
				currentConfig := getProjectConfig(projectPath)
				currentConfig.CurrentAdr++
				updateProjectConfig(projectPath, currentConfig)
				filename := newAdr(projectPath, currentConfig, c.Args().Slice())
				startEditor(filename)
				return nil
			},
		},

		{
			Name:        "init",
			Aliases:     []string{"i"},
			Usage:       "Initializes the ADR configurations",
			UsageText:   "adr init docs/architecture/decisions",
			Description: "Initializes the ADR configuration with an optional ADR base directory\n This is a a prerequisite to running any other adr sub-command",
			Action: func(c *cli.Context) error {
				projectDir, _ := os.Getwd()
				initDir := c.Args().First()
				if initDir == "" {
					initDir = adrProjectConfigFolderName
				}
				initHomeDir()
				initHomeConfigIfNotExists()
				homeConfig := getHomeConfig()
				newConfig, err := addProject(homeConfig, projectDir)
				if err != nil {
					os.Exit(1)
				}
				initProjectBaseDir(initDir)
				writeHomeConfiguration(newConfig)
				initProjectConfig(projectDir, initDir)
				initTemplate(projectDir)
				return nil
			},
		},
	}
}
