package cmd

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/Riven-Spell/advent_of_code_forever/inputs"
	"github.com/Riven-Spell/advent_of_code_forever/solutions"
	"github.com/Riven-Spell/advent_of_code_forever/util"
	"github.com/spf13/cobra"
	"strings"
)

var cacheArgs = struct {
	Day             uint // 1 <= day <= 25
	Year            uint // 2015 <= day
	Replace         bool
	InputComplexity uint64
	Mode            string
}{}

//go:embed help/cache_help.txt
var longCacheHelp string

var cache = &cobra.Command{
	Use:   "cache",
	Short: "Download, Generate, Replace, or Delete inputs",
	Long:  longCacheHelp,

	RunE: func(cmd *cobra.Command, args []string) error {
		cDay, cYear := solutions.Index.GetCurrentDay()

		if cacheArgs.Day != 0 || cacheArgs.Year != 0 {
			if cacheArgs.Day != 0 {
				cDay = cacheArgs.Day
			}

			if cacheArgs.Year != 0 {
				cYear = cacheArgs.Year
			}
		}

		if cDay == 0 {
			cDay++
		}

		var err error
		cache := inputs.Cache
		switch strings.ToLower(cacheArgs.Mode) {
		case "delete":
			err = cache.DeleteInput(cDay, cYear)
		case "generate":
			day := solutions.Index.Get(cDay, cYear)
			if day != nil || day.Generator == nil {
				err = fmt.Errorf("could not generate input: no generator present for day %d/%d", cYear, cDay)
			}

			complexity := cacheArgs.InputComplexity
			if complexity == 0 {
				complexity = day.DefaultComplexity
			}

			inputData, solutions := day.Generator(complexity) // todo: complexity
			err = cache.PutInput(cDay, cYear, strings.NewReader(inputData), cacheArgs.Replace)
			if err == nil && solutions != nil {
				err = cache.PutSolution(cDay, cYear, *solutions, cacheArgs.Replace)
			}
		case "download":
			err = cache.DownloadInput(cDay, cYear, cacheArgs.Replace)
		default:
			return errors.New("unknown cache operation: " + strings.ToLower(cacheArgs.Mode))
		}

		if err != nil {
			fmt.Printf("Failed to %s input: %s", strings.ToLower(cacheArgs.Mode), err.Error())
		} else {
			fmt.Printf("Successfully %s input for %d/%d",
				cacheArgs.Mode+util.Ternary(strings.HasSuffix(cacheArgs.Mode, "e"), "d", "ed"),
				cYear, cDay,
			)
		}

		return nil
	},
}

func init() {
	cache.PersistentFlags().UintVar(&cacheArgs.Day, "day", 0, "Day to target. Current day assumed if not specified.")
	cache.PersistentFlags().UintVar(&cacheArgs.Year, "year", 0, "Year to target. Current year assumed if not specified.")
	cache.PersistentFlags().BoolVar(&cacheArgs.Replace, "replace", false, "Replace the existing input if it was already obtained/generated? (default: false)")
	cache.PersistentFlags().StringVar(&cacheArgs.Mode, "mode", "generate", "Download/Generate/Delete (default: generate)")
	cache.PersistentFlags().Uint64Var(&cacheArgs.InputComplexity, "complexity", 0, "Input complexity to generate with (default: 100)")

	RootCmd.AddCommand(cache)
}
