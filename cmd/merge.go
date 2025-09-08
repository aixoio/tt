package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:     "merge <source> [target]",
	Aliases: []string{"m"},
	Short:   "Merge branches with smart conflict handling",
	Long: `Merge source branch into target branch with intelligent prompts.
If branches not provided, will show interactive selection.`,
	Args: cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceBranch := ""
		targetBranch := ""

		// Get current branch as default target
		currentBranch, err := getCurrentBranch()
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}

		// Handle source branch
		if len(args) > 0 {
			sourceBranch = args[0]
		} else {
			// Interactive selection for source branch
			if err := selectBranch("Select source branch", &sourceBranch); err != nil {
				return fmt.Errorf("failed to select source branch: %w", err)
			}
		}

		if sourceBranch == "" {
			return fmt.Errorf("source branch cannot be empty")
		}

		// Handle target branch
		if len(args) > 1 {
			targetBranch = args[1]
		} else {
			// Use current branch as default if not specified
			targetBranch = currentBranch
		}

		// Confirm merge details
		fmt.Printf("Merging '%s' into '%s'\n", sourceBranch, targetBranch)

		// Check if we're already on the target branch
		if currentBranch != targetBranch {
			fmt.Printf("Switching to target branch '%s'...\n", targetBranch)
			if err := exec.Command("git", "checkout", targetBranch).Run(); err != nil {
				return fmt.Errorf("failed to switch to target branch: %w", err)
			}
		}

		// Perform the merge
		fmt.Printf("Merging '%s' into '%s'...\n", sourceBranch, targetBranch)
		gitMergeCmd := exec.Command("git", "merge", sourceBranch)
		gitMergeCmd.Stdout = os.Stdout
		gitMergeCmd.Stderr = os.Stderr

		if err := gitMergeCmd.Run(); err != nil {
			return fmt.Errorf("merge failed: %w", err)
		}

		fmt.Printf("Successfully merged '%s' into '%s'.\n", sourceBranch, targetBranch)

		// Handle auto-delete flag
		deleteFlag, _ := cmd.Flags().GetBool("delete")
		if deleteFlag {
			fmt.Printf("Deleting source branch '%s'...\n", sourceBranch)
			if err := exec.Command("git", "branch", "-d", sourceBranch).Run(); err != nil {
				fmt.Printf("Warning: failed to delete source branch: %v\n", err)
			} else {
				fmt.Printf("Deleted branch '%s'.\n", sourceBranch)
			}
		}

		// Handle auto-push flag
		pushFlag, _ := cmd.Flags().GetBool("push")
		if pushFlag {
			fmt.Printf("Pushing changes...\n")
			if err := pushChanges(); err != nil {
				return fmt.Errorf("failed to push after merge: %w", err)
			}
		}

		return nil
	},
}

// getCurrentBranch returns the name of the current git branch
func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// selectBranch shows an interactive selection of branches
func selectBranch(prompt string, branch *string) error {
	// Get list of branches
	cmd := exec.Command("git", "branch", "--format", "%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	// Parse branches
	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(branches) == 0 {
		return fmt.Errorf("no branches found")
	}

	// Create options for huh
	options := make([]huh.Option[string], len(branches))
	for i, b := range branches {
		b = strings.TrimSpace(b)
		options[i] = huh.NewOption(b, b)
	}

	// Show selection form
	return huh.NewSelect[string]().
		Title(prompt).
		Options(options...).
		Value(branch).
		Run()
}

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.Flags().BoolP("push", "p", false, "Push after successful merge")
	mergeCmd.Flags().BoolP("delete", "d", false, "Delete source branch after merge")
}
