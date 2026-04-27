package command

import "testing"

func TestDispatch_NoArgs(t *testing.T) {
	err := Dispatch([]string{})
	if err != nil {
		t.Fatalf("expected no error for empty args (help), got: %v", err)
	}
}

func TestDispatch_Help(t *testing.T) {
	for _, arg := range []string{"help", "-h", "--help"} {
		err := Dispatch([]string{arg})
		if err != nil {
			t.Errorf("Dispatch(%q) returned error: %v", arg, err)
		}
	}
}

func TestDispatch_Version(t *testing.T) {
	err := Dispatch([]string{"version"})
	if err != nil {
		t.Fatalf("expected no error for version, got: %v", err)
	}
}

func TestDispatch_UnknownCommand(t *testing.T) {
	// Unknown commands route to runHandler, which prints "not implemented"
	err := Dispatch([]string{"nonexistent"})
	if err != nil {
		t.Fatalf("expected no error for unknown command (routes to run), got: %v", err)
	}
}
