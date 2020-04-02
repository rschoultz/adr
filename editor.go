package main

import (
	"github.com/fatih/color"
	"log"
	"os"
	"os/exec"
)

func startEditor(filename string) {
	editor, foundEditor := os.LookupEnv("VISUAL")

	if !foundEditor || editor == "" {
		editor, foundEditor = os.LookupEnv("EDITOR")
	}

	if !foundEditor || editor == "" {
		color.Green("No VISUAL or EDITOR environment variable set. Edit manually by (e.g.):\nvi " + filename)
		return
	}

	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		log.Fatalf("Cannot open file with editor: %s", err)
	}
}
