package cmd

import (
	"fmt"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/aixoio/tt/styles"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Hard-reset the repository after confirmation",
	Long:  styles.Info.Render("Perform a hard reset of the repository, discarding all uncommitted changes after user confirmation."),
	RunE: func(cmd *cobra.Command, args []string) error {
		var confirm bool

		// Build the confirmation prompt
		prompt := huh.NewConfirm().
			Title(styles.WarningIcon + " " + styles.Warning.Render("Reset Repository")).
			Description(
				"The following commands will be executed:\n\n" +
					"  " + styles.GitCommand.Render("git add .") + "\n" +
					"  " + styles.GitCommand.Render("git reset --hard") + "\n\n" +
					"This will discard all uncommitted changes. Continue?").
			Value(&confirm).
			Affirmative("Yes").
			Negative("No").
			WithTheme(huh.ThemeCharm())

		// Show the prompt
		if err := prompt.Run(); err != nil {
			return fmt.Errorf("failed to show confirmation prompt: %w", err)
		}

		if !confirm {
			fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Aborted – no changes were made."))
			return nil
		}

		// Execute git add .
		fmt.Print(styles.SpinnerIcon + " " + styles.Info.Render("Staging all files... "))
		if out, err := exec.Command("git", "add", ".").CombinedOutput(); err != nil {
			fmt.Println(styles.ErrorIcon)
			fmt.Println(styles.Error.Render("git add failed:"), string(out))
			return fmt.Errorf("git add failed: %w", err)
		}
		fmt.Println(styles.SuccessIcon)

		// Execute git reset --hard
		fmt.Print(styles.SpinnerIcon + " " + styles.Info.Render("Performing hard reset... "))
		if out, err := exec.Command("git", "reset", "--hard").CombinedOutput(); err != nil {
			fmt.Println(styles.ErrorIcon)
			fmt.Println(styles.Error.Render("git reset --hard failed:"), string(out))
			return fmt.Errorf("git reset --hard failed: %w", err)
		}
		fmt.Println(styles.SuccessIcon)

		fmt.Println(styles.Card.Render(
			styles.Success.Render("Repository successfully reset.") + "\n" +
				styles.Neutral.Render("All uncommitted changes have been discarded."),
		))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
