package command

import "fmt"

// Version is the current version of openx.
const Version = "0.1.0"

func versionHandler() {
	fmt.Printf("version %s\n", Version)
}
