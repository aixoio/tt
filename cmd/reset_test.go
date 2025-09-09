package cmd

import (
	"testing"
)

func TestResetCmd(t *testing.T) {
	// Test that the command is properly initialized
	if resetCmd == nil {
		t.Fatal("resetCmd should not be nil")
	}

	if resetCmd.Use != "reset" {
		t.Errorf("expected Use to be 'reset', got %s", resetCmd.Use)
	}

	if resetCmd.Short != "Hard-reset the repository after confirmation" {
		t.Errorf("expected Short to match, got %s", resetCmd.Short)
	}
}

func TestResetCmdFlags(t *testing.T) {
	// Test that the command has no custom flags
	flags := resetCmd.Flags()
	if flags == nil {
		t.Fatal("flags should not be nil")
	}

	// Check that no custom flags are defined (only inherited ones)
	definedFlags := flags.NFlag()
	if definedFlags != 0 {
		t.Errorf("expected 0 custom flags, got %d", definedFlags)
	}
}

func TestResetCmdParent(t *testing.T) {
	// Test that resetCmd is added to rootCmd
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd == resetCmd {
			found = true
			break
		}
	}

	if !found {
		t.Error("resetCmd should be added to rootCmd")
	}
}

// Note: Full integration testing of the interactive prompt and git commands
// would require mocking exec.Command and simulating user input, which is
// complex for unit tests. Consider integration tests for end-to-end validation.
