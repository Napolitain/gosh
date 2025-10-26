package main

import (
	"fmt"
	"os"

	"github.com/Napolitain/gosh/pkg/shell"
)

func main() {
	sh, err := shell.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing shell: %v\n", err)
		os.Exit(1)
	}

	if err := sh.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running shell: %v\n", err)
		os.Exit(1)
	}
}
