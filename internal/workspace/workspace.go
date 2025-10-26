package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultWorkspaceDir = ".gosh"
	internalDir         = "internal"
)

// Workspace manages the shell's working directory and code persistence
type Workspace struct {
	rootPath    string
	internalPath string
	sessionID   string
	codeBlocks  []string
}

// New creates a new workspace in the user's home directory
func New() (*Workspace, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	workspaceDir := filepath.Join(homeDir, defaultWorkspaceDir)
	
	// Create workspace directory if it doesn't exist
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create workspace directory: %w", err)
	}

	// Create internal directory for session code
	internalPath := filepath.Join(workspaceDir, internalDir)
	if err := os.MkdirAll(internalPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create internal directory: %w", err)
	}

	// Initialize go.mod if it doesn't exist
	goModPath := filepath.Join(workspaceDir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		modContent := "module gosh\n\ngo 1.25\n"
		if err := os.WriteFile(goModPath, []byte(modContent), 0644); err != nil {
			return nil, fmt.Errorf("failed to create go.mod: %w", err)
		}
	}

	// Create session ID
	sessionID := time.Now().Format("20060102_150405")

	return &Workspace{
		rootPath:    workspaceDir,
		internalPath: internalPath,
		sessionID:   sessionID,
		codeBlocks:  make([]string, 0),
	}, nil
}

// Path returns the workspace root directory path
func (w *Workspace) Path() string {
	return w.rootPath
}

// InternalPath returns the internal directory path
func (w *Workspace) InternalPath() string {
	return w.internalPath
}

// SessionID returns the current session ID
func (w *Workspace) SessionID() string {
	return w.sessionID
}

// AddCodeBlock adds a compiled code block to the workspace
func (w *Workspace) AddCodeBlock(code string) error {
	w.codeBlocks = append(w.codeBlocks, code)
	
	// Save to session file in internal/
	sessionFile := filepath.Join(w.internalPath, fmt.Sprintf("session_%s.go", w.sessionID))
	
	// Build file content
	content := "package internal\n\nimport (\n\t\"fmt\"\n)\n\n"
	for _, block := range w.codeBlocks {
		content += "// Block\n" + block + "\n\n"
	}
	
	if err := os.WriteFile(sessionFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}
	
	return nil
}

// GetCodeBlocks returns all code blocks from the current session
func (w *Workspace) GetCodeBlocks() []string {
	return w.codeBlocks
}

// Clear clears all code blocks
func (w *Workspace) Clear() error {
	w.codeBlocks = make([]string, 0)
	
	// Remove session file
	sessionFile := filepath.Join(w.internalPath, fmt.Sprintf("session_%s.go", w.sessionID))
	if err := os.Remove(sessionFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove session file: %w", err)
	}
	
	return nil
}

// GenerateCobraCLI generates a Cobra-based CLI tool from the session code
func (w *Workspace) GenerateCobraCLI(name string) error {
	if name == "" {
		return fmt.Errorf("CLI name cannot be empty")
	}
	
	// Create CLI directory
	cliDir := filepath.Join(w.rootPath, "cmd", name)
	if err := os.MkdirAll(cliDir, 0755); err != nil {
		return fmt.Errorf("failed to create CLI directory: %w", err)
	}
	
	// Generate main.go with Cobra
	mainContent := fmt.Sprintf(`package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "%s",
	Short: "Generated CLI from gosh session %s",
	Run: func(cmd *cobra.Command, args []string) {
		// Session code
%s
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
`, name, w.sessionID, w.formatCodeBlocksForCLI())
	
	mainPath := filepath.Join(cliDir, "main.go")
	if err := os.WriteFile(mainPath, []byte(mainContent), 0644); err != nil {
		return fmt.Errorf("failed to write main.go: %w", err)
	}
	
	return nil
}

// formatCodeBlocksForCLI formats code blocks for inclusion in CLI tool
func (w *Workspace) formatCodeBlocksForCLI() string {
	var result strings.Builder
	for _, block := range w.codeBlocks {
		lines := strings.Split(block, "\n")
		for _, line := range lines {
			if line != "" {
				result.WriteString("\t\t")
				result.WriteString(line)
				result.WriteString("\n")
			}
		}
	}
	return result.String()
}
