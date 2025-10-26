# Quick Start Guide

Get started with gosh (Go Shell) in minutes!

## Installation

### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/Napolitain/gosh.git
cd gosh

# Build the shell
go build -o gosh .

# Run gosh
./gosh
```

### Option 2: Go Install

```bash
go install github.com/Napolitain/gosh@latest
gosh
```

## First Steps

### 1. Start the Shell

```bash
./gosh
```

You'll see:
```
Welcome to gosh - Go Shell
Type 'help' for commands, 'exit' to quit

gosh>
```

### 2. Try Your First Command

```go
gosh> fmt.Println("Hello, gosh!")
Hello, gosh!
```

### 3. Work with Variables

```go
gosh> x := 42
gosh> y := 58
gosh> fmt.Printf("x + y = %d\n", x + y)
x + y = 100
```

### 4. Use Loops and Control Flow

```go
gosh> for i := 1; i <= 3; i++ {
...     fmt.Printf("Count: %d\n", i)
... }
Count: 1
Count: 2
Count: 3
```

### 5. Save Your Work

```go
gosh> save my_first_script
Script saved to: my_first_script.go

gosh> workspace
Workspace directory: /home/user/.gosh/sessions/20251026_143022
```

### 6. View Command History

```go
gosh> history
Command history:
   1  fmt.Println("Hello, gosh!")
   2  x := 42
   3  y := 58
   4  fmt.Printf("x + y = %d\n", x + y)
```

## Essential Commands

| Command | Description |
|---------|-------------|
| `help` | Show available commands |
| `history` | Display command history |
| `save <filename>` | Save session to a script file |
| `clear` | Clear history and workspace |
| `workspace` | Show workspace directory path |
| `reload` | Reload workspace code |
| `exit` or `quit` | Exit the shell |

## Tips & Tricks

### Multi-line Statements

You can write multi-line code directly:

```go
gosh> func greet(name string) {
...     fmt.Printf("Hello, %s!\n", name)
... }
gosh> greet("World")
Hello, World!
```

### Working with Slices and Maps

```go
gosh> fruits := []string{"apple", "banana", "cherry"}
gosh> for i, fruit := range fruits {
...     fmt.Printf("%d: %s\n", i+1, fruit)
... }
1: apple
2: banana
3: cherry

gosh> ages := map[string]int{"Alice": 30, "Bob": 25}
gosh> fmt.Printf("Ages: %v\n", ages)
Ages: map[Alice:30 Bob:25]
```

### Finding Your Session Files

All your commands are automatically saved to a session directory:

```bash
# Find all your sessions
ls -la ~/.gosh/sessions/

# Each session has a timestamped directory
# Inside you'll find:
# - session.go (all your commands)
# - any scripts you saved with 'save' command
```

## What's Next?

- Check out [examples documentation](../docs/examples/README.md) for more complex examples
- Read the [full README](../README.md) for detailed features
- Explore the [implementation details](../docs/IMPLEMENTATION.md)
- Contribute! See [CONTRIBUTING.md](../CONTRIBUTING.md)

## Getting Help

- Type `help` in the shell for available commands
- Check the [README.md](../README.md) for comprehensive documentation
- Open an issue on GitHub for bugs or feature requests

## Common Use Cases

### Prototyping and Experimentation

Use gosh to quickly test Go code snippets without creating full projects.

### Learning Go

Practice Go syntax and features in an interactive environment.

### Script Development

Build up complex scripts incrementally, testing each part as you go.

### Code Snippets

Save useful code snippets for later reuse.

## Troubleshooting

### Command Not Found

If you get "undefined" errors for standard functions like `fmt.Println()`, the fmt package should be auto-imported. If issues persist, try the `reload` command.

### Workspace Issues

If you have problems with the workspace, you can clear it:
```
gosh> clear
```

Or manually delete and recreate:
```bash
rm -rf ~/.gosh/sessions/
```

---

Enjoy using gosh! 🚀
