# gosh - Go Shell

Go-based shell which works cross-platform (Linux, Windows), is interactive and intuitive to use.

## Features

- **Interactive REPL**: Execute Go code interactively with immediate feedback
- **Hot Reload**: Leverages Go's compilation through the yaegi interpreter for instant code execution
- **Session Management**: Automatically saves your commands as a Go project in a workspace folder
- **Script Saving**: Save your interactive session as a reusable Go script
- **Cross-Platform**: Works seamlessly on Linux, Windows, and macOS
- **Command History**: Track and review your command history
- **Workspace Persistence**: All session code is saved to `~/.gosh/sessions/` for later reference

## Installation

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
- `save <filename>` - Save current session to a script file
- `clear` - Clear history and workspace
- `workspace` - Show workspace directory path
- `reload` - Reload workspace code
- `exit` or `quit` - Exit the shell

### Example Session

```
gosh> fmt.Println("Hello, World!")
Hello, World!

gosh> x := 42
gosh> fmt.Printf("x = %d\n", x)
x = 42

gosh> for i := 0; i < 3; i++ {
...     fmt.Printf("Iteration %d\n", i)
... }
Iteration 0
Iteration 1
Iteration 2

gosh> save my_script
Script saved to: my_script.go

gosh> history
Command history:
   1  fmt.Println("Hello, World!")
   2  x := 42
   3  fmt.Printf("x = %d\n", x)
   4  for i := 0; i < 3; i++ { fmt.Printf("Iteration %d\n", i) }

gosh> exit
```

## How It Works

### Hot Reload

gosh uses the [yaegi](https://github.com/traefik/yaegi) Go interpreter, which provides instant code execution without traditional compilation. This allows for:

- Immediate feedback on code execution
- Interactive development and experimentation
- Fast iteration cycles

### Session Persistence

Every time you start gosh, it creates a new session directory in `~/.gosh/sessions/` with a timestamp. All commands you execute are automatically appended to a `session.go` file in that directory, making it easy to:

- Review what you did in a session
- Save successful experiments as scripts
- Build up working code incrementally

### Script Saving

Use the `save` command to export your current session to a named Go script file. The script will be saved in your current session's workspace directory.

## Architecture

```
gosh/
├── main.go                 # Entry point
├── pkg/
│   ├── shell/             # Shell REPL implementation
│   │   └── shell.go
│   └── workspace/         # Session and file management
│       └── workspace.go
└── ~/.gosh/               # User workspace (created at runtime)
    └── sessions/          # Session directories
        └── YYYYMMDD_HHMMSS/
            ├── session.go  # Auto-saved session code
            └── *.go        # Saved scripts
```

## Development

### Building

```bash
go build -o gosh .
```

### Dependencies

- [yaegi](https://github.com/traefik/yaegi) - Go interpreter for hot reload functionality

## Future Enhancements

- Tab completion
- Syntax highlighting
- Multi-line input support with proper formatting
- Integration with external Go packages
- Enhanced error messages and debugging
- Configuration file support

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.
