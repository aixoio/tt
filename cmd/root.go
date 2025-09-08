package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/aixoio/tt/styles"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tt",
	Short: "A beautiful git helper tool",
	Long: styles.Header.Render("tt - Git Helper Tool") + `

A beautifully styled, 100% git compatible tool built
to make git operations intuitive and enjoyable for developers.

Features:
• ` + styles.Highlight.Render(`Interactive prompts`) + ` for commit messages and branch names
• ` + styles.Highlight.Render(`Smart defaults`) + ` for common git operations
• ` + styles.Highlight.Render(`Beautiful styling`) + ` with clear visual feedback
• ` + styles.Highlight.Render(`Auto-push`) + ` and upstream management
• ` + styles.Highlight.Render(`Conflict-aware`) + ` merge operations

Get started by running: ` + styles.InlineCode.Render(`tt init`),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tt.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
