package cmd

import (
	"testing"
)

func TestAddCmd(t *testing.T) {
	// Test that the command is properly initialized
	if addCmd == nil {
		t.Fatal("addCmd should not be nil")
	}

	if addCmd.Use != "add [paths...]" {
		t.Errorf("expected Use to be 'add [paths...]', got %s", addCmd.Use)
	}

	if addCmd.Short != "Stage files for commit" {
		t.Errorf("expected Short to match, got %s", addCmd.Short)
	}

	// Test aliases
	if len(addCmd.Aliases) != 1 || addCmd.Aliases[0] != "ad" {
		t.Errorf("expected alias 'ad', got %v", addCmd.Aliases)
	}
}

func TestAddCmdFlags(t *testing.T) {
	// Test that the command has the expected flags
	flags := addCmd.Flags()
	if flags == nil {
		t.Fatal("flags should not be nil")
	}

	// Check for --all flag
	allFlag := flags.Lookup("all")
	if allFlag == nil {
		t.Error("expected --all flag to be defined")
	} else {
		if allFlag.Name != "all" {
			t.Errorf("expected flag name 'all', got %s", allFlag.Name)
		}
	}

	// Check for --path flag
	pathFlag := flags.Lookup("path")
	if pathFlag == nil {
		t.Error("expected --path flag to be defined")
	} else {
		if pathFlag.Name != "path" {
			t.Errorf("expected flag name 'path', got %s", pathFlag.Name)
		}
	}
}

func TestAddCmdParent(t *testing.T) {
	// Test that addCmd is added to rootCmd
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd == addCmd {
			found = true
			break
		}
	}

	if !found {
		t.Error("addCmd should be added to rootCmd")
	}
}

// Note: Full integration testing of git add commands would require
// setting up a temporary git repository and mocking exec.Command,
// which is complex for unit tests. Consider integration tests for
// end-to-end validation.
