package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:     "c [message]",
	Aliases: []string{"commit"},
	Short:   "Commit changes",
	Long:    `Commit changes to git. Can use interactive prompt for message.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		message, _ := cmd.Flags().GetString("message")
		addFlag, _ := cmd.Flags().GetBool("add")

		if addFlag {
			if err := exec.Command("git", "add", ".").Run(); err != nil {
				return fmt.Errorf("failed to add files: %w", err)
			}
		}

		if message == "" {
			if err := huh.NewInput().
				Title("Commit message").
				Placeholder("Enter your commit message...").
				Value(&message).
				Run(); err != nil {
				return fmt.Errorf("failed to get commit message: %w", err)
			}
		}

		if message == "" {
			return fmt.Errorf("commit message cannot be empty")
		}

		gitCmd := exec.Command("git", "commit", "-m", message)
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr

		if err := gitCmd.Run(); err != nil {
			return fmt.Errorf("failed to commit: %w", err)
		}

		pushFlag, _ := cmd.Flags().GetBool("push")
		if pushFlag {
			if err := pushChanges(); err != nil {
				return fmt.Errorf("failed to push after commit: %w", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringP("message", "m", "", "Commit message")
	commitCmd.Flags().BoolP("add", "a", false, "Add all files before committing")
	commitCmd.Flags().BoolP("push", "p", false, "Push after committing")
}
