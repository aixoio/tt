package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/aixoio/tt/styles"
)

var tagPushCmd = &cobra.Command{
	Use:     "push [tag]",
	Aliases: []string{"tp"},
	Short:   "Push git tags to remote repository",
	Long:    styles.Info.Render("Push tags to the remote repository. Push all tags or a specific tag."),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return pushAllTags()
		}

		if len(args) == 1 {
			return pushSpecificTag(args[0])
		}

		// Invalid usage
		fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Invalid usage. Use 'tt tag push' to push all tags, or 'tt tag push <tag>' to push a specific tag."))
		return fmt.Errorf("invalid command arguments")
	},
}

func pushAllTags() error {
	fmt.Println(styles.Header.Render("Push All Tags"))
	fmt.Println()

	fmt.Print(styles.Spinner.Render("⏳") + " " + styles.Info.Render("Pushing all tags to remote... "))
	pushCmd := exec.Command("git", "push", "--tags")
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr

	if err := pushCmd.Run(); err != nil {
		fmt.Println(styles.ErrorIcon)
		return fmt.Errorf("failed to push all tags: %w", err)
	}

	fmt.Println(styles.SuccessIcon)
	fmt.Println()
	fmt.Println(styles.Card.Render(
		styles.Success.Render("All tags pushed successfully!") + "\n" +
			styles.Neutral.Render("Status: ") + styles.Success.Render("All tags pushed to remote"),
	))

	return nil
}

func pushSpecificTag(tagName string) error {
	fmt.Println(styles.Header.Render("Push Specific Tag"))
	fmt.Println()

	// Verify tag exists
	verifyCmd := exec.Command("git", "tag", "-l", tagName)
	if verifyOutput, err := verifyCmd.Output(); err != nil || string(verifyOutput) == "" {
		return fmt.Errorf("tag '%s' does not exist locally", tagName)
	}

	fmt.Printf("%s %s\n", styles.InfoIcon, styles.Info.Render("Pushing tag: "+styles.Branch.Render(tagName)))
	fmt.Print(styles.Spinner.Render("⏳") + " " + styles.Info.Render("Pushing tag to remote... "))

	pushCmd := exec.Command("git", "push", "origin", tagName)
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr

	if err := pushCmd.Run(); err != nil {
		fmt.Println(styles.ErrorIcon)
		return fmt.Errorf("failed to push tag '%s': %w", tagName, err)
	}

	fmt.Println(styles.SuccessIcon)
	fmt.Println()
	fmt.Println(styles.Card.Render(
		styles.Success.Render("Tag pushed successfully!") + "\n" +
			styles.Neutral.Render("Tag: ") + styles.Branch.Render(tagName) + "\n" +
			styles.Neutral.Render("Status: ") + styles.Success.Render("Tag pushed to remote"),
	))

	return nil
}

func init() {
	tagCmd.AddCommand(tagPushCmd)
}
