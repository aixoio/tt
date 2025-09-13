package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/aixoio/tt/styles"
)

type CommitInfo struct {
	Hash    string
	Message string
	Author  string
	Date    string
}

var revertCmd = &cobra.Command{
	Use:     "revert [commit-hash]",
	Aliases: []string{"rv"},
	Short:   "Revert a commit by creating a new commit that undoes the changes",
	Long: styles.Info.Render("Revert a specific commit by creating a new commit that undoes the changes. " +
		"You can specify a commit hash directly, select from recent commits, or search through all commits."),
	RunE: func(cmd *cobra.Command, args []string) error {
		var targetCommit string
		var err error

		// Show header
		fmt.Println(styles.Header.Render("Git Revert"))
		fmt.Println()

		// Check if we're in a git repository
		if _, err := exec.Command("git", "rev-parse", "--git-dir").Output(); err != nil {
			fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Not a git repository"))
			return fmt.Errorf("not a git repository")
		}

		// Check for uncommitted changes
		statusCmd := exec.Command("git", "status", "--porcelain")
		if statusOutput, err := statusCmd.Output(); err == nil && len(strings.TrimSpace(string(statusOutput))) > 0 {
			fmt.Println(styles.WarningIcon + " " + styles.Warning.Render("You have uncommitted changes. Consider committing or stashing them first."))
			fmt.Println()
		}

		// Determine how to get the commit hash
		if len(args) > 0 {
			// Direct commit hash provided
			targetCommit = args[0]
			if err := validateCommitHash(targetCommit); err != nil {
				return err
			}
		} else {
			// Interactive selection
			targetCommit, err = selectCommitInteractively()
			if err != nil {
				return err
			}
		}

		if targetCommit == "" {
			fmt.Println(styles.InfoIcon + " " + styles.Info.Render("No commit selected"))
			return nil
		}

		// Show commit details before reverting
		commitDetails, err := getCommitDetails(targetCommit)
		if err != nil {
			return err
		}

		// Get diff stats for the commit
		diffStats, err := getRevertDiffStats(targetCommit)
		if err != nil {
			return err
		}

		fmt.Println(styles.Card.Render(
			styles.Info.Render("Commit to revert:") + "\n" +
				styles.CommitHash.Render(commitDetails.Hash) + " " + styles.Primary.Render(commitDetails.Message) + "\n" +
				styles.Neutral.Render("Author: ") + styles.Highlight.Render(commitDetails.Author) + "\n" +
				styles.Neutral.Render("Date: ") + styles.Muted.Render(commitDetails.Date) + "\n\n" +
				styles.Info.Render("Files that will be reverted:") + "\n" +
				diffStats,
		))
		fmt.Println()

		// Confirm revert
		var confirm bool
		prompt := huh.NewConfirm().
			Title(styles.WarningIcon + " " + styles.Warning.Render("Confirm Revert")).
			Description("This will create a new commit that undoes the changes from the selected commit. Continue?").
			Value(&confirm).
			Affirmative("Yes, revert").
			Negative("No, cancel").
			WithTheme(huh.ThemeCharm())

		if err := prompt.Run(); err != nil {
			return fmt.Errorf("failed to show confirmation prompt: %w", err)
		}

		if !confirm {
			fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Revert cancelled"))
			return nil
		}

		// Perform the revert
		return performRevert(targetCommit)
	},
}

func validateCommitHash(hash string) error {
	// Validate hash format (7+ characters, hex)
	if len(hash) < 7 {
		return fmt.Errorf("commit hash must be at least 7 characters long")
	}
	if matched, _ := regexp.MatchString("^[0-9a-fA-F]+$", hash); !matched {
		return fmt.Errorf("invalid commit hash format")
	}

	// Check if commit exists
	cmd := exec.Command("git", "cat-file", "-t", hash)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("commit %s does not exist", hash)
	}

	return nil
}

func selectCommitInteractively() (string, error) {
	var selectionType string
	var options = []huh.Option[string]{
		huh.NewOption("Recent commits (last 5)", "recent"),
		huh.NewOption("Search all commits", "search"),
	}

	selectPrompt := huh.NewSelect[string]().
		Title(styles.Primary.Render("How would you like to select a commit?")).
		Options(options...).
		Value(&selectionType).
		WithTheme(huh.ThemeCharm())

	if err := selectPrompt.Run(); err != nil {
		return "", fmt.Errorf("failed to get selection type: %w", err)
	}

	switch selectionType {
	case "recent":
		return selectFromRecentCommits()
	case "search":
		return searchAllCommits()
	default:
		return "", fmt.Errorf("invalid selection")
	}
}

func selectFromRecentCommits() (string, error) {
	// Get last 5 commits
	cmd := exec.Command("git", "log", "--oneline", "-n", "5")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get recent commits: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("no commits found")
	}

	var options []huh.Option[string]
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) >= 2 {
			hash := parts[0]
			message := parts[1]
			display := fmt.Sprintf("%s %s", styles.CommitHash.Render(hash), styles.Primary.Render(message))
			options = append(options, huh.NewOption(display, hash))
		}
	}

	var selectedHash string
	selectPrompt := huh.NewSelect[string]().
		Title(styles.Primary.Render("Select a commit to revert:")).
		Options(options...).
		Value(&selectedHash).
		WithTheme(huh.ThemeCharm())

	if err := selectPrompt.Run(); err != nil {
		return "", fmt.Errorf("failed to select commit: %w", err)
	}

	return selectedHash, nil
}

