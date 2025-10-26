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
	sessionFileName     = "session.go"
)

// Workspace manages the shell's working directory and code persistence
type Workspace struct {
	path        string
	sessionPath string
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

	// Create session file path with timestamp
	sessionID := time.Now().Format("20060102_150405")
	sessionDir := filepath.Join(workspaceDir, "sessions", sessionID)
	
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create session directory: %w", err)
	}

	sessionPath := filepath.Join(sessionDir, sessionFileName)

	// Initialize session file with package main
	initContent := "package main\n\nimport (\n\t\"fmt\"\n)\n\n"
	if err := os.WriteFile(sessionPath, []byte(initContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to initialize session file: %w", err)
	}

	return &Workspace{
		path:        sessionDir,
		sessionPath: sessionPath,
	}, nil
}

// Path returns the workspace directory path
func (w *Workspace) Path() string {
	return w.path
}

// AppendCode appends code to the session file
func (w *Workspace) AppendCode(code string) error {
	// Open file in append mode
	f, err := os.OpenFile(w.sessionPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open session file: %w", err)
	}
	defer f.Close()

	// Add newline if code doesn't end with one
	if !strings.HasSuffix(code, "\n") {
		code += "\n"
	}

	if _, err := f.WriteString(code + "\n"); err != nil {
		return fmt.Errorf("failed to write to session file: %w", err)
	}

	return nil
}

// LoadCode loads all code from the session file
func (w *Workspace) LoadCode() (string, error) {
	data, err := os.ReadFile(w.sessionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("failed to read session file: %w", err)
	}

	return string(data), nil
}

// Clear clears the session file
func (w *Workspace) Clear() error {
	// Reinitialize with package main
	initContent := "package main\n\nimport (\n\t\"fmt\"\n)\n\n"
	if err := os.WriteFile(w.sessionPath, []byte(initContent), 0644); err != nil {
		return fmt.Errorf("failed to clear session file: %w", err)
	}

	return nil
}

// SaveScript saves the current session to a named script file
func (w *Workspace) SaveScript(filename string) error {
	// Ensure filename has .go extension
	if !strings.HasSuffix(filename, ".go") {
		filename += ".go"
	}

	destPath := filepath.Join(w.path, filename)

	// Read current session
	data, err := os.ReadFile(w.sessionPath)
	if err != nil {
		return fmt.Errorf("failed to read session file: %w", err)
	}

	// Write to destination
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write script file: %w", err)
	}

	return nil
}
