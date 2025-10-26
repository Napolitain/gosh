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
	reader      *bufio.Reader
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
		reader:      bufio.NewReader(os.Stdin),
	}, nil
}

// Run starts the interactive shell loop
func (s *Shell) Run() error {
	fmt.Println("Welcome to gosh - Go Shell")
	fmt.Println("Type 'help' for commands, 'exit' to quit")
	fmt.Println()

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nExiting gosh...")
		os.Exit(0)
	}()

	for {
		fmt.Print("gosh> ")
		
		input, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nExiting gosh...")
				return nil
			}
			return fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)
		
		if input == "" {
			continue
		}

		// Handle built-in commands
		if s.handleBuiltinCommand(input) {
			continue
		}

		// Add to history
		s.history = append(s.history, input)

		// Save to workspace
		if err := s.workspace.AppendCode(input); err != nil {
			fmt.Printf("Warning: failed to save code: %v\n", err)
		}

		// Execute the code
		if err := s.execute(input); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
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
		fmt.Println("Exiting gosh...")
		os.Exit(0)
		return true

	case "help":
		s.printHelp()
		return true

	case "history":
		s.printHistory()
		return true

	case "save":
		if len(parts) < 2 {
			fmt.Println("Usage: save <filename>")
			return true
		}
		if err := s.workspace.SaveScript(parts[1]); err != nil {
			fmt.Printf("Error saving script: %v\n", err)
		} else {
			fmt.Printf("Script saved to: %s\n", parts[1])
		}
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
		return true

	case "reload":
		// Reload all workspace code
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

// execute runs the given Go code
func (s *Shell) execute(code string) error {
	_, err := s.interpreter.Eval(code)
	return err
}

// reloadWorkspace reloads all code from the workspace
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

	// Load workspace code
	code, err := s.workspace.LoadCode()
	if err != nil {
		return fmt.Errorf("failed to load workspace code: %w", err)
	}

	if code != "" {
		if _, err := i.Eval(code); err != nil {
			return fmt.Errorf("failed to evaluate workspace code: %w", err)
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
	fmt.Println("  save <file> - Save current session to a script file")
	fmt.Println("  clear       - Clear history and workspace")
	fmt.Println("  workspace   - Show workspace directory")
	fmt.Println("  reload      - Reload workspace code")
	fmt.Println("  exit/quit   - Exit the shell")
	fmt.Println()
	fmt.Println("You can type any valid Go code and it will be executed.")
	fmt.Println("Examples:")
	fmt.Println("  fmt.Println(\"Hello, World!\")")
	fmt.Println("  x := 42")
	fmt.Println("  fmt.Printf(\"x = %d\\n\", x)")
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
