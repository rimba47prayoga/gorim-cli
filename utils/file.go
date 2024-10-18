package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// createProjectStructure generates the necessary files and directories
func createProjectStructure(projectName string) {
    // Define project directories
    projectDir := filepath.Join(".", projectName)
    dirs := []string{
        filepath.Join(projectDir, "views"),
        filepath.Join(projectDir, "models"),
        filepath.Join(projectDir, "settings"),
        filepath.Join(projectDir, "router"),
    }

    // Create directories
    for _, dir := range dirs {
        if err := os.MkdirAll(dir, os.ModePerm); err != nil {
            fmt.Printf("Failed to create directory %s: %v\n", dir, err)
            return
        }
    }

    // Generate settings.go
    // generateFile(filepath.Join(projectDir, "settings", "settings.go"), getSettingsContent(projectName))

    // // Generate main.go
    // generateFile(filepath.Join(projectDir, "main.go"), getMainGoContent(projectName))

    fmt.Printf("Project %s created successfully!\n", projectName)
}

// generateFile creates a file with the given content
func GenerateFile(path, content string) {
    f, err := os.Create(path)
    if err != nil {
        fmt.Printf("Failed to create file %s: %v\n", path, err)
        return
    }
    defer f.Close()

    if _, err := f.WriteString(content); err != nil {
        fmt.Printf("Failed to write to file %s: %v\n", path, err)
        return
    }
}

