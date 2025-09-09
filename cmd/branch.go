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

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:     "branch [command] [args]",
	Aliases: []string{"b"},
	Short:   "Create, switch, and list git branches",
	Long:    styles.Info.Render("Manage git branches: create new branches, switch to existing ones, or list all branches with a commit graph."),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return listBranches()
		}

		if len(args) == 2 && args[0] == "delete" {
			return deleteBranch(args[1], cmd)
		}

		if len(args) == 1 {
			return handleSwitchOrCreate(args[0], cmd)
		}

		// Invalid usage
		fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Invalid usage. Use 'tt branch' to list, 'tt branch <name>' to switch or create, or 'tt branch delete <name>' to delete."))
		return fmt.Errorf("invalid command arguments")
	},
}

func checkoutBranch(name string) error {
	// Verify branch exists
	verifyCmd := exec.Command("git", "rev-parse", "--verify", name)
	if err := verifyCmd.Run(); err != nil {
		return fmt.Errorf("branch '%s' does not exist", name)
	}

	// Switch to branch
	checkoutCmd := exec.Command("git", "checkout", name)
	if err := checkoutCmd.Run(); err != nil {
		return fmt.Errorf("failed to switch to branch '%s': %w", name, err)
	}

	fmt.Println(styles.SuccessIcon + " " + styles.Success.Render("Switched to branch: ") + styles.Branch.Render(name))
	return nil
}

// handleSwitchOrCreate switches to branch if exists, else creates it
func handleSwitchOrCreate(name string, cmd *cobra.Command) error {
	// Verify branch exists
	verifyCmd := exec.Command("git", "rev-parse", "--verify", name)
	if err := verifyCmd.Run(); err != nil {
		// Branch does not exist, create it
		return createBranch(name, cmd)
	}

	// Branch exists, switch to it
	return checkoutBranch(name)
}

// listBranches shows all branches and a graph
func listBranches() error {
	fmt.Println(styles.Header.Render("Git Branches"))
	fmt.Println()

	// Get current branch
	currentBranchCmd := exec.Command("git", "branch", "--show-current")
	if currentOutput, err := currentBranchCmd.Output(); err == nil {
		currentBranch := strings.TrimSpace(string(currentOutput))
		fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Current branch: ") + styles.Branch.Render(currentBranch))
		fmt.Println()
	}

	// Get all branches
	fmt.Println(styles.Neutral.Render("Branches:"))
	branchesCmd := exec.Command("git", "branch", "--list")
	if branchesOutput, err := branchesCmd.Output(); err == nil {
		branches := strings.TrimSpace(string(branchesOutput))
		for branch := range strings.SplitSeq(branches, "\n") {
			branch = strings.TrimSpace(branch)
			if branch != "" {
				if strings.HasPrefix(branch, "*") {
					fmt.Println("  " + styles.SuccessIcon + " " + styles.Branch.Render(branch[1:]) + " (current)")
				} else {
					fmt.Println("  • " + styles.Neutral.Render(branch))
				}
			}
		}
	}
	fmt.Println()

	// Show graph
	fmt.Println(styles.Neutral.Render("Recent Commits Graph:"))
	graphCmd := exec.Command("git", "log", "--graph", "--oneline", "--decorate", "--all", "-n", "10")
	if graphOutput, err := graphCmd.Output(); err == nil {
		fmt.Println(styles.Muted.Render(string(graphOutput)))
	}

	return nil
}

