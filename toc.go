package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func generateToc(projectPath string, currentConfig AdrProjectConfig, prefix string) {
	adrs := getAdrs(projectPath, currentConfig)

	fmt.Printf("# Architecture Decision Records\n\n")

	for _, adr := range adrs {
		filename := strings.Replace(adr.Filename, projectPath+"/", "", 1)
		prefixedFilname := filepath.Join(prefix, filename)
		fmt.Printf("* [%s](%s) (%s)\n", adr.Title, prefixedFilname, adr.Status)
	}
}
