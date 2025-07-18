// Package generator provides tools to generate flag-related code and documentation.
package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"path/filepath"
	"regexp"
	"strings"
)

// FlagInfo represents information about a flag.
type FlagInfo struct {
	Name        string
	Shorthand   string
	Description string
	Default     string
}

// GenerateReadmeTable generates a markdown table of flags for the README.
func GenerateReadmeTable(flagInfos []FlagInfo, title string) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("### %s\n\n", title))
	buf.WriteString("| Flag | Short | Description |\n")
	buf.WriteString("|------|-------|-------------|\n")

	for _, flag := range flagInfos {
		shorthand := flag.Shorthand
		if shorthand != "" {
			shorthand = fmt.Sprintf("`-%s`", shorthand)
		}
		buf.WriteString(fmt.Sprintf("| `--%s` | %s | %s |\n", flag.Name, shorthand, flag.Description))
	}

	return buf.String()
}

// ExtractFlagConstants extracts flag constants from the configuration package.
func ExtractFlagConstants(configDir string) (map[string][]FlagInfo, error) {
	// Parse the flags/constants.go file
	fset := token.NewFileSet()
	flagsFile := filepath.Join(configDir, "flags", "constants.go")
	node, err := parser.ParseFile(fset, flagsFile, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags/constants.go: %w", err)
	}

	// Extract flag names, shorthands, and descriptions
	nameMap := make(map[string]string)
	shorthandMap := make(map[string]string)
	descMap := make(map[string]string)

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, name := range valueSpec.Names {
				if i >= len(valueSpec.Values) {
					continue
				}

				basicLit, ok := valueSpec.Values[i].(*ast.BasicLit)
				if !ok || basicLit.Kind != token.STRING {
					continue
				}

				value := strings.Trim(basicLit.Value, "\"")

				if strings.HasPrefix(name.Name, "FlagName") {
					nameMap[strings.TrimPrefix(name.Name, "FlagName")] = value
				} else if strings.HasPrefix(name.Name, "FlagShort") {
					shorthandMap[strings.TrimPrefix(name.Name, "FlagShort")] = value
				} else if strings.HasPrefix(name.Name, "FlagDesc") {
					descMap[strings.TrimPrefix(name.Name, "FlagDesc")] = value
				}
			}
		}
	}

	// Group flags by category
	globalFlags := []FlagInfo{}
	instanceFlags := []FlagInfo{}

	// Add global flags
	for key, name := range nameMap {
		if key == "LogLevel" || key == "TenancyID" || key == "TenancyName" || key == "Compartment" {
			shorthand := shorthandMap[key]
			desc := descMap[key]
			globalFlags = append(globalFlags, FlagInfo{
				Name:        name,
				Shorthand:   shorthand,
				Description: desc,
			})
		}
	}

	// Add instance flags
	for key, name := range nameMap {
		if key == "List" || key == "Find" || key == "ImageDetails" {
			shorthand := shorthandMap[key]
			desc := descMap[key]
			instanceFlags = append(instanceFlags, FlagInfo{
				Name:        name,
				Shorthand:   shorthand,
				Description: desc,
			})
		}
	}

	result := make(map[string][]FlagInfo)
	result["global"] = globalFlags
	result["instance"] = instanceFlags

	return result, nil
}

// UpdateReadme updates the README.md file with generated flag tables.
func UpdateReadme(readmePath string, flagInfos map[string][]FlagInfo) error {
	// Read the README file
	content, err := os.ReadFile(readmePath)
	if err != nil {
		return fmt.Errorf("failed to read README.md: %w", err)
	}

	// Generate flag tables
	globalTable := GenerateReadmeTable(flagInfos["global"], "Global Flags")
	instanceTable := GenerateReadmeTable(flagInfos["instance"], "Instance Command Flags")

	// Define regex patterns to find and replace the flag tables
	globalPattern := regexp.MustCompile(`(?s)### Global Flags\n\n\|.*?\n\|.*?\n(?:\|.*?\n)+`)
	instancePattern := regexp.MustCompile(`(?s)### Instance Command Flags\n\n\|.*?\n\|.*?\n(?:\|.*?\n)+`)

	// Replace the flag tables
	newContent := globalPattern.ReplaceAllString(string(content), globalTable)
	newContent = instancePattern.ReplaceAllString(newContent, instanceTable)

	// Write the updated content back to the README file
	return os.WriteFile(readmePath, []byte(newContent), 0644)
}
