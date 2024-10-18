/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var cmd *exec.Cmd

// runserverCmd represents the runserver command
var runserverCmd = &cobra.Command{
	Use:   "runserver",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runServerWithHotReload()
	},
}

func init() {
	rootCmd.AddCommand(runserverCmd)
}

func runServerWithHotReload() {
    fmt.Println("Starting server with hot reload...")

    // Create a new file watcher
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    // Watch the current directory and its subdirectories
    if err := watchRecursive(".", watcher); err != nil {
        log.Fatal(err)
    }

    // Start the server for the first time
    startServer()

    // Signal channel to listen for termination (CTRL+C) and exit cleanly
    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

    // Channel to handle file events
    done := make(chan bool)

    go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                // Check if the file is modified, renamed, or created
                if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Rename == fsnotify.Rename || event.Op&fsnotify.Create == fsnotify.Create {
                    fmt.Println("Detected file change:", event.Name)

                    // If a new directory is created, start watching it recursively
                    fileInfo, err := os.Stat(event.Name)
                    if err == nil && fileInfo.IsDir() {
                        if err := watchRecursive(event.Name, watcher); err != nil {
                            log.Printf("Error watching new directory: %s\n", err)
                        }
                    }

                    restartServer() // Restart the server on file change
                }

            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                log.Println("Error:", err)

            case <-signalChan:
                fmt.Println("Received interrupt, shutting down...")
                stopServer()
                done <- true
            }
        }
    }()

    <-done
}

// Start the server by running `go run main.go`
func startServer() {
    fmt.Println("Running `go run main.go`...")

    cmd = exec.Command("go", "run", "main.go")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    err := cmd.Start()
    if err != nil {
        fmt.Printf("Error starting server: %s\n", err)
        return
    }

    fmt.Printf("Server running with PID %d\n", cmd.Process.Pid)
}

// Stop the running server by killing the process
func stopServer() {
    if cmd != nil && cmd.Process != nil {
        fmt.Println("Stopping server...")

        // Send SIGTERM for graceful shutdown
        err := cmd.Process.Signal(syscall.SIGTERM)
        if err != nil {
            fmt.Printf("Error sending SIGTERM: %s\n", err)
        }

        // Wait for the process to exit or kill it if it doesn't stop
        exitChan := make(chan error)
        go func() {
            exitChan <- cmd.Wait()
        }()

        select {
        case err = <-exitChan:
            if err != nil {
                fmt.Printf("Error waiting for server to stop: %s\n", err)
            } else {
                fmt.Println("Server stopped successfully.")
            }
        case <-time.After(5 * time.Second): // Timeout after 5 seconds
            fmt.Println("Graceful shutdown timed out, forcefully killing the server...")

            // If timeout occurs, force kill the process and its children
            killProcessAndChildren(cmd.Process.Pid)
        }

        cmd = nil

        // Add a small delay to ensure port is released
        time.Sleep(1 * time.Second)
    }
}

// Restart the server by stopping and then starting it again
func restartServer() {
    stopServer()
    // Short delay to ensure the file changes are written before restarting
    time.Sleep(500 * time.Millisecond)
    startServer()
}

// Recursively watch directories and subdirectories
func watchRecursive(dir string, watcher *fsnotify.Watcher) error {
    return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            err = watcher.Add(path)
            if err != nil {
                return fmt.Errorf("failed to watch directory %s: %v", path, err)
            }
            fmt.Printf("Watching directory: %s\n", path)
        }
        return nil
    })
}

// killProcessAndChildren forcefully kills the process and all of its children
func killProcessAndChildren(pid int) {
    // Kill the process tree starting from the PID
    fmt.Printf("Killing process and children with PID %d\n", pid)

    process := exec.Command("pkill", "-TERM", "-P", fmt.Sprint(pid))
    process.Stdout = os.Stdout
    process.Stderr = os.Stderr
    if err := process.Run(); err != nil {
        fmt.Printf("Error killing process tree: %s\n", err)
    }
}
