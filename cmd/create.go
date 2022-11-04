package cmd

import (
	"fmt"
	"github.com/Riven-Spell/advent_of_code_forever/solutions"
	"github.com/Riven-Spell/advent_of_code_forever/solutions/solution_templates"
	"github.com/Riven-Spell/advent_of_code_forever/util"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var createArgs = struct {
	Next    string
	Day     uint // 1 <= day <= 25
	Year    uint // 2015 <= day
	Replace bool
}{}

var create = &cobra.Command{
	Use:   "create {--next day/year | --day <day> --year <year>}",
	Short: "Create a new day or year.",
	Long:  "Attempts to create a new day/year if not present.",

	RunE: func(cmd *cobra.Command, args []string) error {
		codeDir, err := util.CodeDirName()
		if err != nil {
			fmt.Printf("failed to get code directory: %w\n", err)
			return nil
		}
		workDir, err := os.Getwd()
		if err != nil {
			if err != nil {
				fmt.Printf("failed to get pwd: %w\n", err)
				return nil
			}
		}

		if filepath.Dir(codeDir) != workDir {
			fmt.Printf("cannot generate code unless in the source directory\n")
			return nil
		}

		dayTemp := solution_templates.SolutionTemplateInfill{}
		cDay, cYear := solutions.Index.GetCurrentDay()
		mode := strings.ToLower(strings.TrimSpace(createArgs.Next))

		if createArgs.Day != 0 || createArgs.Year != 0 {
			mode = ""

			if createArgs.Day != 0 {
				cDay = createArgs.Day
			} else {
				mode = "day"
			}

			if createArgs.Year != 0 {
				cYear = createArgs.Year
			}
		}

		// select the correct day/year
		switch mode {
		case "day":
			cDay++
			if cDay <= 25 {
				break
			}
			fallthrough
		case "year":
			cYear++
			cDay = 1
		default:
			return fmt.Errorf("no such next mode '%s'", mode)
		}

		if cDay == 0 {
			cDay++
		}

		if !createArgs.Replace {
			// check we're not overwriting anything
			if solutions.Index.Get(cDay, cYear) != nil {
				fmt.Printf("Not overwriting day %d/%d as it already exists.", cYear, cDay)
				return nil
			}
		}

		dayTemp.Day = cDay
		dayTemp.Year = cYear
		dayTemp.Package = fmt.Sprintf("day%d", cDay)

		solutionsPackage := filepath.Join(workDir, "solutions/solution_code")
		dayPackage := filepath.Join(solutionsPackage, fmt.Sprint(cYear), dayTemp.Package)
		err = os.MkdirAll(dayPackage, 0755)
		if err != nil {
			fmt.Printf("cannot create folders: %s\n", err.Error())
			return nil
		}

		// Generate the template code
		{
			codeFileName := filepath.Join(dayPackage, dayTemp.Package+".go")

			f, err := os.OpenFile(codeFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				fmt.Printf("cannot open file: %s\n", err.Error())
				return nil
			}

			err = solution_templates.SolutionTemplate.Execute(f, dayTemp)
			if err != nil {
				fmt.Printf("cannot fill template: %s\n", err.Error())
				return nil
			}

			_ = f.Close()
		}

		// Generate the importer code
		err = solution_templates.UpdateImporter()
		if err != nil {
			fmt.Printf("failed to update importer: %s\n", err.Error())
			return nil
		}

		return nil
	},
}

func init() {
	create.PersistentFlags().BoolVar(&createArgs.Replace, "replace", false, "Replace an existing day? (default: false)")
	create.PersistentFlags().StringVar(&createArgs.Next, "next", "day", "Create the next day or year? Falls back to year if already on day 25 of last year.")
	create.PersistentFlags().UintVar(&createArgs.Day, "day", 0, "Specify a day to create (1-25)")
	create.PersistentFlags().UintVar(&createArgs.Year, "year", 0, "Specify a year to create (2015-onward). Current year assumed if not specified.")

	RootCmd.AddCommand(create)
}
