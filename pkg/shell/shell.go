package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Napolitain/gosh/pkg/workspace"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// Shell represents the interactive Go shell
type Shell struct {
	interpreter *interp.Interpreter
	workspace   *workspace.Workspace
	history     []string
}

// New creates a new Shell instance
func New() (*Shell, error) {
	ws, err := workspace.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create workspace: %w", err)
	}

	i := interp.New(interp.Options{})
	if err := i.Use(stdlib.Symbols); err != nil {
		return nil, fmt.Errorf("failed to load standard library: %w", err)
	}

	// Pre-import commonly used packages
	if _, err := i.Eval(`import "fmt"`); err != nil {
		return nil, fmt.Errorf("failed to import fmt: %w", err)
	}

	return &Shell{
		interpreter: i,
		workspace:   ws,
		history:     make([]string, 0),
	}, nil
}

// Run starts the interactive shell loop
func (s *Shell) Run() error {
	fmt.Println("Welcome to gosh - Go Shell")
	fmt.Println("Enter code blocks and press Ctrl+Enter to compile and execute")
	fmt.Println("Type 'help' for commands, 'exit' to quit")
	fmt.Println()

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println()
		s.promptForCLIGeneration()
		os.Exit(0)
	}()

	reader := bufio.NewReader(os.Stdin)
	
	for {
		// Read block-based input
		codeBlock, shouldExit, err := s.readCodeBlock(reader)
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				s.promptForCLIGeneration()
				return nil
			}
			return fmt.Errorf("failed to read input: %w", err)
		}

		if shouldExit {
			s.promptForCLIGeneration()
			return nil
		}

		codeBlock = strings.TrimSpace(codeBlock)
		
		if codeBlock == "" {
			continue
		}

		// Handle built-in commands
		if s.handleBuiltinCommand(codeBlock) {
			continue
		}

		// Add to history
		s.history = append(s.history, codeBlock)

		// Try to compile/execute the code
		if err := s.execute(codeBlock); err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Println("Code not added to project. Fix and try again.")
		} else {
			// If successful, add to workspace
			if err := s.workspace.AddCodeBlock(codeBlock); err != nil {
				fmt.Printf("Warning: failed to save code: %v\n", err)
			} else {
				fmt.Println("✓ Code compiled and added to project")
			}
		}
	}
}

// readCodeBlock reads a multi-line code block until Ctrl+Enter is pressed
func (s *Shell) readCodeBlock(reader *bufio.Reader) (string, bool, error) {
	var lines []string
	fmt.Print("gosh> ")
	
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", false, err
		}
		
		line = strings.TrimRight(line, "\n\r")
		
		// Check for exit command at the beginning
		if len(lines) == 0 && (line == "exit" || line == "quit") {
			return "", true, nil
		}
		
		// Check for special commands on first line
		if len(lines) == 0 && (strings.HasPrefix(line, "help") || 
			strings.HasPrefix(line, "history") || 
			strings.HasPrefix(line, "clear") || 
			strings.HasPrefix(line, "workspace")) {
			return line, false, nil
		}
		
		// Simple heuristic: if line ends with specific patterns, continue reading
		// Otherwise, treat as end of block (simulating Ctrl+Enter behavior)
		lines = append(lines, line)
		
		// Check if this looks like a complete statement
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || 
			strings.HasSuffix(trimmed, "{") || 
			strings.HasSuffix(trimmed, ",") ||
			strings.HasPrefix(trimmed, "//") {
			// Continue reading
			fmt.Print("... ")
			continue
		}
		
		// Check if we have an incomplete block
		blockText := strings.Join(lines, "\n")
		if s.isIncompleteBlock(blockText) {
			fmt.Print("... ")
			continue
		}
		
		// Otherwise, treat as complete block
		return blockText, false, nil
	}
}

// isIncompleteBlock checks if a code block is incomplete
func (s *Shell) isIncompleteBlock(block string) bool {
	openBraces := strings.Count(block, "{")
	closeBraces := strings.Count(block, "}")
	openParens := strings.Count(block, "(")
	closeParens := strings.Count(block, ")")
	
	return openBraces > closeBraces || openParens > closeParens
}

