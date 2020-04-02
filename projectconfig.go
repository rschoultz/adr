package main

import (
	"encoding/json"
	"errors"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// AdrProjectConfig Project ADR configuration, loaded and used by each sub-command
type AdrProjectConfig struct {
	BaseDir    string `json:"base_directory"`
	CurrentAdr int    `json:"current_id"`
}

var adrProjectConfigFolderName = ".adr"
var adrProjectConfigFileName = "config.json"
var adrProjectConfigTemplateName = "template.md"

func projectConfigFolderPath(projectPath string) string {
	return filepath.Join(projectPath, adrProjectConfigFolderName)
}

func projectConfigFilePath(projectPath string) string {
	projectConfigPath := projectConfigFolderPath(projectPath)
	return filepath.Join(projectConfigPath, adrProjectConfigFileName)
}

func projectTemplatePath(projectPath string) string {
	projectConfigPath := projectConfigFolderPath(projectPath)
	adrTemplateFilePath := filepath.Join(projectConfigPath, adrProjectConfigTemplateName)
	return adrTemplateFilePath
}

func addProject(dir string, adrConfig AdrConfig) (AdrConfig, error) {
	projects := adrConfig.Projects

	for _, s := range projects {
		if s == dir {
			color.Green("Project is already listed")
			return adrConfig, nil
		}

		if strings.HasPrefix(dir, s) {
			color.Red("Cannot add sub-directory, project already added: " + s)
			return adrConfig, errors.New("cannot add-subdirectory")
		}

		if strings.HasPrefix(s, dir) {
			color.Red("Cannot add a parent directory to one that is already registered: " + s)
			return adrConfig, errors.New("cannot add a parent directory to one that is already registered")
		}

	}

	newProjects := append(projects, dir)
	return AdrConfig{newProjects}, nil
}

func initProjectBaseDir(baseDir string) {
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		err := os.MkdirAll(baseDir, 0755)
		if err != nil {
			color.Red("Could not create " + baseDir)
			os.Exit(1)
		}
	} else {
		color.Red(baseDir + " already exists, skipping folder creation")
	}
}

func initProjectConfig(projectPath string, baseDir string) {
	projectConfigFolderPath := projectConfigFolderPath(projectPath)
	if _, err := os.Stat(projectConfigFolderPath); os.IsNotExist(err) {
		err := os.Mkdir(projectConfigFolderPath, 0744)
		if err != nil {
			color.Red("Could not create directory " + projectConfigFolderPath)
			panic(err)
		}
	}
	config := AdrProjectConfig{baseDir, 0}
	bytes, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		panic(err)
	}

	configFilePath := projectConfigFilePath(projectPath)
	err = ioutil.WriteFile(configFilePath, bytes, 0644)
	if err != nil {
		color.Red("Could not write config file " + configFilePath)
		panic(err)
	}
}

func initTemplate(projectPath string) {
	body := []byte(`# {{.Number}}. {{.Title}}

Date: {{.Date}}

## Status

{{.Status}}

## Context

_The issue motivating this decision,and any context that influences 
or constrains the decision._

## Decision

_The change that we're proposing or have agreed to implement._

## Consequences

_What becomes easier or more difficult to do and any risks introduced
by the change that will need to be mitigated._
`)

	err := ioutil.WriteFile(projectTemplatePath(projectPath), body, 0644)
	if err != nil {
		color.Red("Could not write project template file")
		panic(err)
	}
}

func updateProjectConfig(projectPath string, config AdrProjectConfig) {
	bytes, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		panic(err)
	}

	configFilePath := projectConfigFilePath(projectPath)
	err = ioutil.WriteFile(configFilePath, bytes, 0644)
	if err != nil {
		color.Red("Could not update project configuration " + configFilePath)
		panic(err)
	}
}

func getProjectConfig(projectPath string) (currentProjectConfig AdrProjectConfig) {

	configFilePath := projectConfigFilePath(projectPath)
	bytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		color.Red("No ADR configuration is found for this project!")
		color.HiGreen("Start by initializing ADR configuration, check 'adr init --help' for more help")
		os.Exit(1)
	}

	err = json.Unmarshal(bytes, &currentProjectConfig)
	if err != nil {
		color.Red("Could not parse config file " + configFilePath)
		panic(err)
	}
	return currentProjectConfig
}

func getProjectPathBySubDir(currentDir string, homeConfig AdrConfig) (string, error) {
	for _, s := range homeConfig.Projects {
		if strings.HasPrefix(currentDir, s) {
			return s, nil
		}
	}
	return "", errors.New("current project not registered")
}
