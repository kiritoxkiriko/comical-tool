package main

import (
	"fmt"
	"os"

	"github.com/kiritoxkiriko/comical-tool/cli/internal/command"
)

func main() {
	if err := command.NewRoot().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
