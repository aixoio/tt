package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/aixoio/tt/styles"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new git repository",
	Long:  styles.Info.Render("Initialize a new git repository in the current directory."),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Show header
		fmt.Println(styles.Header.Render("Git Init"))
		fmt.Println()

		fmt.Print(styles.Spinner.Render("⏳") + " " + styles.Info.Render("Initializing git repository... "))

		err := exec.Command("git", "init").Run()
		if err != nil {
			fmt.Println(styles.ErrorIcon)
			return fmt.Errorf("failed to initialize Git repository: %w", err)
		}
		fmt.Println(styles.SuccessIcon)

		// Show success message with more details
		fmt.Println()
		fmt.Println(styles.Card.Render(
			styles.SuccessIcon + " " + styles.Success.Render("Git repository initialized!") + "\n" +
				styles.Neutral.Render("Status: ") + styles.Success.Render("Ready to start committing"),
		))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
