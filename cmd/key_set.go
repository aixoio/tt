package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aixoio/tt/styles"
)

var keySetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	Long:  styles.Info.Render("Set your API key, base URL, or default model for AI features."),
	RunE: func(cmd *cobra.Command, args []string) error {
		var configOption string
		var value string

		// Select configuration option
		selectForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title(styles.Primary.Render("Select Configuration Option")).
					Options(
						huh.NewOption("API Key", "api_key"),
						huh.NewOption("Base URL", "base_url"),
						huh.NewOption("Default Model", "default_model"),
					).
					Value(&configOption),
			),
		).WithTheme(huh.ThemeCharm())

		if err := selectForm.Run(); err != nil {
			return fmt.Errorf("failed to select option: %w", err)
		}

		// Create input form based on selection
		var input *huh.Input
		switch configOption {
		case "api_key":
			input = huh.NewInput().
				Title(styles.Primary.Render("OpenRouter API Key")).
				Placeholder("sk-or-v1-...").
				Description("Enter your OpenRouter API key (will be stored securely)").
				Value(&value).
				EchoMode(huh.EchoModePassword).
				Validate(func(s string) error {
					if len(s) < 10 {
						return fmt.Errorf("API key seems too short")
					}
					return nil
				})
		case "base_url":
			input = huh.NewInput().
				Title(styles.Primary.Render("Base URL")).
				Placeholder("https://openrouter.ai/api/v1").
				Description("Enter the base URL for the API").
				Value(&value).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("base URL cannot be empty")
					}
					return nil
				})
		case "default_model":
			input = huh.NewInput().
				Title(styles.Primary.Render("Default Model")).
				Placeholder("gpt-3.5-turbo").
				Description("Enter the default model for AI commits").
				Value(&value).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("default model cannot be empty")
					}
					return nil
				})
		}

		form := huh.NewForm(
			huh.NewGroup(input),
		).WithTheme(huh.ThemeCharm())

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to get value: %w", err)
		}

		if value == "" {
			fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Value cannot be empty"))
			return fmt.Errorf("value cannot be empty")
		}

		// Set the value in Viper
		viper.Set(configOption, value)

		// Write the config
		if err := viper.WriteConfig(); err != nil {
			fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Failed to save configuration"))
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Println(styles.SuccessIcon + " " + styles.Success.Render("Configuration updated successfully!"))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(keySetCmd)
}
