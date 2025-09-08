package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:     "push",
	Aliases: []string{"p"},
	Short:   "Push changes to remote",
	Long:    `Push changes to the remote repository. Automatically sets upstream if not configured.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := pushChanges(); err != nil {
			return fmt.Errorf("failed to push: %w", err)
		}
		return nil
	},
}

func pushChanges() error {
	// Check if upstream is set
	upstreamCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err := upstreamCmd.Run(); err != nil {
		// No upstream, set it
		fmt.Println("No upstream branch. Setting upstream to origin/HEAD...")
		pushCmd := exec.Command("git", "push", "--set-upstream", "origin", "HEAD")
		pushCmd.Stdout = os.Stdout
		pushCmd.Stderr = os.Stderr
		if err := pushCmd.Run(); err != nil {
			return fmt.Errorf("failed to push and set upstream: %w", err)
		}
		return nil
	}

	// Upstream exists, just push
	pushCmd := exec.Command("git", "push")
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	if err := pushCmd.Run(); err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
