package main

import (
	"fmt"
	"os"

	"github.com/psuijk/openx/internal/command"

	_ "github.com/psuijk/openx/internal/backend/cmux"
)

func main() {
	err := command.Dispatch(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
