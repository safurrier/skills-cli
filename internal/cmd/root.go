// Package cmd implements the cobra CLI commands for skills-cli.
package cmd

import (
	"github.com/spf13/cobra"
)

var dir string

var rootCmd = &cobra.Command{
	Use:   "skills",
	Short: "A repo-aware CLI for Agent Skills",
	Long:  "Discover, inspect, and invoke Agent Skills (SKILL.md) in your repository.",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dir, "dir", ".", "Root directory to scan for skills")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
