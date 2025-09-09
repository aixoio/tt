package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aixoio/tt/styles"
)

type headerTransport struct {
	rt      http.RoundTripper
	headers map[string]string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Set(k, v)
	}
	return t.rt.RoundTrip(req)
}

var aicCmd = &cobra.Command{
	Use:     "aic",
	Aliases: []string{"ai-commit", "ai", "ac"},
	Short:   "Generate AI-powered commit messages",
	Long:    styles.Info.Render("Use AI to generate conventional commit messages based on your staged changes."),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if API key is set
		apiKey := viper.GetString("api_key")
		if apiKey == "" {
			fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("API key not set. Run 'tt key-set' to configure it."))
			return fmt.Errorf("API key not set")
		}

		// Show header
		fmt.Println(styles.Header.Render("AI Commit"))
		fmt.Println()

		// Get staged changes
		fmt.Print(styles.InfoIcon + " " + styles.Info.Render("Analyzing staged changes... "))
		gitCmd := exec.Command("git", "diff", "--staged")
		diffOutput, err := gitCmd.Output()
		if err != nil {
			fmt.Println(styles.ErrorIcon)
			return fmt.Errorf("failed to get staged changes: %w", err)
		}
		fmt.Println(styles.SuccessIcon)

		if len(diffOutput) == 0 {
			fmt.Println(styles.WarningIcon + " " + styles.Warning.Render("No staged changes found. Stage some files first."))
			return fmt.Errorf("no staged changes")
		}

		// Initialize OpenAI client with OpenRouter
		config := openai.DefaultConfig(apiKey)
		config.BaseURL = viper.GetString("base_url")
		config.HTTPClient = &http.Client{
			Transport: &headerTransport{
				rt: http.DefaultTransport,
				headers: map[string]string{
					"HTTP-Referer": "https://github.com/aixoio/tt",
					"X-Title":      "tt",
				},
			},
		}
		client := openai.NewClientWithConfig(config)

		// Prepare prompt
		prompt := "Generate a conventional commit message for the following changes. Follow conventional commit format (type(scope): description). Keep it concise and descriptive:\n\n" + string(diffOutput)

		// Call AI
		fmt.Print(styles.Spinner.Render("🤖") + " " + styles.Info.Render("Generating commit message... "))
		resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Model: viper.GetString("default_model"),
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		})
		if err != nil {
			fmt.Println(styles.ErrorIcon)
			return fmt.Errorf("failed to generate commit message: %w", err)
		}
		fmt.Println(styles.SuccessIcon)

		aiMessage := resp.Choices[0].Message.Content
		if aiMessage == "" {
			fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("AI returned empty message"))
			return fmt.Errorf("empty AI response")
		}

		// Display the generated message
		fmt.Println()
		fmt.Println(styles.Card.Render(
			styles.Success.Render("Generated Commit Message:") + "\n" +
				styles.Highlight.Render(aiMessage),
		))

		// Confirm and commit
		var confirm bool
		if err := huh.NewConfirm().
			Title(styles.Primary.Render("Use this commit message?")).
			Description("This will commit your staged changes with the AI-generated message").
			Value(&confirm).
			Run(); err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}

		if !confirm {
			fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Commit cancelled"))
			return nil
		}

		// Execute commit
		fmt.Print(styles.Spinner.Render("⏳") + " " + styles.Info.Render("Committing... "))
		commitCmd := exec.Command("git", "commit", "-m", aiMessage)
		commitCmd.Stdout = os.Stdout
		commitCmd.Stderr = os.Stderr

		if err := commitCmd.Run(); err != nil {
			fmt.Println(styles.ErrorIcon)
			return fmt.Errorf("failed to commit: %w", err)
		}
		fmt.Println(styles.SuccessIcon)

		fmt.Println(styles.Card.Render(
			styles.Success.Render("Commit successful!") + "\n" +
				styles.Neutral.Render("Message: ") + styles.Highlight.Render(aiMessage),
		))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(aicCmd)
}
