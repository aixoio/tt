package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aixoio/tt/styles"
)

var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"s"},
	Short:   "Show git repository status",
	Long:    styles.Info.Render("Display the current state of the git repository, including staged, unstaged, and untracked files."),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(styles.Header.Render("Git Status"))
		fmt.Println()

		// Get current branch
		branchCmd := exec.Command("git", "branch", "--show-current")
		if branchOutput, err := branchCmd.Output(); err == nil {
			currentBranch := strings.TrimSpace(string(branchOutput))
			fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Current branch: ") + styles.Branch.Render(currentBranch))
			fmt.Println()
		}

		// Get git status
		statusCmd := exec.Command("git", "status", "--porcelain")
		output, err := statusCmd.Output()
		if err != nil {
			fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Failed to get git status"))
			return fmt.Errorf("failed to get git status: %w", err)
		}

		status := string(output)
		if status == "" {
			fmt.Println(styles.SuccessIcon + " " + styles.Success.Render("Working tree clean"))
			fmt.Println(styles.Neutral.Render("No changes to commit."))
			return nil
		}

		// Parse and display status
		lines := strings.Split(strings.TrimSpace(status), "\n")
		var staged []string
		var unstaged []string
		var untracked []string

		for _, line := range lines {
			if len(line) < 3 {
				continue
			}
			// Git status --porcelain format: XY filename
			// where X is staged status, Y is unstaged status (or space)
			// Filename starts after first 2 characters and a space
			statusCode := line[:2]
			fileName := strings.TrimLeft(line[2:], " ") // Skip status code and trim leading spaces

			// Determine which category this file belongs to
			stagedChar := statusCode[0]
			unstagedChar := statusCode[1]

			// Check staged changes (first character)
			if stagedChar != ' ' && stagedChar != '?' {
				staged = append(staged, fileName)
			}

			// Check unstaged changes (second character)
			if unstagedChar == 'M' || unstagedChar == 'D' {
				// Only add to unstaged if not already in staged
				// (to avoid duplicates for files with both staged and unstaged changes)
				alreadyStaged := false
				for _, s := range staged {
					if s == fileName {
						alreadyStaged = true
						break
					}
				}
				if !alreadyStaged {
					unstaged = append(unstaged, fileName)
				}
			}

			// Check untracked files
			if statusCode == "??" {
				untracked = append(untracked, fileName)
			}
		}

		// Display staged files
		if len(staged) > 0 {
			fmt.Println(styles.SuccessIcon + " " + styles.Success.Render("Changes to be committed:"))
			for _, file := range staged {
				fmt.Println("  " + styles.FilePath.Render(file))
			}
			fmt.Println()
		}

		// Display unstaged files
		if len(unstaged) > 0 {
			fmt.Println(styles.WarningIcon + " " + styles.Warning.Render("Changes not staged for commit:"))
			for _, file := range unstaged {
				fmt.Println("  " + styles.FilePath.Render(file))
			}
			fmt.Println()
		}

		// Display untracked files
		if len(untracked) > 0 {
			fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Untracked files:"))
			for _, file := range untracked {
				fmt.Println("  " + styles.FilePath.Render(file))
			}
			fmt.Println()
		}

		// Show summary
		totalChanges := len(staged) + len(unstaged) + len(untracked)
		if totalChanges > 0 {
			fmt.Println(styles.Neutral.Render(fmt.Sprintf("Total: %d file(s) changed", totalChanges)))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
