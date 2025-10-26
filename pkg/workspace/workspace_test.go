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

	if ws.path == "" {
		t.Error("Workspace path is empty")
	}

	if ws.sessionPath == "" {
		t.Error("Session path is empty")
	}

	// Verify session file exists
	if _, err := os.Stat(ws.sessionPath); os.IsNotExist(err) {
		t.Errorf("Session file does not exist: %s", ws.sessionPath)
	}
}

func TestAppendCode(t *testing.T) {
	ws, err := New()
	if err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	code := `fmt.Println("test")`
	if err := ws.AppendCode(code); err != nil {
		t.Fatalf("Failed to append code: %v", err)
	}

	// Read the file and verify code was appended
	content, err := os.ReadFile(ws.sessionPath)
	if err != nil {
		t.Fatalf("Failed to read session file: %v", err)
	}

	if !strings.Contains(string(content), code) {
		t.Errorf("Code was not appended correctly. Expected to contain: %s", code)
	}
}

func TestLoadCode(t *testing.T) {
	ws, err := New()
	if err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	// Append some code
	testCode := `x := 42`
	if err := ws.AppendCode(testCode); err != nil {
		t.Fatalf("Failed to append code: %v", err)
	}

	// Load the code
	code, err := ws.LoadCode()
	if err != nil {
		t.Fatalf("Failed to load code: %v", err)
	}

	if !strings.Contains(code, testCode) {
		t.Errorf("Loaded code does not contain appended code. Got: %s", code)
	}
}

func TestClear(t *testing.T) {
	ws, err := New()
	if err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	// Append some code
	if err := ws.AppendCode(`x := 42`); err != nil {
		t.Fatalf("Failed to append code: %v", err)
	}

	// Clear the workspace
	if err := ws.Clear(); err != nil {
		t.Fatalf("Failed to clear workspace: %v", err)
	}

	// Load code and verify it's back to initial state
	code, err := ws.LoadCode()
	if err != nil {
		t.Fatalf("Failed to load code: %v", err)
	}

	// Should only contain package main and imports
	if strings.Contains(code, "x := 42") {
		t.Error("Code was not cleared properly")
	}

	if !strings.Contains(code, "package main") {
		t.Error("Cleared workspace should still contain package main")
	}
}

func TestSaveScript(t *testing.T) {
	ws, err := New()
	if err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	// Append some code
	testCode := `fmt.Println("test script")`
	if err := ws.AppendCode(testCode); err != nil {
		t.Fatalf("Failed to append code: %v", err)
	}

	// Save script
	scriptName := "test_script"
	if err := ws.SaveScript(scriptName); err != nil {
		t.Fatalf("Failed to save script: %v", err)
	}

	// Verify script file exists
	scriptPath := filepath.Join(ws.path, scriptName+".go")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		t.Errorf("Script file does not exist: %s", scriptPath)
	}

	// Verify script content
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read script file: %v", err)
	}

	if !strings.Contains(string(content), testCode) {
		t.Error("Script does not contain expected code")
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
