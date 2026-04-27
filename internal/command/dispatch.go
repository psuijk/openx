package command

// Dispatch routes command-line arguments to the appropriate subcommand handler.
func Dispatch(args []string) error {
	if len(args) == 0 {
		helpHandler()
		return nil
	}

	switch args[0] {
	case "add":
		return addHandler(args[1:])
	case "list":
		return listHandler()
	case "show":
		return showHandler(args[1:])
	case "edit":
		return editHandler(args[1:])
	case "remove":
		return removeHandler(args[1:])
	case "help", "-h", "--help", "":
		helpHandler()
	case "version":
		versionHandler()
	default:
		return runHandler(args[1:])
	}

	return nil
}
