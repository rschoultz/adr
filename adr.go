package main

import (
	"github.com/fatih/color"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

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

const (
	layoutISO = "2006-01-02"
)

func newAdr(projectPath string, config AdrProjectConfig, adrName []string) (adrFullPath string) {

	adr := Adr{
		Title:  strings.Join(adrName, " "),
		Date:   time.Now().Format(layoutISO),
		Number: config.CurrentAdr,
		Status: PROPOSED,
	}

	useTemplate, err := template.ParseFiles(projectTemplatePath(projectPath))
	if err != nil {
		panic(err)
	}
	adrFileName := strconv.Itoa(adr.Number) + "-" + strings.Join(strings.Split(strings.Trim(adr.Title, "\n \t"), " "), "-") + ".md"
	adrFullPath = filepath.Join(projectPath, config.BaseDir, adrFileName)
	f, err := os.Create(adrFullPath)
	if err != nil {
		panic(err)
	}
	err = useTemplate.Execute(f, adr)
	if err != nil {
		color.Red("Problems writing ADR on " + adrFullPath)
		panic(err)
	}
	_ = f.Close()
	color.Green("ADR number " + strconv.Itoa(adr.Number) + " was successfully written to : " + adrFullPath)

	return
}
