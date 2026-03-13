// Package main is the entry point for skills-cli.
package main

import (
	"os"

	"github.com/safurrier/skills-cli/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
