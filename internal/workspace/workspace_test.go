package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	ws, err := New()
	if err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	if ws.Path() == "" {
		t.Error("Workspace path is empty")
	}

	if ws.InternalPath() == "" {
		t.Error("Internal path is empty")
	}

	if ws.SessionID() == "" {
		t.Error("Session ID is empty")
	}

	// Verify internal directory exists
	if _, err := os.Stat(ws.InternalPath()); os.IsNotExist(err) {
		t.Errorf("Internal directory does not exist: %s", ws.InternalPath())
	}

	// Verify go.mod exists
	goModPath := filepath.Join(ws.Path(), "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Errorf("go.mod does not exist: %s", goModPath)
	}
}

func TestAddCodeBlock(t *testing.T) {
	ws, err := New()
	if err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	code := `fmt.Println("test")`
	if err := ws.AddCodeBlock(code); err != nil {
		t.Fatalf("Failed to add code block: %v", err)
	}

	// Verify code block was added
	blocks := ws.GetCodeBlocks()
	if len(blocks) != 1 {
		t.Errorf("Expected 1 code block, got %d", len(blocks))
	}

	if blocks[0] != code {
		t.Errorf("Code block mismatch. Expected: %s, Got: %s", code, blocks[0])
	}

	// Verify session file was created
	sessionFile := filepath.Join(ws.InternalPath(), "session_"+ws.SessionID()+".go")
	if _, err := os.Stat(sessionFile); os.IsNotExist(err) {
		t.Errorf("Session file does not exist: %s", sessionFile)
	}
}

func TestGetCodeBlocks(t *testing.T) {
	ws, err := New()
	if err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	// Add multiple code blocks
	blocks := []string{
		`x := 42`,
		`fmt.Println(x)`,
		`y := "test"`,
	}

	for _, block := range blocks {
		if err := ws.AddCodeBlock(block); err != nil {
			t.Fatalf("Failed to add code block: %v", err)
		}
	}

	// Get code blocks
	retrieved := ws.GetCodeBlocks()
	if len(retrieved) != len(blocks) {
		t.Errorf("Expected %d code blocks, got %d", len(blocks), len(retrieved))
	}

	for i, block := range blocks {
		if retrieved[i] != block {
			t.Errorf("Block %d mismatch. Expected: %s, Got: %s", i, block, retrieved[i])
		}
	}
}

func TestClear(t *testing.T) {
	ws, err := New()
	if err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	// Add some code blocks
	if err := ws.AddCodeBlock(`x := 42`); err != nil {
		t.Fatalf("Failed to add code block: %v", err)
	}

	// Clear the workspace
	if err := ws.Clear(); err != nil {
		t.Fatalf("Failed to clear workspace: %v", err)
	}

	// Verify code blocks were cleared
	blocks := ws.GetCodeBlocks()
	if len(blocks) != 0 {
		t.Errorf("Expected 0 code blocks after clear, got %d", len(blocks))
	}

	// Verify session file was removed
	sessionFile := filepath.Join(ws.InternalPath(), "session_"+ws.SessionID()+".go")
	if _, err := os.Stat(sessionFile); !os.IsNotExist(err) {
		t.Error("Session file should not exist after clear")
	}
}

func TestGenerateCobraCLI(t *testing.T) {
	ws, err := New()
	if err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	// Add some code blocks
	testCode := `fmt.Println("test CLI")`
	if err := ws.AddCodeBlock(testCode); err != nil {
		t.Fatalf("Failed to add code block: %v", err)
	}

	// Generate CLI
	cliName := "test_cli"
	if err := ws.GenerateCobraCLI(cliName); err != nil {
		t.Fatalf("Failed to generate CLI: %v", err)
	}

	// Verify CLI directory and main.go exist
	cliDir := filepath.Join(ws.Path(), "cmd", cliName)
	if _, err := os.Stat(cliDir); os.IsNotExist(err) {
		t.Errorf("CLI directory does not exist: %s", cliDir)
	}

	mainPath := filepath.Join(cliDir, "main.go")
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		t.Errorf("main.go does not exist: %s", mainPath)
	}

	// Verify main.go content
	content, err := os.ReadFile(mainPath)
	if err != nil {
		t.Fatalf("Failed to read main.go: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "package main") {
		t.Error("main.go should contain 'package main'")
	}

	if !strings.Contains(contentStr, "github.com/spf13/cobra") {
		t.Error("main.go should import cobra")
	}

	if !strings.Contains(contentStr, testCode) {
		t.Error("main.go should contain the test code")
	}
}

func TestPath(t *testing.T) {
	ws, err := New()
	if err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	path := ws.Path()
	if path == "" {
		t.Error("Path() returned empty string")
	}

	// Verify path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Workspace path does not exist: %s", path)
	}
}

func TestInternalPath(t *testing.T) {
	ws, err := New()
	if err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	internalPath := ws.InternalPath()
	if internalPath == "" {
		t.Error("InternalPath() returned empty string")
	}

	// Verify path exists
	if _, err := os.Stat(internalPath); os.IsNotExist(err) {
		t.Errorf("Internal path does not exist: %s", internalPath)
	}

	// Verify it's a subdirectory of workspace
	if !strings.HasPrefix(internalPath, ws.Path()) {
		t.Error("Internal path should be a subdirectory of workspace path")
	}
}

func TestSessionID(t *testing.T) {
	ws, err := New()
	if err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	sessionID := ws.SessionID()
	if sessionID == "" {
		t.Error("SessionID() returned empty string")
	}

	// Verify session ID format (YYYYMMDD_HHMMSS)
	if len(sessionID) != 15 {
		t.Errorf("Session ID should be 15 characters long, got %d", len(sessionID))
	}
}