func searchAllCommits() (string, error) {
	// Get all commits
	cmd := exec.Command("git", "log", "--oneline", "--all")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get all commits: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("no commits found")
	}

	var options []huh.Option[string]
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) >= 2 {
			hash := parts[0]
			message := parts[1]
			display := fmt.Sprintf("%s %s", styles.CommitHash.Render(hash), styles.Primary.Render(message))
			options = append(options, huh.NewOption(display, hash))
		}
	}

	var selectedHash string
	selectPrompt := huh.NewSelect[string]().
		Title(styles.Primary.Render("Search and select a commit to revert:")).
		Options(options...).
		Value(&selectedHash).
		WithTheme(huh.ThemeCharm())

	if err := selectPrompt.Run(); err != nil {
		return "", fmt.Errorf("failed to select commit: %w", err)
	}

	return selectedHash, nil
}

func getCommitDetails(hash string) (*CommitInfo, error) {
	// Get commit details
	cmd := exec.Command("git", "show", "--no-patch", "--format=%H%n%s%n%an%n%ad", "--date=short", hash)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit details: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 4 {
		return nil, fmt.Errorf("invalid commit details format")
	}

	return &CommitInfo{
		Hash:    lines[0],
		Message: lines[1],
		Author:  lines[2],
		Date:    lines[3],
	}, nil
}

func getRevertDiffStats(hash string) (string, error) {
	// Get diff stats for the commit that will be reverted
	cmd := exec.Command("git", "show", "--stat", "--format=", hash)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff stats: %w", err)
	}

	// Parse and format the diff stats
	stats := strings.TrimSpace(string(output))
	if stats == "" {
		return "No file changes in this commit", nil
	}

	// Format the stats with styling
	lines := strings.Split(stats, "\n")
	var formattedStats []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if this is a file line (contains | and numbers)
		if strings.Contains(line, "|") && (strings.Contains(line, "+") || strings.Contains(line, "-")) {
			// This is a file stat line like "file.txt | 5 +--"
			parts := strings.Split(line, "|")
			if len(parts) >= 2 {
				fileName := strings.TrimSpace(parts[0])
				changeStats := strings.TrimSpace(parts[1])

				// Style the file name
				styledFile := styles.FilePath.Render(fileName)

				// Style the change stats (additions in green, deletions in red)
				changeStats = strings.ReplaceAll(changeStats, "+", styles.Success.Render("+"))
				changeStats = strings.ReplaceAll(changeStats, "-", styles.Error.Render("-"))

				formattedStats = append(formattedStats, fmt.Sprintf("  %s | %s", styledFile, changeStats))
			}
		} else if strings.Contains(line, "file") && strings.Contains(line, "changed") {
			// This is the summary line like "2 files changed, 15 insertions(+), 8 deletions(-)"
			summaryLine := strings.ReplaceAll(line, "insertions(+)", styles.Success.Render("insertions(+)"))
			summaryLine = strings.ReplaceAll(summaryLine, "deletions(-)", styles.Error.Render("deletions(-)"))
			formattedStats = append(formattedStats, styles.Highlight.Render(summaryLine))
		} else {
			// Regular line
			formattedStats = append(formattedStats, styles.Neutral.Render(line))
		}
	}

	return strings.Join(formattedStats, "\n"), nil
}

func performRevert(hash string) error {
	fmt.Print(styles.SpinnerIcon + " " + styles.Info.Render("Preparing revert... "))

	// First, try to revert without committing
	revertCmd := exec.Command("git", "revert", "--no-commit", hash)
	revertCmd.Stdout = os.Stdout
	revertCmd.Stderr = os.Stderr

	if err := revertCmd.Run(); err != nil {
		fmt.Println(styles.ErrorIcon)
		return fmt.Errorf("failed to revert commit: %w", err)
	}

	fmt.Println(styles.SuccessIcon)

	// Check for conflicts
	statusCmd := exec.Command("git", "status", "--porcelain")
	statusOutput, err := statusCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check status: %w", err)
	}

	if strings.Contains(string(statusOutput), "UU ") {
		fmt.Println(styles.WarningIcon + " " + styles.Warning.Render("Merge conflicts detected. Please resolve them and then run 'git commit' to complete the revert."))
		return nil
	}

	// No conflicts, create the revert commit
	fmt.Print(styles.SpinnerIcon + " " + styles.Info.Render("Creating revert commit... "))

	// Get commit details for the revert message
	commitDetails, err := getCommitDetails(hash)
	if err != nil {
		fmt.Println(styles.ErrorIcon)
		return fmt.Errorf("failed to get commit details for revert message: %w", err)
	}

	// Create revert commit message in format: "revert [hash]: [original message]"
	revertMessage := fmt.Sprintf("revert %s: %s", hash, commitDetails.Message)
	commitCmd := exec.Command("git", "commit", "-m", revertMessage)
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr

	if err := commitCmd.Run(); err != nil {
		fmt.Println(styles.ErrorIcon)
		return fmt.Errorf("failed to create revert commit: %w", err)
	}

	fmt.Println(styles.SuccessIcon)

	fmt.Println(styles.Card.Render(
		styles.Success.Render("Revert successful!") + "\n" +
			styles.Neutral.Render("Created commit: ") + styles.Highlight.Render(revertMessage),
	))

	return nil
}

func init() {
	rootCmd.AddCommand(revertCmd)
}
