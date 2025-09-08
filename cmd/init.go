package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "creates a new tt repo",
	Long:  `creates a new git repo with tt`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := exec.Command("git", "init").Run()
		if err != nil {
			return fmt.Errorf("failed to initialize Git repository: %w", err)
		}
		fmt.Println("Initialized a new Git repository.")
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
