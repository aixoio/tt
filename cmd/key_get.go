package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aixoio/tt/styles"
)

var keyGetCmd = &cobra.Command{
	Use:   "key-get",
	Short: "Get the current OpenRouter API key (redacted)",
	Long:  styles.Info.Render("Display the current API key with sensitive parts redacted."),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := viper.GetString("api_key")

		if apiKey == "" {
			fmt.Println(styles.WarningIcon + " " + styles.Warning.Render("No API key set"))
			return nil
		}

		// Redact the API key (show first 4 and last 4 characters)
		redacted := redactAPIKey(apiKey)

		fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Current API Key: ") + styles.Highlight.Render(redacted))

		return nil
	},
}

func redactAPIKey(key string) string {
	if len(key) <= 8 {
		return strings.Repeat("*", len(key))
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}

func init() {
	rootCmd.AddCommand(keyGetCmd)
}
