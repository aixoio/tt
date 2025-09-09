package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aixoio/tt/styles"
)

var keyGetCmd = &cobra.Command{
	Use:   "key-get",
	Short: "Get the current configuration values",
	Long:  styles.Info.Render("Display the current API key, base URL, and default model."),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := viper.GetString("api_key")
		baseURL := viper.GetString("base_url")
		defaultModel := viper.GetString("default_model")

		fmt.Println(styles.Header.Render("Current Configuration"))
		fmt.Println()

		if apiKey == "" {
			fmt.Println(styles.WarningIcon + " " + styles.Warning.Render("API Key: Not set"))
		} else {
			fmt.Println(styles.InfoIcon + " " + styles.Info.Render("API Key: ") + styles.Highlight.Render(apiKey))
		}

		fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Base URL: ") + styles.Highlight.Render(baseURL))
		fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Default Model: ") + styles.Highlight.Render(defaultModel))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(keyGetCmd)
}
