package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:     "branch [name]",
	Aliases: []string{"b"},
	Short:   "Create a new git branch based on the current one",
	Long:    `Create a new git branch based on the current one. Use -p or --push to auto-push and set upstream.`,
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branchName := ""
		if len(args) > 0 {
			branchName = args[0]
		}

		if branchName == "" {
			if err := huh.NewInput().
				Title("Branch name").
				Placeholder("Enter new branch name...").
				Value(&branchName).
				Run(); err != nil {
				return fmt.Errorf("failed to get branch name: %w", err)
			}
		}

		if branchName == "" {
			return fmt.Errorf("branch name cannot be empty")
		}

		// Create the new branch
		gitCreateCmd := exec.Command("git", "checkout", "-b", branchName)
		gitCreateCmd.Stdout = os.Stdout
		gitCreateCmd.Stderr = os.Stderr
		if err := gitCreateCmd.Run(); err != nil {
			return fmt.Errorf("failed to create branch: %w", err)
		}

		fmt.Printf("Created branch '%s'.\n", branchName)

		pushFlag, _ := cmd.Flags().GetBool("push")
		if pushFlag {
			if err := pushChangesToNewBranch(branchName); err != nil {
				return fmt.Errorf("failed to push branch: %w", err)
			}
		}

		return nil
	},
}

func pushChangesToNewBranch(branchName string) error {
	pushCmd := exec.Command("git", "push", "--set-upstream", "origin", branchName)
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	if err := pushCmd.Run(); err != nil {
		return fmt.Errorf("failed to push and set upstream: %w", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(branchCmd)
	branchCmd.Flags().BoolP("push", "p", false, "Auto-push the new branch and set upstream")
}
