package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aixoio/tt/styles"
)

var apCmd = &cobra.Command{
	Use:     "ap",
	Aliases: []string{"aip", "aicommitpush"},
	Short:   "Generate AI commit message and push changes",
	Long:    styles.Info.Render("Generate an AI-powered commit message and push changes to the remote repository in one command."),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Show header
		fmt.Println(styles.Header.Render("AI Commit & Push"))
		fmt.Println()

		// Execute aicommit with auto-commit enabled
		fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Generating and creating AI commit..."))
		aicCommand := aicCmd
		// Set flags for auto-commit and add
		aicCommand.Flags().Set("commit", "true")
		aicCommand.Flags().Set("add", "true")

		// Execute the aic command
		if err := aicCommand.RunE(cmd, args); err != nil {
			return fmt.Errorf("failed to create AI commit: %w", err)
		}

		// Execute push
		fmt.Println()
		fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Pushing changes to remote..."))
		if err := pushChanges(); err != nil {
			return fmt.Errorf("failed to push changes: %w", err)
		}

		// Show final success message
		fmt.Println()
		fmt.Println(styles.Card.Render(
			styles.Success.Render("AI Commit & Push completed!") + "\n" +
				styles.Neutral.Render("Status: ") + styles.Success.Render("Changes committed and pushed"),
		))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(apCmd)
}
