package main

import (
	"fmt"
	"os"
	"text/template"

	_ "embed"
)

//go:embed changelog.tmpl
var changelogTemplateRaw string

type ChangelogData struct {
	ReleaseName string
	FromTag     string
	ToTag       string
	Fixes       map[string]string
	Ci          map[string]string
	Docs        map[string]string
	Authors     map[string]string
	AddEmail    bool
	GhMarkdown  bool
}

func executeTemplate(data ChangelogData, stdout bool, outputPath string) error {
	parsedTemplate, err := template.New("changelog").Parse(changelogTemplateRaw)
	if err != nil {
		return fmt.Errorf("error parsing the template: %w", err)
	}

	if stdout {
		if err := parsedTemplate.Execute(os.Stdout, data); err != nil {
			return fmt.Errorf("error executing template onto stdout: %w", err)
		}
	}

	if outputPath != "" {
		fd, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("error opening the output file: %w")
		}
		defer fd.Close()

		if err := parsedTemplate.Execute(fd, data); err != nil {
			return fmt.Errorf("error executing template onto %q: %w", outputPath, err)
		}
	}

	return nil
}
