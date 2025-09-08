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
	Use:     "branch [name]",
	Aliases: []string{"b"},
	Short:   "Create and manage git branches",
	Long:    styles.Info.Render("Create a new git branch based on the current one. Use -p or --push to auto-push and set upstream."),
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branchName := ""
		if len(args) > 0 {
			branchName = args[0]
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

		if branchName == "" {
			// Create styled input form
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title(styles.Primary.Render("New Branch Name")).
						Placeholder(styles.Muted.Render("feature/my-awesome-feature")).
						Description(styles.Neutral.Render("Choose a descriptive branch name")).
						Value(&branchName).
						Validate(func(s string) error {
							if len(s) < 2 {
								return fmt.Errorf("branch name too short")
							}
							if strings.Contains(s, " ") {
								return fmt.Errorf("branch name cannot contain spaces")
							}
							return nil
						}),
				),
			).WithTheme(huh.ThemeCharm())

			if err := form.Run(); err != nil {
				return fmt.Errorf("failed to get branch name: %w", err)
			}
		}

		if branchName == "" {
			fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Branch name cannot be empty"))
			return fmt.Errorf("branch name cannot be empty")
		}

		// Show branch creation details
		fmt.Println()
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
			fmt.Print(styles.InfoIcon + " " + styles.Info.Render("Pushing branch to remote... "))
			if err := pushChangesToNewBranch(branchName); err != nil {
				fmt.Println(styles.ErrorIcon)
				fmt.Println(styles.WarningIcon + " " + styles.Warning.Render("Push failed, but branch was created locally"))
				return fmt.Errorf("failed to push branch: %w", err)
			}
			fmt.Println(styles.SuccessIcon)
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
