package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "aocf",
	Short: "Advent of Code Forever",
	Long:  "Advent of Code codebase for the long-term.",
}
