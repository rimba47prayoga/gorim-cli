package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var startProjectFiles = map[string]string{
	".gitignore":         "templates/.gitignore.tmpl",
	".env":               "templates/.env.tmpl",
	"LICENSE":            "templates/LICENSE.tmpl",
	"api/routes.go":      "templates/api/routes.go.tmpl",
	"main.go":            "templates/main.go.tmpl",
	"migrations/init.go": "templates/migrations/init.go.tmpl",
	"migrations/register.go": "templates/migrations/register.go.tmpl",
	"settings/config.go": "templates/settings/config.go.tmpl",
}

// initGoMod initializes a Go module with the given project name.
func initGoMod(projectName string) error {
	fmt.Println("Initializing Go module...")

	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize go module: %w", err)
	}

	fmt.Println("Go module initialized successfully!")

	// Install Gorim framework
	if err := installGorim(); err != nil {
		return err
	}

	if err := installSqlite(); err != nil {
		return err
	}

	return nil
}

func installSqlite() error {
	fmt.Println("Installing sqlite as default database..")

	cmd := exec.Command("go", "get", "github.com/glebarez/sqlite")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install sqlite: %w", err)
	}

	fmt.Println("Sqlite installed successfully!")
	return nil
}

// installGorim runs `go get gorim.org/gorim`
func installGorim() error {
	fmt.Println("Installing gorim.org/gorim")

	cmd := exec.Command("go", "get", "gorim.org/gorim")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Gorim: %w", err)
	}

	fmt.Println("Gorim installed successfully!")
	return nil
}

func StartProject(projectName string) {
	if projectName == "." {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			return
		}
		projectName = filepath.Base(dir)
	} else {
		if err := os.MkdirAll(projectName, os.ModePerm); err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}

		if err := os.Chdir(projectName); err != nil {
			fmt.Println("Error changing directory:", err)
			return
		}
	}

	// Initialize Go module
	if err := initGoMod(projectName); err != nil {
		fmt.Println(err)
		return
	}

	context := TemplateContext{ProjectName: projectName}

	for path, tmplPath := range startProjectFiles {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Println("Error creating directory:", err)
			continue
		}

		if err := CreateFileFromTemplate(path, tmplPath, context); err != nil {
			fmt.Println("Error creating file:", err)
		}
	}

	fmt.Println("Project generated successfully!")
}