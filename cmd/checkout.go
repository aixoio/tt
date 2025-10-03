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

var checkoutCmd = &cobra.Command{
	Use:     "checkout",
	Aliases: []string{"co"},
	Short:   "Interactively checkout branches or commits",
	Long: styles.Info.Render("Checkout a branch or specific commit through an interactive menu. " +
		"Select from local branches or search through commits to checkout."),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Show header
		fmt.Println(styles.Header.Render("Git Checkout"))
		fmt.Println()

		// Check if we're in a git repository
		if _, err := exec.Command("git", "rev-parse", "--git-dir").Output(); err != nil {
			fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Not a git repository"))
			return fmt.Errorf("not a git repository")
		}

		// Check for uncommitted changes (warning only, not blocking)
		statusCmd := exec.Command("git", "status", "--porcelain")
		if statusOutput, err := statusCmd.Output(); err == nil && len(strings.TrimSpace(string(statusOutput))) > 0 {
			fmt.Println(styles.WarningIcon + " " + styles.Warning.Render("You have uncommitted changes."))
			fmt.Println(styles.Neutral.Render("  Git will carry them over if they don't conflict, or block the checkout if they would be overwritten."))
			fmt.Println()
		}

		// Show current branch/state
		currentBranchCmd := exec.Command("git", "branch", "--show-current")
		if currentOutput, err := currentBranchCmd.Output(); err == nil {
			currentBranch := strings.TrimSpace(string(currentOutput))
			if currentBranch != "" {
				fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Current branch: ") + styles.Branch.Render(currentBranch))
			} else {
				// Detached HEAD state
				headCmd := exec.Command("git", "rev-parse", "--short", "HEAD")
				if headOutput, err := headCmd.Output(); err == nil {
					headHash := strings.TrimSpace(string(headOutput))
					fmt.Println(styles.WarningIcon + " " + styles.Warning.Render("Currently in detached HEAD state at: ") + styles.CommitHash.Render(headHash))
				}
			}
			fmt.Println()
		}

		// Determine checkout type
		checkoutType, err := selectCheckoutType()
		if err != nil {
			return err
		}

		if checkoutType == "" {
			fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Checkout cancelled"))
			return nil
		}

		var target string

		switch checkoutType {
		case "branch":
			target, err = selectBranchForCheckout()
			if err != nil {
				return err
			}
			if target == "" {
				fmt.Println(styles.InfoIcon + " " + styles.Info.Render("No branch selected"))
				return nil
			}
		case "commit":
			target, err = selectCommit()
			if err != nil {
				return err
			}
			if target == "" {
				fmt.Println(styles.InfoIcon + " " + styles.Info.Render("No commit selected"))
				return nil
			}

			// Show detached HEAD warning for commits
			fmt.Println()
			fmt.Println(styles.Card.Render(
				styles.WarningIcon + " " + styles.Warning.Render("Detached HEAD Warning") + "\n\n" +
					styles.Neutral.Render("Checking out a commit will put you in 'detached HEAD' state.\n") +
					styles.Neutral.Render("You can look around, make experimental changes and commit them,\n") +
					styles.Neutral.Render("but any commits you make will be lost when you checkout a branch.\n\n") +
					styles.Info.Render("To keep your changes, create a new branch with: tt branch <name>"),
			))
			fmt.Println()

			// Confirm detached HEAD checkout
			var confirm bool
			prompt := huh.NewConfirm().
				Title(styles.WarningIcon + " " + styles.Warning.Render("Continue with Detached HEAD Checkout?")).
				Description("This will detach your HEAD from any branch.").
				Value(&confirm).
				Affirmative("Yes, checkout commit").
				Negative("No, cancel").
				WithTheme(huh.ThemeCharm())

			if err := prompt.Run(); err != nil {
				return fmt.Errorf("failed to show confirmation prompt: %w", err)
			}

			if !confirm {
				fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Checkout cancelled"))
				return nil
			}
		default:
			return fmt.Errorf("invalid checkout type")
		}

		// Perform the checkout
		return performCheckout(target, checkoutType)
	},
}

func selectCheckoutType() (string, error) {
	var checkoutType string
	var options = []huh.Option[string]{
		huh.NewOption("Checkout a branch", "branch"),
		huh.NewOption("Checkout a specific commit (detached HEAD)", "commit"),
	}

	selectPrompt := huh.NewSelect[string]().
		Title(styles.Primary.Render("What would you like to checkout?")).
		Options(options...).
		Value(&checkoutType).
		WithTheme(huh.ThemeCharm())

	if err := selectPrompt.Run(); err != nil {
		return "", fmt.Errorf("failed to get checkout type: %w", err)
	}

	return checkoutType, nil
}

