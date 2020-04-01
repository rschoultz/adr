package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// AdrConfig Home directory configuration lists directories where ADR Projects has been configured
type AdrConfig struct {
	Projects []string `json:"projects"`
}

// AdrProjectConfig Project ADR configuration, loaded and used by each sub-command
type AdrProjectConfig struct {
	BaseDir    string `json:"base_directory"`
	CurrentAdr int    `json:"current_id"`
}

// Adr basic structure
type Adr struct {
	Number int
	Title  string
	Date   string
	Status AdrStatus
}

// AdrStatus type
type AdrStatus string

// ADR status enums
const (
	PROPOSED   AdrStatus = "Proposed"
	ACCEPTED   AdrStatus = "Accepted"
	DEPRECATED AdrStatus = "Deprecated"
	SUPERSEDED AdrStatus = "Superseded"
)

var usr, _ = user.Current()
var adrConfigFolderName = ".adr"
var adrConfigFolderPath = filepath.Join(usr.HomeDir, adrConfigFolderName)
var adrConfigFilePath = filepath.Join(adrConfigFolderPath, "projects.json")

var adrProjectConfigFolderName = ".adr"
var adrProjectConfigFileName = "config.json"
var adrProjectConfigTemplateName = "template.md"

func initHomeDir() {
	if _, err := os.Stat(adrConfigFolderPath); os.IsNotExist(err) {
		err := os.Mkdir(adrConfigFolderPath, 0744)
		if err != nil {
			color.Red("Could not create directory " + adrConfigFolderPath)
		}
	}
}

func initHomeConfigIfNotExists() {
	if _, err := os.Stat(adrConfigFilePath); os.IsNotExist(err) {
		color.Green("Creating new home directory configuration")
		newHomeConfig := AdrConfig{nil}
		writeHomeConfiguration(newHomeConfig)
	}
}

func getHomeConfig() AdrConfig {
	var currentHomeConfig AdrConfig

	bytes, err := ioutil.ReadFile(adrConfigFilePath)
	if err != nil {
		color.Red("No ADR home directory configuration is found!")
		color.HiGreen("Start by initializing ADR configuration, check 'adr init --help' for more help")
		os.Exit(1)
	}

	json.Unmarshal(bytes, &currentHomeConfig)
	return currentHomeConfig
}

func addProject(adrConfig AdrConfig, dir string) (AdrConfig, error) {
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

func writeHomeConfiguration(homeConfig AdrConfig) {
	bytes, err := json.MarshalIndent(homeConfig, "", " ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(adrConfigFilePath, bytes, 0644)
	if err != nil {
		panic(err)
	}
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
	adrProjectConfigFolderPath := filepath.Join(projectPath, adrProjectConfigFolderName)
	if _, err := os.Stat(adrProjectConfigFolderPath); os.IsNotExist(err) {
		os.Mkdir(adrProjectConfigFolderPath, 0744)
	}
	config := AdrProjectConfig{baseDir, 0}
	bytes, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		panic(err)
	}

	var adrProjectConfigFilePath = filepath.Join(adrProjectConfigFolderPath, adrProjectConfigFileName)
	ioutil.WriteFile(adrProjectConfigFilePath, bytes, 0644)
}

func initTemplate(projectPath string) {
	body := []byte(`
# {{.Number}}. {{.Title}}
======
Date: {{.Date}}

## Status
======
{{.Status}}

## Context
======

## Decision
======

## Consequences
======

`)

	adrFullProjectConfigFolderPath := filepath.Join(projectPath, adrProjectConfigFolderName)
	adrTemplateFilePath := filepath.Join(adrFullProjectConfigFolderPath, adrProjectConfigTemplateName)

	ioutil.WriteFile(adrTemplateFilePath, body, 0644)
}

func updateProjectConfig(projectPath string, config AdrProjectConfig) {
	bytes, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		panic(err)
	}

	adrFullProjectConfigFolderPath := filepath.Join(projectPath, adrProjectConfigFolderName)
	projectConfigFilePath := filepath.Join(adrFullProjectConfigFolderPath, adrProjectConfigFileName)

	ioutil.WriteFile(projectConfigFilePath, bytes, 0644)
}

func getProjectConfig(projectPath string) AdrProjectConfig {
	var currentProjectConfig AdrProjectConfig

	adrFullProjectConfigFolderPath := filepath.Join(projectPath, adrProjectConfigFolderName)
	projectConfigFilePath := filepath.Join(adrFullProjectConfigFolderPath, adrProjectConfigFileName)

	bytes, err := ioutil.ReadFile(projectConfigFilePath)
	if err != nil {
		color.Red("No ADR configuration is found for this project!")
		color.HiGreen("Start by initializing ADR configuration, check 'adr init --help' for more help")
		os.Exit(1)
	}

	json.Unmarshal(bytes, &currentProjectConfig)
	return currentProjectConfig
}

func getProjectPathByCwd() (string, error) {
	getwd, _ := os.Getwd()
	projectPath, err := getProjectPathBySubDir(getwd)
	return projectPath, err
}

func getProjectPathBySubDir(currentDir string) (string, error) {
	homeConfig := getHomeConfig()

	for _, s := range homeConfig.Projects {
		if strings.HasPrefix(currentDir, s) {
			return s, nil
		}
	}
	return "", errors.New("Current project not registered. Use 'init' subcommand.")
}

func newAdr(projectPath string, config AdrProjectConfig, adrName []string) {

	adrFullProjectConfigFolderPath := filepath.Join(projectPath, adrProjectConfigFolderName)
	adrTemplateFilePath := filepath.Join(adrFullProjectConfigFolderPath, adrProjectConfigTemplateName)

	adr := Adr{
		Title:  strings.Join(adrName, " "),
		Date:   time.Now().Format(time.RFC3339),
		Number: config.CurrentAdr,
		Status: PROPOSED,
	}

	useTemplate, err := template.ParseFiles(adrTemplateFilePath)
	if err != nil {
		panic(err)
	}
	adrFileName := strconv.Itoa(adr.Number) + "-" + strings.Join(strings.Split(strings.Trim(adr.Title, "\n \t"), " "), "-") + ".md"
	adrFullPath := filepath.Join(projectPath, config.BaseDir, adrFileName)
	f, err := os.Create(adrFullPath)
	if err != nil {
		panic(err)
	}
	useTemplate.Execute(f, adr)
	f.Close()
	color.Green("ADR number " + strconv.Itoa(adr.Number) + " was successfully written to : " + adrFullPath)
}
