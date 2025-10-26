# Example Usage Sessions

This directory contains example usage scenarios for gosh.

## Hello World

```go
fmt.Println("Hello, World!")
fmt.Println("Welcome to gosh!")
```

## Variables and Operations

```go
// Variables
x := 42
y := 3.14
name := "gosh"

fmt.Printf("x = %d\n", x)
fmt.Printf("y = %.2f\n", y)
fmt.Printf("name = %s\n", name)

// Operations
sum := x + 10
fmt.Printf("x + 10 = %d\n", sum)

// String operations
greeting := fmt.Sprintf("Hello from %s!", name)
fmt.Println(greeting)
```

## Loops and Control Flow

```go
// For loop
fmt.Println("Counting to 5:")
for i := 1; i <= 5; i++ {
    fmt.Printf("%d ", i)
}
fmt.Println()

// Range over slice
fruits := []string{"apple", "banana", "cherry"}
fmt.Println("\nFruits:")
for i, fruit := range fruits {
    fmt.Printf("%d: %s\n", i+1, fruit)
}

// Conditional
x := 42
if x > 40 {
    fmt.Println("\nx is greater than 40")
} else {
    fmt.Println("\nx is not greater than 40")
}
```

## Functions

```go
// Define a function
add := func(a, b int) int {
    return a + b
}

result := add(5, 3)
fmt.Printf("5 + 3 = %d\n", result)

// Anonymous function
multiply := func(a, b int) int {
    return a * b
}

fmt.Printf("4 * 7 = %d\n", multiply(4, 7))
```

## Working with Slices and Maps

```go
// Slices
numbers := []int{1, 2, 3, 4, 5}
fmt.Printf("Numbers: %v\n", numbers)

// Append to slice
numbers = append(numbers, 6, 7)
fmt.Printf("After append: %v\n", numbers)

// Maps
ages := map[string]int{
    "Alice": 30,
    "Bob": 25,
}
fmt.Printf("Ages: %v\n", ages)

// Add to map
ages["Charlie"] = 35
fmt.Printf("After adding Charlie: %v\n", ages)
```

## Saving Your Work

After experimenting with code, save your session:

```
gosh> save my_experiment
Script saved to: my_experiment.go
```

The saved script will be in your workspace directory (check with `workspace` command).
