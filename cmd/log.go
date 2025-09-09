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
		full, _ := cmd.Flags().GetBool("full")
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

		// Get current branch
		branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		branchOutput, err := branchCmd.Output()
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}
		currentBranch := strings.TrimSpace(string(branchOutput))
		fmt.Println(styles.Primary.Render("On branch: ") + styles.Branch.Render(currentBranch))
		fmt.Println()

		// Get upstream commit if exists
		upstreamHash := ""
		upstreamCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
		if upstreamOutput, err := upstreamCmd.Output(); err == nil {
			upstreamBranch := strings.TrimSpace(string(upstreamOutput))
			if upstreamBranch != "" {
				hashCmd := exec.Command("git", "rev-parse", upstreamBranch)
				if hashOutput, err := hashCmd.Output(); err == nil {
					upstreamHash = strings.TrimSpace(string(hashOutput))[:7] // Short hash
				}
			}
		}

		// Build git log command
		gitArgs := []string{"log", "--oneline"}

		// Add reverse flag to show in reverse chronological order
		gitArgs = append(gitArgs, "--reverse")

		// Set default count to 10 if not full
		if !full && count == 10 {
			gitArgs = append(gitArgs, "-n10")
		} else if !full && count > 0 {
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
		hasOutput := false
		for scanner.Scan() {
			hasOutput = true
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

					// Check if this is the upstream commit
					shortHash := strings.Trim(hash, "* ")
					if len(shortHash) > 7 {
						shortHash = shortHash[:7]
					}

					// Style the output
					styledLine := styles.CommitHash.Render(hash) + " " + styles.Primary.Render(message)

					// Add origin marker if this is the upstream commit
					if shortHash == upstreamHash {
						styledLine += " " + styles.Highlight.Render("[origin]")
					}

					fmt.Println(styledLine)
				} else {
					// Fallback for any other format
					fmt.Println(styles.Neutral.Render(line))
				}
			}
		}

		// Show message if no commits
		if !hasOutput {
			fmt.Println(styles.InfoIcon + " " + styles.Info.Render("No commits found"))
		}

		if err := gitCmd.Wait(); err != nil {
			// Git log might exit with error if there are no commits
			return nil
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading git log output: %w", err)
		}

		// Check for uncommitted changes
		statusCmd := exec.Command("git", "status", "--porcelain")
		if statusOutput, err := statusCmd.Output(); err == nil && len(strings.TrimSpace(string(statusOutput))) > 0 {
			fmt.Println(styles.WarningIcon + " " + styles.Warning.Render("Uncommitted changes"))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.Flags().BoolP("full", "f", false, "Show all commits (default: show last 10)")
	logCmd.Flags().IntP("count", "n", 10, "Number of commits to show (ignored with --full)")
	logCmd.Flags().BoolP("all", "a", false, "Show all branches")
	logCmd.Flags().BoolP("graph", "g", false, "Show graph")
}
