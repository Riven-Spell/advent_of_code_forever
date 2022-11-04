package solution_templates

import (
	"embed"
	"fmt"
	"path/filepath"
	"text/template"
)

//go:embed template_files
var templateFS embed.FS

func prepareTemplate(tName string) *template.Template {
	buf, err := templateFS.ReadFile(filepath.Join("template_files", tName))

	if err != nil {
		panic(fmt.Errorf("failed opening %s: %w", tName, err))
	}

	out, err := template.New(tName).Parse(string(buf))

	if err != nil {
		panic(fmt.Errorf("failed parsing %s: %w", tName, err))
	}

	return out
}

var SolutionTemplate = prepareTemplate("solution.go.template")

type SolutionTemplateInfill struct {
	Package string
	Day     uint
	Year    uint
}

var ImporterTemplate = prepareTemplate("importer.go.template")

type ImporterTemplateInfill struct {
	Imports []string
}
