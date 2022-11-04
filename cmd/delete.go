package cmd

import (
	"fmt"
	"github.com/Riven-Spell/advent_of_code_forever/inputs"
	"github.com/Riven-Spell/advent_of_code_forever/solutions/solution_templates"
	"github.com/Riven-Spell/advent_of_code_forever/util"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var deleteArgs = struct {
	Year uint
	Day  uint
}{}

var deleteCMD = &cobra.Command{
	Use:   "delete [--day <day>] --year <year>",
	Short: "Delete an existing day's code. WARNING: this action is irrecoverable.",

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

		if deleteArgs.Year < 2015 {
			return fmt.Errorf("must specify a year on or after 2015")
		}

		if deleteArgs.Day > 25 {
			return fmt.Errorf("must specify a day before Dec. 25th")
		}

		solutionsPackage := filepath.Join(workDir, "solutions/solution_code")
		rmTarget := filepath.Join(solutionsPackage, fmt.Sprint(deleteArgs.Year))

		if deleteArgs.Day != 0 {
			rmTarget = filepath.Join(rmTarget, fmt.Sprintf("day%d", deleteArgs.Day))
		}

		err = os.RemoveAll(rmTarget)
		if err != nil {
			fmt.Printf("Failed to remove directory %s: %w", rmTarget, err)
			return nil
		}

		// Generate the importer code
		err = solution_templates.UpdateImporter()
		if err != nil {
			fmt.Printf("failed to update importer: %s\n", err.Error())
			return nil
		}

		err = inputs.Cache.DeleteInput(deleteArgs.Day, deleteArgs.Year)
		if err != nil {
			fmt.Printf("failed to delete inputs: %s", err.Error())
			return nil
		}

		return nil
	},
}

func init() {
	deleteCMD.PersistentFlags().UintVar(&deleteArgs.Year, "year", 0, "Year to remove (from). If used alone, removes the entire year. Must be set.")
	deleteCMD.PersistentFlags().UintVar(&deleteArgs.Day, "day", 0, "Day to remove from the target year. Optional.")

	RootCmd.AddCommand(deleteCMD)
}
