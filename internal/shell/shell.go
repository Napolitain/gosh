package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Napolitain/gosh/internal/workspace"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"golang.org/x/term"
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
	fmt.Println("Write multi-line code blocks - press Enter for new lines")
	fmt.Println("Press Ctrl+D (Cmd+D on Mac) to execute your code block")
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

// readCodeBlock reads a multi-line code block
// Press Enter for new lines, Ctrl+D (Cmd+D on Mac) to submit
func (s *Shell) readCodeBlock(reader *bufio.Reader) (string, bool, error) {
	fmt.Print("gosh> ")
	
	// Check if stdin is a terminal
	fd := int(os.Stdin.Fd())
	isTerminal := term.IsTerminal(fd)
	
	if isTerminal {
		// Use raw mode for better control
		return s.readCodeBlockRaw()
	}
	
	// Fallback for non-terminal (pipes, redirects, etc.)
	return s.readCodeBlockBuffered(reader)
}

// readCodeBlockRaw reads input in raw terminal mode with Ctrl+D submit
func (s *Shell) readCodeBlockRaw() (string, bool, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		// Fall back to buffered if raw mode fails
		return s.readCodeBlockBuffered(bufio.NewReader(os.Stdin))
	}
	defer term.Restore(fd, oldState)
	
	var buffer []rune
	lineStart := 0 // Track start of current line for backspace
	
	buf := make([]byte, 3) // Read up to 3 bytes for potential escape sequences
	
	for {
		n, err := os.Stdin.Read(buf[:1])
		if err != nil {
			term.Restore(fd, oldState)
			if err == io.EOF {
				return "", false, err
			}
			return "", false, err
		}
		
		if n == 0 {
			continue
		}
		
		ch := buf[0]
		
		switch ch {
		case 3: // Ctrl+C
			fmt.Print("^C\r\n")
			term.Restore(fd, oldState)
			return "", false, io.EOF
			
		case 4: // Ctrl+D - Submit block or EOF
			if len(buffer) == 0 {
				// Empty buffer - treat as EOF
				fmt.Print("^D\r\n")
				term.Restore(fd, oldState)
				return "", false, io.EOF
			}
			// Submit the current block
			fmt.Print("\r\n")
			term.Restore(fd, oldState)
			result := string(buffer)
			result = strings.TrimSpace(result)
			
			// Check for exit commands
			if result == "exit" || result == "quit" {
				return "", true, nil
			}
			
			return result, false, nil
			
		case 13, 10: // Enter (CR or LF)
			// Add newline to buffer
			buffer = append(buffer, '\n')
			lineStart = len(buffer)
			fmt.Print("\r\n...  ")
			
		case 127, 8: // Backspace or DEL
			if len(buffer) > lineStart {
				buffer = buffer[:len(buffer)-1]
				fmt.Print("\b \b")
			}
			
		case 9: // Tab
			// Insert 4 spaces for tab
			buffer = append(buffer, ' ', ' ', ' ', ' ')
			fmt.Print("    ")
			
		case 27: // ESC - potential escape sequence
			// Try to read more bytes for escape sequences
			// Set a short timeout or just skip for now
			// For simplicity, we'll ignore escape sequences
			continue
			
		default:
			if ch >= 32 && ch < 127 {
				// Printable ASCII character
				buffer = append(buffer, rune(ch))
				fmt.Printf("%c", ch)
			}
			// Ignore other control characters
		}
	}
}

// readCodeBlockBuffered reads input using buffered reader (fallback mode)
// Uses empty line to submit
func (s *Shell) readCodeBlockBuffered(reader *bufio.Reader) (string, bool, error) {
	var lines []string
	
	firstLine := true
	for {
		var line string
		var err error
		
		line, err = reader.ReadString('\n')
		if err != nil {
			return "", false, err
		}
		
		line = strings.TrimRight(line, "\n\r")
		
		// Check for exit command at the beginning
		if firstLine && (line == "exit" || line == "quit") {
			return "", true, nil
		}
		
		// Check for special commands on first line
		if firstLine && (strings.HasPrefix(line, "help") || 
			strings.HasPrefix(line, "history") || 
			strings.HasPrefix(line, "clear") || 
			strings.HasPrefix(line, "workspace") ||
			strings.HasPrefix(line, "reload")) {
			return line, false, nil
		}
		
		firstLine = false
		
		// Empty line submits the block
		if strings.TrimSpace(line) == "" {
			if len(lines) > 0 {
				// Submit the accumulated block
				return strings.Join(lines, "\n"), false, nil
			}
			// Empty input, start over
			fmt.Print("gosh> ")
			firstLine = true
			continue
		}
		
		// Add line to the block
		lines = append(lines, line)
		fmt.Print("...  ")
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
	fmt.Println("  - Type or paste multi-line Go code")
	fmt.Println("  - Press Enter to add new lines within your code block")
	fmt.Println("  - Press Ctrl+D (Cmd+D on Mac) to execute the code block")
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
