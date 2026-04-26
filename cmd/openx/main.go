package main

import (
	"fmt"
	"os"

	"github.com/psuijk/openx/internal/command"
)

func main() {
	err := command.Dispatch(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
