package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/aixoio/tt/styles"
)

var addCmd = &cobra.Command{
	Use:     "add [paths...]",
	Aliases: []string{"ad"},
	Short:   "Stage files for commit",
	Long:    styles.Info.Render("Stage files in the repository. Works like `git add` with optional `--all` or explicit paths."),
	RunE: func(cmd *cobra.Command, args []string) error {
		allFlag, _ := cmd.Flags().GetBool("all")
		paths, _ := cmd.Flags().GetStringArray("path")

		// Show header
		fmt.Println(styles.Header.Render("Git Add"))
		fmt.Println()

		// Determine what to add
		var gitArgs []string
		if allFlag {
			gitArgs = []string{"add", "."}
		} else if len(paths) > 0 {
			gitArgs = append([]string{"add"}, paths...)
		} else if len(args) > 0 {
			gitArgs = append([]string{"add"}, args...)
		} else {
			fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("No files specified. Use --all or provide paths."))
			return fmt.Errorf("no files specified")
		}

		// Show what we're doing
		if allFlag {
			fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Staging all changes..."))
		} else {
			fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Staging specified files..."))
		}

		// Execute git add with spinner
		err := runWithSpinner(styles.InfoIcon+" "+styles.Info.Render("Adding files..."), func() error {
			gitCmd := exec.Command("git", gitArgs...)
			gitCmd.Stdout = os.Stdout
			gitCmd.Stderr = os.Stderr
			return gitCmd.Run()
		})

		if err != nil {
			fmt.Println(styles.ErrorIcon)
			return fmt.Errorf("failed to add files: %w", err)
		}

		fmt.Println(styles.SuccessIcon)

		// Show success message
		fmt.Println(styles.Card.Render(
			styles.Success.Render("Files staged successfully!") + "\n" +
				styles.Neutral.Render("Ready for commit."),
		))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("all", "A", false, "Stage all changes (git add .)")
	addCmd.Flags().StringArrayP("path", "p", []string{}, "Specific file or directory to stage (repeatable)")
}
