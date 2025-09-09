package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/aixoio/tt/styles"
)

var tagCmd = &cobra.Command{
	Use:     "tag [command] [args]",
	Aliases: []string{"t"},
	Short:   "Create and manage git tags",
	Long:    styles.Info.Render("Manage git tags: create lightweight or annotated tags, list existing tags, or delete tags."),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return listTags()
		}

		if len(args) == 2 && args[0] == "delete" {
			return deleteTag(args[1])
		}

		if len(args) == 1 {
			return createTagInteractive(args[0], cmd)
		}

		// Invalid usage
		fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Invalid usage. Use 'tt tag' to list, 'tt tag <name>' to create, or 'tt tag delete <name>' to delete."))
		return fmt.Errorf("invalid command arguments")
	},
}

func listTags() error {
	fmt.Println(styles.Header.Render("Git Tags"))
	fmt.Println()

	// Get all tags
	fmt.Println(styles.Neutral.Render("Tags:"))
	tagsCmd := exec.Command("git", "tag", "-l", "--sort=-creatordate", "--format=%(refname:short)|%(creatordate:relative)|%(subject)")
	if tagsOutput, err := tagsCmd.Output(); err == nil {
		tags := strings.TrimSpace(string(tagsOutput))
		if tags != "" {
			for tagLine := range strings.SplitSeq(tags, "\n") {
				tagLine = strings.TrimSpace(tagLine)
				if tagLine != "" {
					parts := strings.Split(tagLine, "|")
					if len(parts) >= 2 {
						tagName := parts[0]
						tagDate := parts[1]
						tagMessage := ""
						if len(parts) > 2 {
							tagMessage = parts[2]
						}

						fmt.Printf("  • %s (%s)\n", styles.Branch.Render(tagName), styles.Muted.Render(tagDate))
						if tagMessage != "" {
							fmt.Printf("    └─ %s\n", styles.Info.Render(tagMessage))
						}
					}
				}
			}
		} else {
			fmt.Println("  " + styles.Muted.Render("No tags found"))
		}
	} else {
		fmt.Println("  " + styles.Muted.Render("Failed to retrieve tags"))
	}

	return nil
}

func createTagInteractive(name string, cmd *cobra.Command) error {
	message, _ := cmd.Flags().GetString("message")

	// If no message provided, create lightweight tag
	if message == "" {
		fmt.Println(styles.Header.Render("Create Lightweight Tag"))
		fmt.Println()
		fmt.Printf("%s %s\n", styles.InfoIcon, styles.Info.Render("Creating lightweight tag: "+styles.Branch.Render(name)))

		fmt.Print(styles.Spinner.Render("⏳") + " " + styles.Info.Render("Creating tag... "))
		tagCmd := exec.Command("git", "tag", name)
		if output, err := tagCmd.CombinedOutput(); err != nil {
			fmt.Println(styles.ErrorIcon)
			return fmt.Errorf("failed to create tag: %v", string(output))
		}
		fmt.Println(styles.SuccessIcon)

		fmt.Println()
		fmt.Println(styles.Card.Render(
			styles.Success.Render("Lightweight tag created successfully!") + "\n" +
				styles.Neutral.Render("Tag: ") + styles.Branch.Render(name),
		))
		return nil
	}

	// Create annotated tag with message
	fmt.Println(styles.Header.Render("Create Annotated Tag"))
	fmt.Println()
	fmt.Printf("%s %s\n", styles.InfoIcon, styles.Info.Render("Creating annotated tag: "+styles.Branch.Render(name)))

	fmt.Print(styles.Spinner.Render("⏳") + " " + styles.Info.Render("Creating annotated tag... "))
	tagCmd := exec.Command("git", "tag", "-a", name, "-m", message)
	if output, err := tagCmd.CombinedOutput(); err != nil {
		fmt.Println(styles.ErrorIcon)
		return fmt.Errorf("failed to create annotated tag: %v", string(output))
	}
	fmt.Println(styles.SuccessIcon)

	fmt.Println()
	fmt.Println(styles.Card.Render(
		styles.Success.Render("Annotated tag created successfully!") + "\n" +
			styles.Neutral.Render("Tag: ") + styles.Branch.Render(name) + "\n" +
			styles.Neutral.Render("Message: ") + styles.Highlight.Render(message),
	))

	return nil
}

func deleteTag(name string) error {
	// Verify tag exists
	verifyCmd := exec.Command("git", "tag", "-l", name)
	if verifyOutput, err := verifyCmd.Output(); err != nil || strings.TrimSpace(string(verifyOutput)) == "" {
		return fmt.Errorf("tag '%s' does not exist", name)
	}

	// Show header
	fmt.Println(styles.Header.Render("Delete Git Tag"))
	fmt.Println()

	// Get tag details
	fmt.Printf("%s %s\n", styles.WarningIcon, styles.Warning.Render("Deleting tag: "+styles.Branch.Render(name)))

	// Confirm deletion
	var confirm bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(styles.Warning.Render("Confirm Tag Deletion")).
				Description(styles.Neutral.Render("Are you sure you want to delete tag '" + name + "'? This action cannot be undone.")).
				Value(&confirm),
		),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}

	if !confirm {
		fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Tag deletion cancelled."))
		return nil
	}

	// Delete local tag
	fmt.Print(styles.Spinner.Render("⏳") + " " + styles.Info.Render("Deleting tag... "))
	deleteCmd := exec.Command("git", "tag", "-d", name)
	if output, err := deleteCmd.CombinedOutput(); err != nil {
		fmt.Println(styles.ErrorIcon)
		return fmt.Errorf("failed to delete tag: %v", string(output))
	}
	fmt.Println(styles.SuccessIcon)

	fmt.Println()
	fmt.Println(styles.Card.Render(
		styles.Success.Render("Tag deleted successfully!") + "\n" +
			styles.Neutral.Render("Tag: ") + styles.Branch.Render(name),
	))

	return nil
}

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.Flags().StringP("message", "m", "", "Tag annotation message")
}
