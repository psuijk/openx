package command

import "fmt"

const Version = "0.1.0"

func versionHandler() {
	fmt.Printf("version %s\n", Version)
}
