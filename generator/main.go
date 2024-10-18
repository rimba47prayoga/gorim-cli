package generator

import "fmt"

// getMainGoContent returns the content for the main.go file
func getMainGoContent(projectName string) string {
    return fmt.Sprintf(`package main

import (
    "fmt"
    "net/http"
    "%s/settings"
)

func main() {
    fmt.Printf("Starting %s on port %s\n", settings.Config.AppName, settings.Config.Port)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, world!")
    })

    if err := http.ListenAndServe(":" + settings.Config.Port, nil); err != nil {
        fmt.Printf("Failed to start server: %v\n", err)
    }
}
`, projectName, projectName)
}