package plan

// Step represents a single command to execute, with a human-readable description.
type Step struct {
	Command     string
	Description string
}

// Plan represents the ordered sequence of steps to open a project workspace.
type Plan struct {
	PreOpen  []Step
	Backend  []Step
	PostOpen []Step
}
