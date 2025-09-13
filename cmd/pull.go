package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/aixoio/tt/styles"
)

var pullCmd = &cobra.Command{
	Use:     "pull",
	Aliases: []string{"pl"},
	Short:   "Pull changes from remote repository",
	Long:    styles.Info.Render("Pull changes from the remote repository. Automatically sets upstream if not configured."),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Show header
		fmt.Println(styles.Header.Render("Git Pull"))
		fmt.Println()

		if err := pullChanges(); err != nil {
			return fmt.Errorf("failed to pull: %w", err)
		}
		return nil
	},
}

func pullChanges() error {
	// Check if upstream is set
	upstreamCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err := upstreamCmd.Run(); err != nil {
		// No upstream, set it
		fmt.Println(styles.InfoIcon + " " + styles.Info.Render("No upstream branch configured"))

		err := runWithSpinner(styles.InfoIcon+" "+styles.Info.Render("Setting upstream and pulling from origin/HEAD"), func() error {
			fmt.Println() // Ensure newline before git output
			pullCmd := exec.Command("git", "pull", "--set-upstream", "origin", "HEAD")
			pullCmd.Stdout = os.Stdout
			pullCmd.Stderr = os.Stderr
			return pullCmd.Run()
		})

		if err != nil {
			fmt.Println(styles.ErrorIcon)
			return fmt.Errorf("failed to pull and set upstream: %w", err)
		}

		// Show success message
		fmt.Println()
		fmt.Println(styles.Card.Render(
			styles.Success.Render("Pull completed!") + "\n" +
				styles.Neutral.Render("Status: ") + styles.Success.Render("Upstream set and changes pulled"),
		))
		return nil
	}

	// Upstream exists, just pull
	err := runWithSpinner(styles.InfoIcon+" "+styles.Info.Render("Pulling changes from remote..."), func() error {
		fmt.Println() // Ensure newline before git output
		pullCmd := exec.Command("git", "pull")
		pullCmd.Stdout = os.Stdout
		pullCmd.Stderr = os.Stderr
		return pullCmd.Run()
	})

	if err != nil {
		fmt.Println(styles.ErrorIcon)
		return fmt.Errorf("failed to pull: %w", err)
	}

	// Show success message
	fmt.Println()
	fmt.Println(styles.Card.Render(
		styles.Success.Render("Pull completed!") + "\n" +
			styles.Neutral.Render("Status: ") + styles.Success.Render("Changes pulled from remote"),
	))
	return nil
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