// deleteBranch deletes a branch after confirmation
func deleteBranch(name string, cmd *cobra.Command) error {
	remote, _ := cmd.Flags().GetBool("remote")

	if remote {
		// Delete remote branch
		// First, confirm with phrase
		var phrase string
		// Build description string safely
		var descBuilder strings.Builder
		descBuilder.WriteString("Type 'confirm delete remote ")
		descBuilder.WriteString(name)
		descBuilder.WriteString("' to continue")
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title(styles.Warning.Render("Enter confirmation phrase")).
					Description(styles.Neutral.Render(descBuilder.String())).
					Value(&phrase).
					Validate(func(s string) error {
						expected := fmt.Sprintf("confirm delete remote %s", name)
						if s != expected {
							return fmt.Errorf("phrase does not match")
						}
						return nil
					}),
			),
		).WithTheme(huh.ThemeCharm())

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}

		if phrase != fmt.Sprintf("confirm delete remote %s", name) {
			fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Remote branch deletion cancelled."))
			return nil
		}

		// Delete remote
		fmt.Print(styles.Spinner.Render("⏳") + " " + styles.Info.Render("Deleting remote branch... "))
		deleteRemoteCmd := exec.Command("git", "push", "origin", "--delete", name)
		output, err := deleteRemoteCmd.CombinedOutput()
		if err != nil {
			fmt.Println(styles.ErrorIcon)
			return fmt.Errorf("failed to delete remote branch: %v", output)
		}
		fmt.Println(styles.SuccessIcon)
		fmt.Println()
		fmt.Println(styles.SuccessIcon + " " + styles.Success.Render("Remote branch deleted successfully!") + " " + styles.Branch.Render(name))
		return nil
	}

	// Local delete logic as before
	// Verify branch exists
	verifyCmd := exec.Command("git", "rev-parse", "--verify", name)
	if err := verifyCmd.Run(); err != nil {
		return fmt.Errorf("branch '%s' does not exist", name)
	}

	// Get current branch to check if deleting current
	currentBranchCmd := exec.Command("git", "branch", "--show-current")
	var currentBranch string
	if currentOutput, err := currentBranchCmd.Output(); err == nil {
		currentBranch = strings.TrimSpace(string(currentOutput))
	}

	if name == currentBranch {
		return fmt.Errorf("cannot delete the current branch '%s'", name)
	}

	// Try to delete first with -d
	fmt.Print(styles.Spinner.Render("⏳") + " " + styles.Info.Render("Deleting branch... "))
	deleteCmd := exec.Command("git", "branch", "-d", name)
	output, err := deleteCmd.CombinedOutput()
	if err != nil {
		fmt.Println(styles.ErrorIcon)
		outputStr := strings.TrimSpace(string(output))
		if strings.Contains(outputStr, "not fully merged") {
			// Prompt for force delete
			var confirmForce bool
			// Build description string safely
			var forceDescBuilder strings.Builder
			forceDescBuilder.WriteString("Branch '")
			forceDescBuilder.WriteString(name)
			forceDescBuilder.WriteString("' is not fully merged. Force delete anyway?")
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title(styles.Warning.Render("Force delete unmerged branch")).
						Description(styles.Neutral.Render(forceDescBuilder.String())).
						Value(&confirmForce),
				),
			).WithTheme(huh.ThemeCharm())

			if err := form.Run(); err != nil {
				return fmt.Errorf("failed to get confirmation: %w", err)
			}

			if !confirmForce {
				fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Branch deletion cancelled."))
				return nil
			}

			// Force delete with -D
			fmt.Print(styles.Spinner.Render("⏳") + " " + styles.Info.Render("Force deleting branch... "))
			forceDeleteCmd := exec.Command("git", "branch", "-D", name)
			if output, err := forceDeleteCmd.CombinedOutput(); err != nil {
				fmt.Println(styles.ErrorIcon)
				return fmt.Errorf("failed to force delete branch: %v", output)
			}
			fmt.Println(styles.SuccessIcon)
		} else {
			return fmt.Errorf("failed to delete branch: %v", output)
		}
	} else {
		fmt.Println(styles.SuccessIcon)
	}

	// Success
	fmt.Println()
	fmt.Println(styles.SuccessIcon + " " + styles.Success.Render("Branch deleted successfully!") + " " + styles.Branch.Render(name))

	return nil
}

// createBranch creates a new branch with the given name
func createBranch(branchName string, cmd *cobra.Command) error {
	if branchName == "" {
		fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Branch name cannot be empty"))
		return fmt.Errorf("branch name cannot be empty")
	}

	// Show header
	fmt.Println(styles.Header.Render("Git Branch"))
	fmt.Println()

	// Get current branch for context
	currentBranchCmd := exec.Command("git", "branch", "--show-current")
	if currentOutput, err := currentBranchCmd.Output(); err == nil {
		currentBranch := strings.TrimSpace(string(currentOutput))
		fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Current branch: ") + styles.Branch.Render(currentBranch))
		fmt.Println()
	}

	// Show branch creation details
	fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Creating branch: ") + styles.Branch.Render(branchName))

	// Create the new branch
	fmt.Print(styles.Spinner.Render("⏳") + " " + styles.Info.Render("Creating branch... "))
	gitCreateCmd := exec.Command("git", "checkout", "-b", branchName)
	gitCreateCmd.Stdout = os.Stdout
	gitCreateCmd.Stderr = os.Stderr
	if err := gitCreateCmd.Run(); err != nil {
		fmt.Println(styles.ErrorIcon)
		return fmt.Errorf("failed to create branch: %w", err)
	}
	fmt.Println(styles.SuccessIcon)

	// Show success message
	fmt.Println()
	fmt.Println(styles.Card.Render(
		styles.SuccessIcon + " " + styles.Success.Render("Branch created successfully!") + "\n" +
			styles.Neutral.Render("Branch: ") + styles.Branch.Render(branchName) + "\n" +
			styles.Neutral.Render("Status: ") + styles.Success.Render("Switched to new branch"),
	))

	pushFlag, _ := cmd.Flags().GetBool("push")
	if pushFlag {
		fmt.Println()
		fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Pushing branch to remote... "))
		if err := pushChangesToNewBranch(branchName); err != nil {
			fmt.Println(styles.ErrorIcon)
			fmt.Println(styles.WarningIcon + " " + styles.Warning.Render("Push failed, but branch was created locally"))
			return fmt.Errorf("failed to push branch: %w", err)
		}
		fmt.Println(styles.SuccessIcon)
	}

	return nil
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
	branchCmd.Flags().Bool("remote", false, "Delete remote branch instead of local")
}
