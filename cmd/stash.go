package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/aixoio/tt/styles"
)

var stashCmd = &cobra.Command{
	Use:     "stash [message]",
	Aliases: []string{"st"},
	Short:   "Stash changes with style",
	Long:    styles.Info.Render("Stash your changes with an interactive prompt for stash messages. Always includes untracked files for simplicity."),
	RunE: func(cmd *cobra.Command, args []string) error {
		message := ""
		if len(args) > 0 {
			message = args[0]
		}

		// Show header
		fmt.Println(styles.Header.Render("Git Stash"))
		fmt.Println()

		// Preview changes
		statusCmd := exec.Command("git", "status", "--porcelain")
		if output, err := statusCmd.Output(); err == nil && len(output) > 0 {
			fmt.Println(styles.Card.Render(
				styles.Info.Render("Files to be stashed:") + "\n" +
					styles.FilePath.Render(string(output)),
			))
		} else {
			fmt.Println(styles.InfoIcon + " " + styles.Info.Render("No changes to stash"))
			return nil
		}

		// Get stash message
		if message == "" {
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title(styles.Primary.Render("Stash Message")).
						Placeholder("Describe your work in progress...").
						Description("Optional message for your stash").
						Value(&message),
				),
			).WithTheme(huh.ThemeCharm())

			if err := form.Run(); err != nil {
				return fmt.Errorf("failed to get stash message: %w", err)
			}
		}

		// Execute stash
		fmt.Println()
		fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Stashing changes..."))
		stashArgs := []string{"stash", "push", "--include-untracked"}
		if message != "" {
			stashArgs = append(stashArgs, "-m", message)
		}
		gitCmd := exec.Command("git", stashArgs...)
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr

		if err := gitCmd.Run(); err != nil {
			fmt.Println(styles.ErrorIcon)
			return err
		}

		fmt.Println(styles.Card.Render(
			styles.Success.Render("Stash successful!") + "\n" +
				styles.Neutral.Render("Message: ") + styles.Highlight.Render(message),
		))

		return nil
	},
}

var stashPopCmd = &cobra.Command{
	Use:   "pop",
	Short: "Apply and remove the latest stash",
	Long:  styles.Info.Render("Apply the latest stash and remove it from the stash list. Shows confirmation and warns about potential conflicts."),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if stash exists
		listCmd := exec.Command("git", "stash", "list")
		output, err := listCmd.Output()
		if err != nil || len(output) == 0 {
			fmt.Println(styles.WarningIcon + " " + styles.Warning.Render("No stashes available"))
			return nil
		}

		// Show latest stash
		fmt.Println(styles.Header.Render("Stash Pop"))
		fmt.Println()
		index := bytes.IndexByte(output, '\n')
		var firstLine string
		if index == -1 {
			firstLine = string(output)
		} else {
			firstLine = string(output[:index+1])
		}
		fmt.Println(styles.Card.Render(
			styles.Info.Render("Latest stash:") + "\n" +
				styles.FilePath.Render(firstLine),
		))

		// Confirm
		confirm := false
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title(styles.Warning.Render("Apply this stash?")).
					Description("This may cause conflicts if files have changed").
					Value(&confirm),
			),
		).WithTheme(huh.ThemeCharm())

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}

		if !confirm {
			fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Stash pop cancelled"))
			return nil
		}

		// Execute pop
		fmt.Println()
		fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Applying stash..."))
		gitCmd := exec.Command("git", "stash", "pop")
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr

		if err := gitCmd.Run(); err != nil {
			fmt.Println(styles.WarningIcon + " " + styles.Warning.Render("Conflicts detected. Resolve them and commit when ready."))
			return err
		}

		fmt.Println(styles.SuccessIcon + " " + styles.Success.Render("Stash applied successfully"))

		return nil
	},
}

var stashListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stashes",
	Long:  styles.Info.Render("Show a simplified list of all your stashes with dates and messages."),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(styles.Header.Render("Stash List"))
		fmt.Println()

		gitCmd := exec.Command("git", "stash", "list", "--pretty=format:%C(yellow)%gd%C(reset) %C(green)%ci%C(reset) %s")
		output, err := gitCmd.Output()
		if err != nil {
			fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Failed to list stashes"))
			return err
		}

		if len(output) == 0 {
			fmt.Println(styles.InfoIcon + " " + styles.Info.Render("No stashes found"))
			return nil
		}

		fmt.Println(styles.Card.Render(
			styles.Info.Render("Your stashes:") + "\n" +
				styles.FilePath.Render(string(output)),
		))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(stashCmd)
	stashCmd.AddCommand(stashPopCmd)
	stashCmd.AddCommand(stashListCmd)
}
