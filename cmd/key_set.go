package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aixoio/tt/styles"
)

var keySetCmd = &cobra.Command{
	Use:   "key-set",
	Short: "Set the OpenRouter API key",
	Long:  styles.Info.Render("Securely set your OpenRouter API key for AI commit messages."),
	RunE: func(cmd *cobra.Command, args []string) error {
		var apiKey string

		// Create styled password input form
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title(styles.Primary.Render("OpenRouter API Key")).
					Placeholder("sk-or-v1-...").
					Description("Enter your OpenRouter API key (will be stored securely)").
					Value(&apiKey).
					EchoMode(huh.EchoModePassword).
					Validate(func(s string) error {
						if len(s) < 10 {
							return fmt.Errorf("API key seems too short")
						}
						return nil
					}),
			),
		).WithTheme(huh.ThemeCharm())

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		if apiKey == "" {
			fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("API key cannot be empty"))
			return fmt.Errorf("API key cannot be empty")
		}

		// Set the API key in Viper
		viper.Set("api_key", apiKey)

		// Write the config
		if err := viper.WriteConfig(); err != nil {
			fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Failed to save API key"))
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Println(styles.SuccessIcon + " " + styles.Success.Render("API key set successfully!"))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(keySetCmd)
}
