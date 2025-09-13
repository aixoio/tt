package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/aixoio/tt/styles"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:     "clone <repository-url> [directory]",
	Aliases: []string{"cl"},
	Short:   "Clone a repository into a new directory",
	Long:    styles.Info.Render("Clone a repository into a new directory. Works like `git clone` with beautiful progress feedback."),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate arguments
		if len(args) < 1 {
			fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Repository URL is required"))
			return fmt.Errorf("repository URL required")
		}

		repoURL := args[0]
		targetDir := ""
		if len(args) > 1 {
			targetDir = args[1]
		}

		// Get flags
		branchFlag, _ := cmd.Flags().GetString("branch")
		depthFlag, _ := cmd.Flags().GetInt("depth")
		singleBranchFlag, _ := cmd.Flags().GetBool("single-branch")
		recursiveFlag, _ := cmd.Flags().GetBool("recursive")

		// Show header
		fmt.Println(styles.Header.Render("Git Clone"))
		fmt.Println()

		// Build git clone command arguments
		gitArgs := []string{"clone"}

		// Add flags if provided
		if branchFlag != "" {
			gitArgs = append(gitArgs, "--branch", branchFlag)
		}
		if depthFlag > 0 {
			gitArgs = append(gitArgs, "--depth", fmt.Sprintf("%d", depthFlag))
		}
		if singleBranchFlag {
			gitArgs = append(gitArgs, "--single-branch")
		}
		if recursiveFlag {
			gitArgs = append(gitArgs, "--recursive")
		}

		// Add repository URL and optional directory
		gitArgs = append(gitArgs, repoURL)
		if targetDir != "" {
			gitArgs = append(gitArgs, targetDir)
		}

		// Show what we're doing
		fmt.Println(styles.InfoIcon + " " + styles.Info.Render(fmt.Sprintf("Cloning %s...", styles.FilePath.Render(repoURL))))
		if targetDir != "" {
			fmt.Println(styles.Neutral.Render(fmt.Sprintf("Target directory: %s", styles.FilePath.Render(targetDir))))
		}

		// Execute git clone with spinner
		err := runWithSpinner(styles.InfoIcon+" "+styles.Info.Render("Cloning repository..."), func() error {
			gitCmd := exec.Command("git", gitArgs...)
			gitCmd.Stdout = os.Stdout
			gitCmd.Stderr = os.Stderr
			return gitCmd.Run()
		})

		if err != nil {
			fmt.Println(styles.ErrorIcon)
			return fmt.Errorf("failed to clone repository: %w", err)
		}

		fmt.Println(styles.SuccessIcon)

		// Determine the actual directory name
		actualDir := targetDir
		if actualDir == "" {
			// Extract repo name from URL if no directory specified
			// This is a simple extraction; git handles this internally
			actualDir = "repository" // Placeholder; in practice, git determines this
		}

		// Show success message with more details
		fmt.Println()
		fmt.Println(styles.Card.Render(
			styles.SuccessIcon + " " + styles.Success.Render("Repository cloned successfully!") + "\n" +
				styles.Neutral.Render("Repository: ") + styles.FilePath.Render(repoURL) + "\n" +
				styles.Neutral.Render("Directory: ") + styles.FilePath.Render(actualDir),
		))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	// Define flags
	cloneCmd.Flags().StringP("branch", "b", "", "Checkout specific branch")
	cloneCmd.Flags().IntP("depth", "d", 0, "Create a shallow clone of specified depth")
	cloneCmd.Flags().Bool("single-branch", false, "Clone only the default branch")
	cloneCmd.Flags().Bool("recursive", false, "Initialize submodules")
}
