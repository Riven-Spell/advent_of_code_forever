package cmd

import (
	"fmt"
	"github.com/Riven-Spell/advent_of_code_forever/inputs"
	"github.com/Riven-Spell/advent_of_code_forever/solutions"
	"github.com/Riven-Spell/advent_of_code_forever/util"
	"github.com/spf13/cobra"
	"reflect"
	"strings"
	"time"
)

var runArgs = struct {
	Year, Day       uint
	Part            int
	All             bool
	CacheAnswers    bool
	InputMode       string // cache, download, generate
	InputComplexity uint64
}{}

func runDay(cDay, cYear uint) error {
	day := solutions.Index.Get(cDay, cYear)
	if day == nil {
		return fmt.Errorf("day %d/%d is not available")
	}

	var input string
	var solution *inputs.Solution
	var err error

	inputMode := strings.ToLower(runArgs.InputMode)
	switch inputMode {
	case "download":
		err = inputs.Cache.DownloadInput(cDay, cYear, true)
		if err != nil {
			fmt.Printf("Day %d/%d: Failed to download input: %s\n", cYear, cDay, err.Error())
			return nil
		}
		fallthrough
	case "cache":
		input, solution, err = inputs.Cache.GetInputAndSolution(cDay, cYear)
		if err != nil {
			fmt.Printf("Day %d/%d: Failed to pull input from cache: %s\n", cYear, cDay, err.Error())
			return nil
		}
	case "generate":
		if day.Generator == nil {
			return fmt.Errorf("day %d/%d does not contain an input generator", cYear, cDay)
		}

		complexity := runArgs.InputComplexity
		if complexity == 0 {
			complexity = day.DefaultComplexity
		}

		input, solution = day.Generator(complexity)
	}

	runner := day.Solution
	if solution == nil {
		solution = &inputs.Solution{}
	}

	getResult := func(result, expected any) string {
		if expected == nil {
			return ""
		}

		if reflect.DeepEqual(result, expected) {
			return " (PASSED)"
		} else {
			return " (FAILED: expected " + fmt.Sprint(expected) + ")"
		}
	}

	if runArgs.Part == -1 || runArgs.Part == 1 {
		// Part 1
		runner.Prepare(input)
		startTime := time.Now() // time the run
		result := runner.Part1()
		if result != nil {
			endTime := time.Now().Sub(startTime)
			fmt.Printf("PART 1: %v%s in %s\n", result, getResult(result, solution.A), endTime.String())
		}
	}

	if runArgs.Part == -1 || runArgs.Part == 2 {
		// Part 1
		runner.Prepare(input)
		startTime := time.Now() // time the run
		result := runner.Part2()
		if result != nil {
			endTime := time.Now().Sub(startTime)
			fmt.Printf("PART 2: %v%s in %s\n", result, getResult(result, solution.B), endTime.String())
		}
	}

	return nil
}

var runCommand = &cobra.Command{
	Use:   "run [--year <year> --day <day> | --all] [--part <1/2>]",
	Short: "Runs a day with it's input. Can generate or download input on the fly. If no year/day is specified, both parts of the most recent day will be ran if available.",

	RunE: func(cmd *cobra.Command, args []string) error {
		if runArgs.All {
			_, maxYear := solutions.Index.GetCurrentDay()
			if runArgs.Year != 0 {
				maxYear = runArgs.Year
			}

			for y := util.Ternary(runArgs.Year != 0, runArgs.Year, 2015); y <= maxYear; y++ {
				maxDay := solutions.Index.GetCurrentDayForYear(y)
				for d := uint(1); d <= maxDay; d++ {
					if err := runDay(d, y); err != nil {
						return err
					}
				}
			}
		} else {
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

			if err := runDay(cDay, cYear); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	runCommand.PersistentFlags().UintVar(&runArgs.Year, "year", 0, "Specified year of solutions to run. If specified with --all, runs every day of that year.")
	runCommand.PersistentFlags().UintVar(&runArgs.Day, "day", 0, "Specified day of solutions to run.")
	runCommand.PersistentFlags().IntVar(&runArgs.Part, "part", -1, "0 or 1. Runs both by default.")
	runCommand.PersistentFlags().BoolVar(&runArgs.All, "all", false, "Run all days available (of all years if year is unspecified).")
	runCommand.PersistentFlags().BoolVar(&runArgs.CacheAnswers, "cache-answers", false, "Overwrite existing input & solution with new data from these runs.")
	runCommand.PersistentFlags().StringVar(&runArgs.InputMode, "input-mode", "cache", "How to obtain an input for the run. (cache/download/generate) (defaults to cache).")
	runCommand.PersistentFlags().Uint64Var(&runArgs.InputComplexity, "input-complexity", 0, "Input complexity to generate at. Defaults to that specified by the problem's code.")

	RootCmd.AddCommand(runCommand)
}
