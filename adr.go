package main

import (
	"bufio"
	"github.com/fatih/color"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Adr basic structure
type Adr struct {
	Number   int
	Title    string
	Date     string
	Status   AdrStatus
	Filename string
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
	adrFileName := strconv.Itoa(adr.Number) + "-" + strings.Join(strings.Split(strings.Trim(strings.ToLower(adr.Title), "\n \t"), " "), "-") + ".md"
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

func visitAdrFile(adrs *[]Adr) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		var adr Adr
		firstLine := ""
		status := ""

		if err != nil {
			log.Fatal(err)
		}
		if info.IsDir() {
			return nil
		}

		file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
		defer file.Close()

		sc := bufio.NewScanner(file)
		for sc.Scan() {
			line := sc.Text()
			matchTitle, _ := regexp.MatchString("# \\d+\\.", line)
			if firstLine == "" && matchTitle {
				firstLine = line
				adr.Title = strings.Replace(firstLine, "# ", "", 1)
			}

			matchStatus, _ := regexp.MatchString(string("^"+PROPOSED+"|^"+ACCEPTED+"|^"+DEPRECATED+"|^"+SUPERSEDED), line)
			if status == "" && matchStatus {
				status = line
				adr.Status = AdrStatus(status)
			}

			if strings.HasPrefix(line, "Date: ") {
				date := strings.Replace(line, "Date: ", "", 1)
				adr.Date = date
			}
		}

		adr.Filename = path
		*adrs = append(*adrs, adr)

		return nil
	}
}

func getAdrs(projectPath string, currentConfig AdrProjectConfig) []Adr {
	var adrs []Adr
	adrDir := filepath.Join(projectPath, currentConfig.BaseDir)
	err := filepath.Walk(adrDir, visitAdrFile(&adrs))
	if err != nil {
		color.Red("Could not open directory of ADRs: " + adrDir)
		panic(err)
	}
	return adrs
}
