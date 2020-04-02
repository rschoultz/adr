package main

import (
	"encoding/json"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

// AdrConfig Home directory configuration lists directories where ADR Projects has been configured
type AdrConfig struct {
	Projects []string `json:"projects"`
}

var usr, _ = user.Current()
var adrHomeConfigFolderName = ".adr"
var adrHomeConfigFolderPath = filepath.Join(usr.HomeDir, adrHomeConfigFolderName)
var adrHomeConfigFilePath = filepath.Join(adrHomeConfigFolderPath, "projects.json")

func initHomeDir() {
	if _, err := os.Stat(adrHomeConfigFolderPath); os.IsNotExist(err) {
		err := os.Mkdir(adrHomeConfigFolderPath, 0744)
		if err != nil {
			color.Red("Could not create directory " + adrHomeConfigFolderPath)
		}
	}
}

func initHomeConfig() (newHomeConfig AdrConfig) {
	if _, err := os.Stat(adrHomeConfigFilePath); os.IsNotExist(err) {
		color.Green("Creating new home directory configuration")
		newHomeConfig = AdrConfig{nil}
		writeHomeConfiguration(newHomeConfig)
	}
	return
}

func getHomeConfig() AdrConfig {
	var currentHomeConfig AdrConfig

	bytes, err := ioutil.ReadFile(adrHomeConfigFilePath)
	if err != nil {
		color.Red("No ADR home directory configuration is found!")
		color.HiGreen("Start by initializing ADR configuration, check 'adr init --help' for more help")
		os.Exit(1)
	}

	err = json.Unmarshal(bytes, &currentHomeConfig)
	if err != nil {
		panic(err)
	}
	return currentHomeConfig
}

func writeHomeConfiguration(homeConfig AdrConfig) {
	bytes, err := json.MarshalIndent(homeConfig, "", " ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(adrHomeConfigFilePath, bytes, 0644)
	if err != nil {
		panic(err)
	}
}