// handleBuiltinCommand handles shell built-in commands
func (s *Shell) handleBuiltinCommand(input string) bool {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return false
	}

	command := parts[0]

	switch command {
	case "exit", "quit":
		s.promptForCLIGeneration()
		os.Exit(0)
		return true

	case "help":
		s.printHelp()
		return true

	case "history":
		s.printHistory()
		return true

	case "clear":
		s.history = make([]string, 0)
		if err := s.workspace.Clear(); err != nil {
			fmt.Printf("Error clearing workspace: %v\n", err)
		} else {
			fmt.Println("History and workspace cleared")
		}
		return true

	case "workspace":
		fmt.Printf("Workspace directory: %s\n", s.workspace.Path())
		fmt.Printf("Internal directory: %s\n", s.workspace.InternalPath())
		fmt.Printf("Session ID: %s\n", s.workspace.SessionID())
		return true

	case "reload":
		// Reload workspace - recreate interpreter
		if err := s.reloadWorkspace(); err != nil {
			fmt.Printf("Error reloading workspace: %v\n", err)
		} else {
			fmt.Println("Workspace reloaded successfully")
		}
		return true

	default:
		return false
	}
}

// promptForCLIGeneration prompts the user to save session as a Cobra CLI tool
func (s *Shell) promptForCLIGeneration() {
	if len(s.workspace.GetCodeBlocks()) == 0 {
		fmt.Println("No code blocks to save. Exiting...")
		return
	}

	fmt.Print("\nWould you like to save this session as a CLI tool? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Exiting...")
		return
	}
	response = strings.TrimSpace(strings.ToLower(response))

	if response == "y" || response == "yes" {
		fmt.Print("Enter CLI tool name: ")
		name, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Exiting...")
			return
		}
		name = strings.TrimSpace(name)

		if name != "" {
			if err := s.workspace.GenerateCobraCLI(name); err != nil {
				fmt.Printf("Error generating CLI tool: %v\n", err)
			} else {
				fmt.Printf("✓ CLI tool '%s' generated successfully!\n", name)
				fmt.Printf("  Location: %s/cmd/%s/\n", s.workspace.Path(), name)
				fmt.Printf("  To build: cd %s/cmd/%s && go build\n", s.workspace.Path(), name)
			}
		}
	}

	fmt.Println("Exiting gosh...")
}

// execute runs the given Go code
func (s *Shell) execute(code string) error {
	_, err := s.interpreter.Eval(code)
	return err
}

// reloadWorkspace reloads workspace by creating a new interpreter
func (s *Shell) reloadWorkspace() error {
	// Create a new interpreter
	i := interp.New(interp.Options{})
	if err := i.Use(stdlib.Symbols); err != nil {
		return fmt.Errorf("failed to load standard library: %w", err)
	}

	// Pre-import commonly used packages
	if _, err := i.Eval(`import "fmt"`); err != nil {
		return fmt.Errorf("failed to import fmt: %w", err)
	}

	// Re-execute all code blocks
	for _, block := range s.workspace.GetCodeBlocks() {
		if _, err := i.Eval(block); err != nil {
			return fmt.Errorf("failed to evaluate code block: %w", err)
		}
	}

	s.interpreter = i
	return nil
}

// printHelp displays help information
func (s *Shell) printHelp() {
	fmt.Println("gosh - Go Shell Commands:")
	fmt.Println("  help        - Show this help message")
	fmt.Println("  history     - Show command history")
	fmt.Println("  clear       - Clear history and workspace")
	fmt.Println("  workspace   - Show workspace information")
	fmt.Println("  reload      - Reload workspace code")
	fmt.Println("  exit/quit   - Exit the shell (prompts to save as CLI tool)")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  - Type or paste Go code")
	fmt.Println("  - Press Enter to add new lines")
	fmt.Println("  - Code block is executed when it appears complete")
	fmt.Println("  - On exit, you can save your session as a Cobra-based CLI tool")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  fmt.Println(\"Hello, World!\")")
	fmt.Println("  x := 42")
	example := `  fmt.Printf("x = %d\n", x)`
	fmt.Println(example)
}

// printHistory displays command history
func (s *Shell) printHistory() {
	if len(s.history) == 0 {
		fmt.Println("No history")
		return
	}

	fmt.Println("Command history:")
	for i, cmd := range s.history {
		fmt.Printf("%4d  %s\n", i+1, cmd)
	}
}
