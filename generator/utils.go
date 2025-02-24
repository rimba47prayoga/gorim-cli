package generator

import (
	"embed"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/*
var TemplatesFS embed.FS

// Template context struct
type TemplateContext struct {
	ProjectName string
	AppName     string
}

func CreateFileFromTemplate(destPath, tmplPath string, context TemplateContext) error {
	// Read template from embedded filesystem
	content, err := TemplatesFS.ReadFile(tmplPath)
	if err != nil {
		return err
	}

	tmpl, err := template.New(filepath.Base(tmplPath)).Parse(string(content))
	if err != nil {
		return err
	}

	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Pass context when executing the template
	return tmpl.Execute(file, context)
}