package solution_templates

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// UpdateImporter regenerates the importer.
// It should never be run if the working directory is not the root of advent_of_code_forever.
func UpdateImporter() error {
	workDir, err := os.Getwd() // wd is assumed to be code dir.
	if err != nil {
		if err != nil {
			return fmt.Errorf("failed to get pwd: %w\n", err)
		}
	}

	solutionsPackage := filepath.Join(workDir, "solutions/solution_code")

	importerName := filepath.Join(workDir, "solutions/solution_code/importer.go")
	f, err := os.OpenFile(importerName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("cannot open file %s: %w\n", importerName, err)
	}

	includeSet := map[string]bool{}

	err = filepath.WalkDir(
		solutionsPackage,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && strings.HasSuffix(d.Name(), ".go") {
				if dir := filepath.Dir(path); dir != solutionsPackage {
					pkgStartIndex := strings.Index(dir, "solution_code/")
					dir = strings.TrimPrefix(dir[pkgStartIndex:], "solution_code/")
					includeSet[dir] = true
				}
			}

			return nil
		})

	includes := make([]string, 0)

	for k := range includeSet {
		includes = append(includes, k)
	}

	err = ImporterTemplate.Execute(f, ImporterTemplateInfill{Imports: includes})
	if err != nil {
		return err
	}

	return f.Close()
}
