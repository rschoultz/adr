package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

func setCommands(app *cli.App) {
	var prefix string
	app.Commands = []*cli.Command{
		{
			Name:    "new",
			Aliases: []string{"c"},
			Usage:   "Create a new ADR",
			Action: func(c *cli.Context) error {

				currentDir, _ := os.Getwd()
				homeConfig := getHomeConfig()
				projectPath, _ := getProjectPathBySubDir(currentDir, homeConfig)
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
				homeConfig := initHomeConfig()
				newConfig, err := addProject(projectDir, homeConfig)
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
		{
			Name:        "generate",
			Aliases:     []string{"g"},
			Usage:       "Generate content",
			Description: "Generate additional content, like table of contents, based on ADRs already written",
			Subcommands: []*cli.Command{
				{
					Name:  "toc",
					Usage: "generate table of contents",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "prefix",
							Aliases:     []string{"p"},
							Usage:       "prefix each file link",
							Required:    false,
							Destination: &prefix,
						},
					},
					Action: func(c *cli.Context) error {
						currentDir, _ := os.Getwd()
						homeConfig := getHomeConfig()
						projectPath, _ := getProjectPathBySubDir(currentDir, homeConfig)
						currentConfig := getProjectConfig(projectPath)
						generateToc(projectPath, currentConfig, prefix)
						return nil
					},
				},
			},
		},
	}
}
