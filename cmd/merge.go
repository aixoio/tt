package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/aixoio/tt/styles"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:     "merge <source> [target]",
	Aliases: []string{"m"},
	Short:   "Merge branches with intelligent conflict handling",
	Long:    styles.Info.Render("Merge source branch into target branch with intelligent prompts. If branches not provided, will show interactive selection."),
	Args:    cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Show header
		fmt.Println(styles.Header.Render("Git Merge"))
		fmt.Println()
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
		fmt.Println(styles.Card.Render(
			styles.InfoIcon + " " + styles.Info.Render("Merge Operation") + "\n" +
				styles.Neutral.Render("Source: ") + styles.Branch.Render(sourceBranch) + "\n" +
				styles.Neutral.Render("Target: ") + styles.Branch.Render(targetBranch),
		))

		// Check if we're already on the target branch
		if currentBranch != targetBranch {
			fmt.Print(styles.Spinner.Render("⏳") + " " + styles.Info.Render("Switching to target branch ") + styles.Branch.Render(targetBranch) + "... ")
			if err := exec.Command("git", "checkout", targetBranch).Run(); err != nil {
				fmt.Println(styles.ErrorIcon)
				return fmt.Errorf("failed to switch to target branch: %w", err)
			}
			fmt.Println(styles.SuccessIcon)
		}

		// Perform the merge
		fmt.Print(styles.Spinner.Render("⏳") + " " + styles.Info.Render("Merging ") + styles.Branch.Render(sourceBranch) + styles.Info.Render(" into ") + styles.Branch.Render(targetBranch) + "... ")
		gitMergeCmd := exec.Command("git", "merge", sourceBranch)
		gitMergeCmd.Stdout = os.Stdout
		gitMergeCmd.Stderr = os.Stderr

		if err := gitMergeCmd.Run(); err != nil {
			fmt.Println(styles.ErrorIcon)
			return fmt.Errorf("merge failed: %w", err)
		}
		fmt.Println(styles.SuccessIcon)

		// Show success message
		fmt.Println()
		fmt.Println(styles.Card.Render(
			styles.SuccessIcon + " " + styles.Success.Render("Merge completed successfully!") + "\n" +
				styles.Neutral.Render("Merged: ") + styles.Branch.Render(sourceBranch) + " → " + styles.Branch.Render(targetBranch),
		))

		// Handle auto-delete flag
		deleteFlag, _ := cmd.Flags().GetBool("delete")
		if deleteFlag {
			fmt.Println()
			fmt.Print(styles.InfoIcon + " " + styles.Info.Render("Deleting source branch ") + styles.Branch.Render(sourceBranch) + "... ")
			if err := exec.Command("git", "branch", "-d", sourceBranch).Run(); err != nil {
				fmt.Println(styles.ErrorIcon)
				fmt.Println(styles.WarningIcon + " " + styles.Warning.Render("Could not delete branch: ") + styles.Muted.Render(err.Error()))
			} else {
				fmt.Println(styles.SuccessIcon)
				fmt.Println(styles.SuccessIcon + " " + styles.Success.Render("Branch ") + styles.Branch.Render(sourceBranch) + styles.Success.Render(" deleted"))
			}
		}

		// Handle auto-push flag
		pushFlag, _ := cmd.Flags().GetBool("push")
		if pushFlag {
			fmt.Println()
			fmt.Print(styles.InfoIcon + " " + styles.Info.Render("Pushing merged changes... "))
			if err := pushChanges(); err != nil {
				fmt.Println(styles.ErrorIcon)
				return fmt.Errorf("failed to push after merge: %w", err)
			}
			fmt.Println(styles.SuccessIcon)
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
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(styles.Primary.Render(prompt)).
				Description(styles.Neutral.Render("Choose a branch from the list below")).
				Options(options...).
				Value(branch),
		),
	).WithTheme(huh.ThemeCharm())

	return form.Run()
}

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.Flags().BoolP("push", "p", false, "Push after successful merge")
	mergeCmd.Flags().BoolP("delete", "d", false, "Delete source branch after merge")
}
