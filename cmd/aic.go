package cmd

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

// getGitDiff gets the current changes in the git repository
func getGitDiff() (string, error) {
	// Check if git is installed
	_, err := exec.LookPath("git")
	if err != nil {
		return "", fmt.Errorf("git is not installed or not in PATH")
	}

	// Check if current directory is a git repository
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("current directory is not a git repository")
	}

	// Get staged changes
	stagedCmd := exec.Command("git", "diff", "--staged")
	stagedOutput, err := stagedCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get staged changes: %w", err)
	}

	// Get unstaged changes if no staged changes
	if len(stagedOutput) == 0 {
		unstagedCmd := exec.Command("git", "diff")
		unstagedOutput, err := unstagedCmd.Output()
		if err != nil {
			return "", fmt.Errorf("failed to get unstaged changes: %w", err)
		}

		if len(unstagedOutput) == 0 {
			return "", fmt.Errorf("no changes detected in the repository")
		}

		return string(unstagedOutput), nil
	}

	return string(stagedOutput), nil
}

// getChangedFiles gets the names of files that have been changed
func getChangedFiles() ([]string, error) {
	// Check if git is installed
	_, err := exec.LookPath("git")
	if err != nil {
		return nil, fmt.Errorf("git is not installed or not in PATH")
	}

	// Get staged files
	stagedCmd := exec.Command("git", "diff", "--staged", "--name-only")
	stagedOutput, err := stagedCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}

	// Get unstaged files if no staged files
	if len(stagedOutput) == 0 {
		unstagedCmd := exec.Command("git", "diff", "--name-only")
		unstagedOutput, err := unstagedCmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to get unstaged files: %w", err)
		}

		if len(unstagedOutput) == 0 {
			return nil, fmt.Errorf("no changed files detected in the repository")
		}

		return strings.Split(strings.TrimSpace(string(unstagedOutput)), "\n"), nil
	}

	return strings.Split(strings.TrimSpace(string(stagedOutput)), "\n"), nil
}

// getProjectInfo gets information about the project
func getProjectInfo() (string, error) {
	// Try to determine the project type based on files
	files, err := filepath.Glob("*")
	if err != nil {
		return "", fmt.Errorf("failed to list files: %w", err)
	}

	var projectInfo strings.Builder
	projectInfo.WriteString("Project files include: ")

	// Look for specific project indicators
	hasGoMod := false
	hasPackageJSON := false
	hasPomXML := false
	hasCMake := false
	hasPyProject := false

	for _, file := range files {
		switch file {
		case "go.mod":
			hasGoMod = true
		case "package.json":
			hasPackageJSON = true
		case "pom.xml":
			hasPomXML = true
		case "CMakeLists.txt":
			hasCMake = true
		case "pyproject.toml":
			hasPyProject = true
		}
	}

	if hasGoMod {
		projectInfo.WriteString("Go project. ")
	}
	if hasPackageJSON {
		projectInfo.WriteString("JavaScript/Node.js project. ")
	}
	if hasPomXML {
		projectInfo.WriteString("Java/Maven project. ")
	}
	if hasCMake {
		projectInfo.WriteString("C/C++ project with CMake. ")
	}
	if hasPyProject {
		projectInfo.WriteString("Python project. ")
	}

	return projectInfo.String(), nil
}

// readUserInput reads a single character from the user
func readUserInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return "y", nil
	}
	return strings.ToLower(input), nil
}

