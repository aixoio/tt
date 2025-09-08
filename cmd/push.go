package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/aixoio/tt/styles"
)

var pushCmd = &cobra.Command{
	Use:     "push",
	Aliases: []string{"p"},
	Short:   "Push changes to remote repository",
	Long:    styles.Info.Render("Push changes to the remote repository. Automatically sets upstream if not configured."),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Show header
		fmt.Println(styles.Header.Render("Git Push"))
		fmt.Println()

		if err := pushChanges(); err != nil {
			return fmt.Errorf("failed to push: %w", err)
		}
		return nil
	},
}

func pushChanges() error {
	// Check if upstream is set
	upstreamCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err := upstreamCmd.Run(); err != nil {
		// No upstream, set it
		fmt.Println(styles.InfoIcon + " " + styles.Info.Render("No upstream branch configured"))
		fmt.Print(styles.Spinner.Render("⏳") + " " + styles.Info.Render("Setting upstream to origin/HEAD... "))

		pushCmd := exec.Command("git", "push", "--set-upstream", "origin", "HEAD")
		pushCmd.Stdout = os.Stdout
		pushCmd.Stderr = os.Stderr
		if err := pushCmd.Run(); err != nil {
			fmt.Println(styles.ErrorIcon)
			return fmt.Errorf("failed to push and set upstream: %w", err)
		}
		fmt.Println(styles.SuccessIcon)

		// Show success message
		fmt.Println()
		fmt.Println(styles.Card.Render(
			styles.SuccessIcon + " " + styles.Success.Render("Push completed!") + "\n" +
				styles.Neutral.Render("Status: ") + styles.Success.Render("Upstream set and changes pushed"),
		))
		return nil
	}

	// Upstream exists, just push
	fmt.Print(styles.Spinner.Render("⏳") + " " + styles.Info.Render("Pushing changes to remote... "))

	pushCmd := exec.Command("git", "push")
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	if err := pushCmd.Run(); err != nil {
		fmt.Println(styles.ErrorIcon)
		return fmt.Errorf("failed to push: %w", err)
	}
	fmt.Println(styles.SuccessIcon)

	// Show success message
	fmt.Println()
	fmt.Println(styles.Card.Render(
		styles.SuccessIcon + " " + styles.Success.Render("Push completed!") + "\n" +
			styles.Neutral.Render("Status: ") + styles.Success.Render("Changes pushed to remote"),
	))
	return nil
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