func selectBranchForCheckout() (string, error) {
	// Get all local branches
	cmd := exec.Command("git", "branch", "--list")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get branches: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("no branches found")
	}

	// Get current branch
	currentBranchCmd := exec.Command("git", "branch", "--show-current")
	var currentBranch string
	if currentOutput, err := currentBranchCmd.Output(); err == nil {
		currentBranch = strings.TrimSpace(string(currentOutput))
	}

	var options []huh.Option[string]
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove the * prefix if present
		branchName := strings.TrimPrefix(line, "* ")
		branchName = strings.TrimSpace(branchName)

		// Get last commit message for this branch
		commitCmd := exec.Command("git", "log", "-1", "--pretty=format:%s", branchName)
		var lastCommit string
		if commitOutput, err := commitCmd.Output(); err == nil {
			lastCommit = strings.TrimSpace(string(commitOutput))
			if len(lastCommit) > 60 {
				lastCommit = lastCommit[:60] + "..."
			}
		}

		// Format display string
		var display string
		if branchName == currentBranch {
			display = fmt.Sprintf("%s %s %s",
				styles.SuccessIcon,
				styles.Branch.Render(branchName),
				styles.Muted.Render("(current)"))
		} else {
			display = fmt.Sprintf("%s %s",
				styles.Branch.Render(branchName),
				styles.Muted.Render("- "+lastCommit))
		}

		options = append(options, huh.NewOption(display, branchName))
	}

	if len(options) == 0 {
		return "", fmt.Errorf("no branches available")
	}

	var selectedBranch string
	selectPrompt := huh.NewSelect[string]().
		Title(styles.Primary.Render("Select a branch to checkout:")).
		Options(options...).
		Value(&selectedBranch).
		WithTheme(huh.ThemeCharm())

	if err := selectPrompt.Run(); err != nil {
		return "", fmt.Errorf("failed to select branch: %w", err)
	}

	return selectedBranch, nil
}

func selectCommit() (string, error) {
	var selectionType string
	var options = []huh.Option[string]{
		huh.NewOption("Recent commits (last 20)", "recent"),
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
		return selectFromRecentCommitsForCheckout()
	case "search":
		return searchAllCommitsForCheckout()
	default:
		return "", fmt.Errorf("invalid selection")
	}
}

func selectFromRecentCommitsForCheckout() (string, error) {
	// Get last 20 commits
	cmd := exec.Command("git", "log", "--oneline", "-n", "20")
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

			// Get author and date for additional context
			detailCmd := exec.Command("git", "show", "--no-patch", "--format=%an, %ar", hash)
			var detail string
			if detailOutput, err := detailCmd.Output(); err == nil {
				detail = strings.TrimSpace(string(detailOutput))
			}

			display := fmt.Sprintf("%s %s %s",
				styles.CommitHash.Render(hash),
				styles.Primary.Render(message),
				styles.Muted.Render("("+detail+")"))
			options = append(options, huh.NewOption(display, hash))
		}
	}

	var selectedHash string
	selectPrompt := huh.NewSelect[string]().
		Title(styles.Primary.Render("Select a commit to checkout:")).
		Options(options...).
		Value(&selectedHash).
		WithTheme(huh.ThemeCharm())

	if err := selectPrompt.Run(); err != nil {
		return "", fmt.Errorf("failed to select commit: %w", err)
	}

	return selectedHash, nil
}

func searchAllCommitsForCheckout() (string, error) {
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

			// Get author and date for additional context
			detailCmd := exec.Command("git", "show", "--no-patch", "--format=%an, %ar", hash)
			var detail string
			if detailOutput, err := detailCmd.Output(); err == nil {
				detail = strings.TrimSpace(string(detailOutput))
			}

			display := fmt.Sprintf("%s %s %s",
				styles.CommitHash.Render(hash),
				styles.Primary.Render(message),
				styles.Muted.Render("("+detail+")"))
			options = append(options, huh.NewOption(display, hash))
		}
	}

	var selectedHash string
	selectPrompt := huh.NewSelect[string]().
		Title(styles.Primary.Render("Search and select a commit to checkout:")).
		Options(options...).
		Value(&selectedHash).
		WithTheme(huh.ThemeCharm())

	if err := selectPrompt.Run(); err != nil {
		return "", fmt.Errorf("failed to select commit: %w", err)
	}

	return selectedHash, nil
}

func performCheckout(target string, checkoutType string) error {
	// Validate target exists
	if checkoutType == "commit" {
		// Validate commit hash format
		if len(target) < 7 {
			return fmt.Errorf("commit hash must be at least 7 characters long")
		}
		if matched, _ := regexp.MatchString("^[0-9a-fA-F]+$", target); !matched {
			return fmt.Errorf("invalid commit hash format")
		}

		// Check if commit exists
		cmd := exec.Command("git", "cat-file", "-t", target)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("commit %s does not exist", target)
		}
	} else {
		// Validate branch exists
		verifyCmd := exec.Command("git", "rev-parse", "--verify", target)
		if err := verifyCmd.Run(); err != nil {
			return fmt.Errorf("branch '%s' does not exist", target)
		}
	}

	// Perform checkout
	fmt.Print(styles.SpinnerIcon + " " + styles.Info.Render("Checking out... "))

	checkoutCmd := exec.Command("git", "checkout", target)
	checkoutCmd.Stdout = os.Stdout
	checkoutCmd.Stderr = os.Stderr

	if err := checkoutCmd.Run(); err != nil {
		fmt.Println(styles.ErrorIcon)
		return fmt.Errorf("failed to checkout: %w", err)
	}

	fmt.Println(styles.SuccessIcon)
	fmt.Println()

	// Show success message
	if checkoutType == "branch" {
		fmt.Println(styles.Card.Render(
			styles.Success.Render("Checkout successful!") + "\n" +
				styles.Neutral.Render("Switched to branch: ") + styles.Branch.Render(target),
		))
	} else {
		fmt.Println(styles.Card.Render(
			styles.Success.Render("Checkout successful!") + "\n" +
				styles.Neutral.Render("HEAD is now at: ") + styles.CommitHash.Render(target) + "\n" +
				styles.Warning.Render("You are in 'detached HEAD' state"),
		))
	}

	return nil
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
}
