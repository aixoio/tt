package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/aixoio/tt/styles"
)

var commitCmd = &cobra.Command{
	Use:     "c [message]",
	Aliases: []string{"commit"},
	Short:   "Commit changes with style",
	Long:    styles.Info.Render("Commit changes to git with an interactive prompt for commit messages. Supports automatic file staging and pushing."),
	RunE: func(cmd *cobra.Command, args []string) error {
		message, _ := cmd.Flags().GetString("message")
		addFlag, _ := cmd.Flags().GetBool("add")

		// Show header
		fmt.Println(styles.Header.Render("Git Commit"))
		fmt.Println()

		// Handle file staging
		if addFlag {
			fmt.Print(styles.InfoIcon + " " + styles.Info.Render("Staging all files... "))
			if err := exec.Command("git", "add", ".").Run(); err != nil {
				fmt.Println(styles.ErrorIcon)
				return fmt.Errorf("failed to add files: %w", err)
			}
			fmt.Println(styles.SuccessIcon)
		}

		// Get commit message
		if message == "" {
			// Show current status before prompting
			statusCmd := exec.Command("git", "status", "--porcelain")
			if output, err := statusCmd.Output(); err == nil && len(output) > 0 {
				fmt.Println(styles.Card.Render(
					styles.Info.Render("Files to be committed:") + "\n" +
						styles.FilePath.Render(string(output)),
				))
			}

			// Create styled input form
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title(styles.Primary.Render("Commit Message")).
						Placeholder("Describe your changes...").
						Description("Write a clear, concise commit message").
						Value(&message).
						Validate(func(s string) error {
							if len(s) < 3 {
								return fmt.Errorf("commit message too short")
							}
							return nil
						}),
				),
			).WithTheme(huh.ThemeCharm())

			if err := form.Run(); err != nil {
				return fmt.Errorf("failed to get commit message: %w", err)
			}
		}

		if message == "" {
			fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Commit message cannot be empty"))
			return fmt.Errorf("commit message cannot be empty")
		}

		// Show commit details
		fmt.Println()
		fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Committing with message: ") + styles.Highlight.Render("\""+message+"\""))

		// Execute commit
		fmt.Print(styles.Spinner.Render("⏳") + " " + styles.Info.Render("Creating commit... "))
		gitCmd := exec.Command("git", "commit", "-m", message)
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr

		if err := gitCmd.Run(); err != nil {
			fmt.Println(styles.ErrorIcon)
			return fmt.Errorf("failed to commit: %w", err)
		}
		fmt.Println(styles.SuccessIcon)

		// Show success message with commit info
		fmt.Println()
		fmt.Println(styles.Card.Render(
			styles.SuccessIcon + " " + styles.Success.Render("Commit successful!") + "\n" +
				styles.Neutral.Render("Message: ") + styles.Highlight.Render(message),
		))

		// Handle auto-push
		pushFlag, _ := cmd.Flags().GetBool("push")
		if pushFlag {
			fmt.Println()
			fmt.Print(styles.InfoIcon + " " + styles.Info.Render("Pushing changes... "))
			if err := pushChanges(); err != nil {
				fmt.Println(styles.ErrorIcon)
				fmt.Println(styles.WarningIcon + " " + styles.Warning.Render("Push failed, but commit was successful"))
				return fmt.Errorf("failed to push after commit: %w", err)
			}
			fmt.Println(styles.SuccessIcon)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringP("message", "m", "", "Commit message")
	commitCmd.Flags().BoolP("add", "a", false, "Add all files before committing")
	commitCmd.Flags().BoolP("push", "p", false, "Push after committing")
}
