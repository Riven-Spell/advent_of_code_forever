package cmd

import (
	"fmt"
	"github.com/Riven-Spell/advent_of_code_forever/core"
	"github.com/Riven-Spell/advent_of_code_forever/util"
	"github.com/spf13/cobra"
)

var envCommand = &cobra.Command{
	Use:   "env",
	Short: "Print all environment variables relevant to aocf",

	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range core.EnvironmentVariables {
			fmt.Println(v.Name)
			val, defaulted := v.Get()
			fmt.Printf("Value: %s", util.Ternary(v.Secret, "[REDACTED]", val))
			if defaulted {
				fmt.Print(" (default)")
			}
			fmt.Print("\n")
			if v.Default != "" {
				fmt.Printf("Default: %s\n", v.Default)
			}
			fmt.Println()
		}
	},
}

func init() {
	RootCmd.AddCommand(envCommand)
}
