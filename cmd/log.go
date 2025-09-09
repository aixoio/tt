package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aixoio/tt/styles"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:     "log [options]",
	Aliases: []string{"l"},
	Short:   "Show commit history with beautiful formatting",
	Long:    "Display git commit history in a styled, readable format similar to git log --oneline",
	RunE: func(cmd *cobra.Command, args []string) error {
		count, _ := cmd.Flags().GetInt("count")
		all, _ := cmd.Flags().GetBool("all")
		graph, _ := cmd.Flags().GetBool("graph")

		// Show header
		fmt.Println(styles.Header.Render("Git Log"))
		fmt.Println()

		// Check if we're in a git repository
		if _, err := exec.Command("git", "rev-parse", "--git-dir").Output(); err != nil {
			fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Not a git repository"))
			return fmt.Errorf("not a git repository")
		}

		// Build git log command
		gitArgs := []string{"log", "--oneline"}

		if count > 0 {
			gitArgs = append(gitArgs, fmt.Sprintf("-n%d", count))
		}

		if all {
			gitArgs = append(gitArgs, "--all")
		}

		if graph {
			gitArgs = append(gitArgs, "--graph")
		}

		// Execute git log command
		gitCmd := exec.Command("git", gitArgs...)
		stdout, err := gitCmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to create stdout pipe: %w", err)
		}

		if err := gitCmd.Start(); err != nil {
			return fmt.Errorf("failed to start git log: %w", err)
		}

		// Process output line by line
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()

			if graph {
				// For graph mode, we need special handling
				fmt.Println(styles.Neutral.Render(line))
			} else {
				// Parse commit hash and message
				parts := strings.SplitN(line, " ", 2)
				if len(parts) >= 2 {
					hash := parts[0]
					message := parts[1]

					// Style the output
					styledLine := styles.CommitHash.Render(hash) + " " + styles.Primary.Render(message)
					fmt.Println(styledLine)
				} else {
					// Fallback for any other format
					fmt.Println(styles.Neutral.Render(line))
				}
			}
		}

		if err := gitCmd.Wait(); err != nil {
			// Git log might exit with error if there are no commits
			return nil
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading git log output: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.Flags().IntP("count", "n", 0, "Limit the number of commits to show")
	logCmd.Flags().BoolP("all", "a", false, "Show all branches")
	logCmd.Flags().BoolP("graph", "g", false, "Show graph")
}
