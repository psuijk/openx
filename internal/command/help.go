package command

import "fmt"

func helpHandler() {
	fmt.Println("Commands:")
	fmt.Println("\tadd\tCreate a project config")
	fmt.Println("\tlist\tList all projects")
	fmt.Println("\tshow\tPrint a project config")
	fmt.Println("\tedit\tOpen config in $EDITOR")
	fmt.Println("\tremove\tDelete a project config")
	fmt.Println("\trun\tOpen a project workspace")
	fmt.Println("\tversion\tPrint version")
}
