# gosh - Go Shell

[![CI](https://github.com/Napolitain/gosh/actions/workflows/ci.yml/badge.svg)](https://github.com/Napolitain/gosh/actions/workflows/ci.yml)
[![Release](https://github.com/Napolitain/gosh/actions/workflows/auto-release.yml/badge.svg)](https://github.com/Napolitain/gosh/actions/workflows/auto-release.yml)

Go-based shell which works cross-platform (Linux, Windows), is interactive and intuitive to use.

## Features

- **Block-Based Input**: Write multi-line Go code naturally - press Enter for new lines, code executes when block is complete
- **Smart Compilation**: Code only added to project when it compiles successfully
- **Hot Reload**: Leverages yaegi interpreter for instant code execution without compilation overhead
- **Workspace as Monorepo**: Maintains a Go monorepo in `~/.gosh/` with session code in `internal/` directory
- **CLI Tool Generation**: On exit, converts your session into a Cobra-based CLI tool
- **Cross-Platform**: Works seamlessly on Linux, Windows, and macOS
- **Command History**: Track and review your command history
- **Session Persistence**: All session code saved to `~/.gosh/internal/session_TIMESTAMP.go`

## Installation

### Download Pre-built Binary

Download the latest release from the [Releases](https://github.com/Napolitain/gosh/releases) page:

- **Linux**: `gosh-linux-amd64`
- **macOS Intel**: `gosh-darwin-amd64`
- **macOS Apple Silicon**: `gosh-darwin-arm64`
- **Windows**: `gosh-windows-amd64.exe`

### From Source

```bash
git clone https://github.com/Napolitain/gosh.git
cd gosh
go build -o gosh .
```

### Run

```bash
./gosh
```

## Usage

### Basic Commands

Start the shell:
```bash
./gosh
```

### Shell Commands

- `help` - Show available commands
- `history` - Display command history
- `clear` - Clear history and workspace
- `workspace` - Show workspace information (path, internal path, session ID)
- `reload` - Reload workspace code
- `exit` or `quit` - Exit the shell (prompts to save as CLI tool)

### Example Session

```
Welcome to gosh - Go Shell
Write multi-line code blocks - press Enter for new lines
Press Ctrl+Enter (Cmd+Enter on Mac) to execute your code block
Type 'help' for commands, 'exit' to quit

gosh> fmt.Println("Hello, World!")
Hello, World!
✓ Code compiled and added to project

gosh> x := 42
✓ Code compiled and added to project

gosh> y := 58
✓ Code compiled and added to project

gosh> total := x + y
✓ Code compiled and added to project

gosh> fmt.Printf("Total: %d\n", total)
Total: 100
✓ Code compiled and added to project

gosh> for i := 0; i < 3; i++ {
...     fmt.Printf("Iteration %d\n", i)
... }
Iteration 0
Iteration 1
Iteration 2
✓ Code compiled and added to project

gosh> history
Command history:
   1  fmt.Println("Hello, World!")
   2  x := 42
   3  y := 58
   4  total := x + y
   5  fmt.Printf("Total: %d\n", total)
   6  for i := 0; i < 3; i++ { fmt.Printf("Iteration %d\n", i) }

gosh> workspace
Workspace directory: /home/user/.gosh
Internal directory: /home/user/.gosh/internal
Session ID: 20251026_143022

gosh> exit

Would you like to save this session as a CLI tool? (y/n): y
Enter CLI tool name: hello_tool
✓ CLI tool 'hello_tool' generated successfully!
  Location: /home/user/.gosh/cmd/hello_tool/
  To build: cd /home/user/.gosh/cmd/hello_tool && go build
Exiting gosh...
```

## How It Works

### Block-Based Input

gosh provides true block-based input:
- Press **Enter** to add new lines within your code block
- Press **Ctrl+Enter** (**Cmd+Enter** on Mac) to execute your code block
- Perfect for writing multi-line Go code naturally

### Smart Compilation

- Code is compiled and executed when you complete a block
- **Only on success** is code added to the project workspace
- Failed compilation shows errors without corrupting your project
- Success shows "✓ Code compiled and added to project"

### Hot Reload

gosh uses the [yaegi](https://github.com/traefik/yaegi) Go interpreter, which provides instant code execution without traditional compilation. This allows for:

- Immediate feedback on code execution
- Interactive development and experimentation
- Fast iteration cycles

### Workspace as Monorepo

gosh maintains a Go monorepo structure in `~/.gosh/`:

```
~/.gosh/
├── go.mod              # Module definition
├── internal/           # Session code
│   └── session_TIMESTAMP.go
└── cmd/                # Generated CLI tools
    └── <tool_name>/
        └── main.go
```

Each session creates a file in `internal/` with all successfully compiled code blocks. This structure allows you to:
- Maintain a clean Go module
- Easily reference code across sessions
- Build complete applications from sessions

### CLI Tool Generation

When exiting gosh, you can save your session as a Cobra-based CLI tool:

1. Exit with `exit` or `quit` command
2. Answer "y" when prompted to save as CLI tool
3. Provide a name for your tool
4. Tool is generated in `~/.gosh/cmd/<name>/`
5. Build with: `cd ~/.gosh/cmd/<name> && go build`

The generated CLI reproduces all your session code, making it easy to share or deploy your experiments.

## Architecture

```
gosh/
├── main.go                 # Entry point
├── internal/
│   ├── shell/             # Shell REPL implementation
│   │   ├── shell.go       # Block-based input & execution
│   │   └── shell_test.go
│   └── workspace/         # Workspace management
│       ├── workspace.go   # Monorepo & CLI generation
│       └── workspace_test.go
└── ~/.gosh/               # User workspace (created at runtime)
    ├── go.mod             # Go module definition
    ├── internal/          # Session code storage
    │   └── session_TIMESTAMP.go
    └── cmd/               # Generated CLI tools
        └── <tool_name>/
            └── main.go
```

## Development

### Building

```bash
go build -o gosh .
```

### Running Tests

```bash
go test ./... -v
```

### Dependencies

- [yaegi](https://github.com/traefik/yaegi) - Go interpreter for hot reload functionality
- [cobra](https://github.com/spf13/cobra) - CLI framework for generated tools

## Future Enhancements

- Syntax highlighting in terminal
- Tab completion for Go keywords and functions
- Integration with external Go packages
- Import management UI
- Configuration file support
- Multi-user workspace support

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.
