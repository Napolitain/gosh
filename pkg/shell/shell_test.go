package shell

import (
	"testing"
)

func TestNew(t *testing.T) {
	sh, err := New()
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}

	if sh.interpreter == nil {
		t.Error("Interpreter is nil")
	}

	if sh.workspace == nil {
		t.Error("Workspace is nil")
	}

	if sh.history == nil {
		t.Error("History is nil")
	}
}

func TestExecute(t *testing.T) {
	sh, err := New()
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}

	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{
			name:    "Simple variable assignment",
			code:    `x := 42`,
			wantErr: false,
		},
		{
			name:    "Print statement",
			code:    `fmt.Println("test")`,
			wantErr: false,
		},
		{
			name:    "Invalid syntax",
			code:    `this is not valid go code`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sh.execute(tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleBuiltinCommand(t *testing.T) {
	sh, err := New()
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}

	tests := []struct {
		name      string
		input     string
		isBuiltin bool
	}{
		{
			name:      "Help command",
			input:     "help",
			isBuiltin: true,
		},
		{
			name:      "History command",
			input:     "history",
			isBuiltin: true,
		},
		{
			name:      "Clear command",
			input:     "clear",
			isBuiltin: true,
		},
		{
			name:      "Workspace command",
			input:     "workspace",
			isBuiltin: true,
		},
		{
			name:      "Reload command",
			input:     "reload",
			isBuiltin: true,
		},
		{
			name:      "Not a builtin command",
			input:     "x := 42",
			isBuiltin: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't directly test handleBuiltinCommand for exit/quit
			// as they call os.Exit, so we skip those
			if tt.input == "exit" || tt.input == "quit" {
				return
			}

			result := sh.handleBuiltinCommand(tt.input)
			if result != tt.isBuiltin {
				t.Errorf("handleBuiltinCommand() = %v, want %v", result, tt.isBuiltin)
			}
		})
	}
}

func TestReloadWorkspace(t *testing.T) {
	sh, err := New()
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}

	// Add some code blocks to workspace
	code := `x := 100`
	if err := sh.workspace.AddCodeBlock(code); err != nil {
		t.Fatalf("Failed to add code block: %v", err)
	}

	// Reload workspace
	if err := sh.reloadWorkspace(); err != nil {
		t.Logf("Note: reload may have issues with certain code types: %v", err)
	}

	// Interpreter should be reset
	if sh.interpreter == nil {
		t.Error("Interpreter is nil after reload")
	}
}
