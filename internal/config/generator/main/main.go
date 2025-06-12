// Command main is the entry point for the flag generator.
// It extracts flag constants from the config package and updates the README.md file.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rozdolsky33/ocloud/internal/config/generator"
)

func main() {
	// Get the project root directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	// Find the project root (where go.mod is located)
	projectRoot := findProjectRoot(wd)
	if projectRoot == "" {
		fmt.Fprintf(os.Stderr, "Error: could not find project root (go.mod file)\n")
		os.Exit(1)
	}

	// Extract flag constants
	configDir := filepath.Join(projectRoot, "internal", "config")
	flagInfos, err := generator.ExtractFlagConstants(configDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting flag constants: %v\n", err)
		os.Exit(1)
	}

	// Update README.md
	readmePath := filepath.Join(projectRoot, "README.md")
	err = generator.UpdateReadme(readmePath, flagInfos)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating README.md: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully updated README.md with flag documentation")
}

// findProjectRoot finds the project root directory by looking for a go.mod file.
func findProjectRoot(dir string) string {
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