// generateCommitMessage uses OpenAI to generate a commit message based on git diff and project information
func generateCommitMessage(apiKey, baseURL, model, diff string) (string, error) {
	if model == "" {
		model = viper.GetString("default_model")
	}

	// Get changed files for more context
	changedFiles, err := getChangedFiles()
	if err != nil {
		// Non-fatal error, we can continue without this info
		fmt.Printf("%s Warning: couldn't get changed files: %v\n", styles.WarningIcon, err)
	}

	// Get project information for more context
	projectInfo, err := getProjectInfo()
	if err != nil {
		// Non-fatal error, we can continue without this info
		fmt.Printf("%s Warning: couldn't get project info: %v\n", styles.WarningIcon, err)
	}

	// Build file list string
	var fileListStr string
	if len(changedFiles) > 0 {
		fileListStr = fmt.Sprintf("Changed files: %s\n\n", strings.Join(changedFiles, ", "))
	}

	// Prepare the prompt with more context
	prompt := "Generate a short, concise git commit message based on the following changes. " +
		"Follow the conventional commit format (e.g., feat:, fix:, docs:, style:, refactor:, test:, chore:). " +
		"Keep it under 50 characters if possible. " +
		"Only respond with the commit message, nothing else.\n\n"

	if projectInfo != "" {
		prompt += "Project information: " + projectInfo + "\n\n"
	}

	prompt += fileListStr + "Changes:\n" + diff

	// Initialize OpenAI client with OpenRouter
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL
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

	// Create request
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate commit message: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI model")
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

// makeCommit creates a git commit with the provided message
func makeCommit(message string) error {
	// Stage all changes
	addCmd := exec.Command("git", "add", ".")
	addCmd.Stdout = os.Stdout
	addCmd.Stderr = os.Stderr
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// Create commit
	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr
	return commitCmd.Run()
}

var (
	autoCommit bool
	model      string
)

var aicCmd = &cobra.Command{
	Use:     "aic",
	Aliases: []string{"ai-commit", "ai", "ac"},
	Short:   "Generate AI-powered commit messages",
	Long:    styles.Info.Render("Use AI to generate conventional commit messages based on your changes."),
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

		// Get git diff
		fmt.Print(styles.InfoIcon + " " + styles.Info.Render("Analyzing changes... "))
		diff, err := getGitDiff()
		if err != nil {
			fmt.Println(styles.ErrorIcon)
			return err
		}
		fmt.Println(styles.SuccessIcon)

		// Print which model is being used
		modelToUse := model
		if model == "" {
			modelToUse = viper.GetString("default_model")
		}

		fmt.Println()
		fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Using model: ") + styles.Highlight.Render(modelToUse))

		// Generate commit message
		fmt.Print(styles.Spinner.Render("🤖") + " " + styles.Info.Render("Generating commit message... "))
		message, err := generateCommitMessage(apiKey, viper.GetString("base_url"), model, diff)
		if err != nil {
			fmt.Println(styles.ErrorIcon)
			return err
		}
		fmt.Println(styles.SuccessIcon)

		// Output commit message with prominent formatting
		fmt.Println()
		fmt.Println(styles.Card.Render(
			styles.Success.Render("Generated Commit Message:") + "\n" +
				styles.Highlight.Render(message),
		))

		// Handle commit based on auto-commit flag or user confirmation
		if autoCommit {
			// Auto-commit mode - commit without confirmation
			if err := makeCommit(message); err != nil {
				return err
			}
			fmt.Println(styles.SuccessIcon + " " + styles.Success.Render("Commit created successfully"))
		} else {
			// Ask for confirmation with additional options
			fmt.Println()
			fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Options:"))
			fmt.Println(styles.Neutral.Render("  y/yes - Create commit with this message"))
			fmt.Println(styles.Neutral.Render("  n/no  - Cancel commit"))
			fmt.Println(styles.Neutral.Render("  g     - Generate more detailed message"))
			fmt.Println(styles.Neutral.Render("  r     - Retry with new generation"))
			fmt.Println(styles.Neutral.Render("  s     - Summarize message"))
			fmt.Println(styles.Neutral.Render("  p     - Provide feedback for the message"))

			for {
				fmt.Print(styles.Primary.Render("Create commit with this message? [y/n/g/r/s/p]: "))

				response, err := readUserInput()
				if err != nil {
					return fmt.Errorf("error reading user input: %w", err)
				}

				if response == "y" || response == "yes" {
					if err := makeCommit(message); err != nil {
						return err
					}
					fmt.Println(styles.SuccessIcon + " " + styles.Success.Render("Commit created successfully"))
					break
				} else if response == "n" || response == "no" {
					fmt.Println(styles.InfoIcon + " " + styles.Info.Render("Commit canceled"))
					break
				} else if response == "g" {
					fmt.Print(styles.Spinner.Render("🔍") + " " + styles.Info.Render("Generating a more detailed commit message... "))
					message, err = generateCommitMessage(apiKey, viper.GetString("base_url"), model, diff+"\n\nPlease provide a more detailed commit message with additional context and explanations.")
					if err != nil {
						fmt.Println(styles.ErrorIcon)
						return err
					}
					fmt.Println(styles.SuccessIcon)
					fmt.Println()
					fmt.Println(styles.Card.Render(
						styles.Success.Render("Generated Detailed Commit Message:") + "\n" +
							styles.Highlight.Render(message),
					))
				} else if response == "r" {
					fmt.Print(styles.Spinner.Render("🔄") + " " + styles.Info.Render("Retrying with a new generation... "))
					message, err = generateCommitMessage(apiKey, viper.GetString("base_url"), model, diff)
					if err != nil {
						fmt.Println(styles.ErrorIcon)
						return err
					}
					fmt.Println(styles.SuccessIcon)
					fmt.Println()
					fmt.Println(styles.Card.Render(
						styles.Success.Render("Regenerated Commit Message:") + "\n" +
							styles.Highlight.Render(message),
					))
				} else if response == "s" {
					fmt.Print(styles.Spinner.Render("📝") + " " + styles.Info.Render("Summarizing the commit message... "))
					message, err = generateCommitMessage(apiKey, viper.GetString("base_url"), model, "Please summarize this commit message in 50 characters or less:\n\n"+message)
					if err != nil {
						fmt.Println(styles.ErrorIcon)
						return err
					}
					fmt.Println(styles.SuccessIcon)
					fmt.Println()
					fmt.Println(styles.Card.Render(
						styles.Success.Render("Summarized Commit Message:") + "\n" +
							styles.Highlight.Render(message),
					))
				} else if response == "p" {
					fmt.Print(styles.Spinner.Render("🔍") + " " + styles.Info.Render("Enter your feedback for the commit message: "))
					reader := bufio.NewReader(os.Stdin)
					feedbackLine, err := reader.ReadString('\n')
					if err != nil {
						return fmt.Errorf("error reading feedback: %w", err)
					}
					feedback := strings.TrimSpace(feedbackLine)

					fmt.Print(styles.Spinner.Render("🎯") + " " + styles.Info.Render("Generating commit message based on your feedback... "))
					promptWithGuidance := "Based on this diff:\n\n" + diff + "\n\nAnd considering this feedback: " + feedback + "\n\nGenerate an appropriate commit message."
					message, err = generateCommitMessage(apiKey, viper.GetString("base_url"), model, promptWithGuidance)
					if err != nil {
						fmt.Println(styles.ErrorIcon)
						return err
					}
					fmt.Println(styles.SuccessIcon)
					fmt.Println()
					fmt.Println(styles.Card.Render(
						styles.Success.Render("Feedback-Based Commit Message:") + "\n" +
							styles.Highlight.Render(message),
					))
				} else {
					fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Invalid option. Please choose y (yes), n (no), g (generate detailed), r (retry), s (shorter), or p (custom prompt)."))
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(aicCmd)
	aicCmd.Flags().BoolVarP(&autoCommit, "commit", "c", false, "Automatically create commit with generated message")
	aicCmd.Flags().StringVarP(&model, "model", "m", "", "OpenRouter model to use for generation (overrides default_model from config)")
}
