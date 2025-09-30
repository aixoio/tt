package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aixoio/tt/styles"
)

var (
	aiFlag   bool
	statFlag bool
	nameOnly bool
)

var diffCmd = &cobra.Command{
	Use:     "diff",
	Aliases: []string{"d"},
	Short:   "Show changes between commits",
	Long:    styles.Info.Render("Display the differences between commits or working tree, with optional AI summary and enhanced styling."),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoCmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
		if err := repoCmd.Run(); err != nil {
			return fmt.Errorf("not a git repository: %w", err)
		}

		fullDiffCmd := exec.Command("git", append([]string{"diff", "--no-color"}, args...)...)
		fullDiffOutput, err := fullDiffCmd.Output()
		if err != nil {
			return fmt.Errorf("failed to get git diff: %w", err)
		}
		diffContent := string(fullDiffOutput)

		if diffContent == "" {
			fmt.Println(styles.SuccessIcon + " " + styles.Success.Render("No changes to display"))
			return nil
		}

		fmt.Println(styles.Header.Render("Git Diff"))
		fmt.Println()

		if statFlag {
			statCmd := exec.Command("git", append([]string{"diff", "--stat"}, args...)...)
			statOutput, err := statCmd.Output()
			if err != nil {
				fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Failed to get diff stat"))
			} else {
				fmt.Print(styles.Info.Render("Changes summary:\n"))
				fmt.Println(string(statOutput))
				fmt.Println()
			}
		}

		if nameOnly {
			nameCmd := exec.Command("git", append([]string{"diff", "--name-only"}, args...)...)
			nameOutput, err := nameCmd.Output()
			if err != nil {
				return fmt.Errorf("failed to get file names: %w", err)
			}
			names := strings.Split(strings.TrimSpace(string(nameOutput)), "\n")
			for _, name := range names {
				if name != "" {
					fmt.Println("  " + styles.FilePath.Render(name))
				}
			}
			return nil
		}

		lines := strings.Split(diffContent, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "diff --git") || strings.HasPrefix(line, "index") || strings.HasPrefix(line, "@@") {
				fmt.Println(styles.DiffHeader.Render(line))
			} else if strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---") {
				fmt.Println(styles.DiffHeader.Render(line))
			} else if strings.HasPrefix(line, "+") {
				fmt.Println(styles.Add.Render(line))
			} else if strings.HasPrefix(line, "-") {
				fmt.Println(styles.Del.Render(line))
			} else {
				fmt.Println(styles.Neutral.Render(line))
			}
		}

		if aiFlag {
			apiKey := viper.GetString("api_key")
			if apiKey == "" {
				fmt.Println(styles.WarningIcon + " " + styles.Warning.Render("API key not set. Skipping AI overview. Run 'tt set' to configure."))
			} else {
				model := viper.GetString("default_model")

				// Inline getProjectInfo
				projectInfoStr := ""
				files, err := filepath.Glob("*")
				if err != nil {
					fmt.Printf("%s Warning: couldn't get project info: %v\n", styles.WarningIcon, err)
				} else {
					var projectInfo strings.Builder
					projectInfo.WriteString("Project files include: ")

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
					projectInfoStr = projectInfo.String()
				}

				// Inline getChangedFiles
				changedFiles := []string{}
				changedCmd := exec.Command("git", append([]string{"diff", "--name-only"}, args...)...)
				changedOutput, err := changedCmd.Output()
				if err != nil {
					fmt.Printf("%s Warning: couldn't get changed files: %v\n", styles.WarningIcon, err)
				} else {
					namesList := strings.Split(strings.TrimSpace(string(changedOutput)), "\n")
					for _, n := range namesList {
						if n != "" {
							changedFiles = append(changedFiles, n)
						}
					}
				}

				// Build prompt
				var sb strings.Builder
				sb.WriteString("Provide a concise summary (2-3 sentences) of the code changes in this git diff. Focus on what the changes achieve, such as new features, bug fixes, or refactors. Only respond with the summary, nothing else.\n\n")
				if projectInfoStr != "" {
					sb.WriteString(projectInfoStr + "\n\n")
				}
				if len(changedFiles) > 0 {
					sb.WriteString(fmt.Sprintf("Changed files: %s\n\n", strings.Join(changedFiles, ", ")))
				}
				sb.WriteString("Changes:\n" + diffContent)
				basePrompt := sb.String()

				summary, err := runWithSpinnerForMessage("🤖 Generating AI overview...", func() (string, error) {
					return generateAIResponse(apiKey, viper.GetString("base_url"), model, basePrompt)
				})
				if err != nil {
					fmt.Println(styles.ErrorIcon + " " + styles.Error.Render("Failed to generate AI overview"))
				} else {
					fmt.Println()
					fmt.Println(styles.Info.Render("AI Overview:"))
					fmt.Println(styles.Info.Render(summary))
					fmt.Println()
				}
			}
		}

		return nil
	},
}

func generateAIResponse(apiKey, baseURL, model, prompt string) (string, error) {
	if model == "" {
		model = viper.GetString("default_model")
	}

	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithHeader("HTTP-Referer", "https://github.com/aixoio/tt"),
		option.WithHeader("X-Title", "tt"),
		option.WithAPIKey(apiKey),
	)

	resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model:    model,
		Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage(prompt)},
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate AI response: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI model")
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().BoolVarP(&aiFlag, "ai", "a", false, "Generate AI-powered overview of changes")
	diffCmd.Flags().BoolVarP(&statFlag, "stat", "s", false, "Show stat summary")
	diffCmd.Flags().BoolVarP(&nameOnly, "name-only", "n", false, "Show only names of changed files")
}
