package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/aixoio/tt/styles"
)

var pushCmd = &cobra.Command{
	Use:     "push",
	Aliases: []string{"p"},
	Short:   "Push changes to remote repository",
	Long:    styles.Info.Render("Push changes to the remote repository. Automatically sets upstream if not configured."),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Show header
		fmt.Println(styles.Header.Render("Git Push"))
		fmt.Println()

		if err := pushChanges(); err != nil {
			return fmt.Errorf("failed to push: %w", err)
		}
		return nil
	},
}

func runWithSpinner(title string, action func() error) error {
	fmt.Print("\033[?25l") // hide cursor
	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		ticker := time.Tick(100 * time.Millisecond)
		for {
			select {
			case <-stop:
				return
			case <-ticker:
				fmt.Printf("\r%s %s", title, spinners[i])
				i = (i + 1) % len(spinners)
			}
		}
	}()
	err := action()
	close(stop)
	wg.Wait()
	fmt.Println()
	fmt.Print("\033[?25h") // show cursor
	return err
}

func pushChanges() error {
	// Check if upstream is set
	upstreamCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err := upstreamCmd.Run(); err != nil {
		// No upstream, set it
		fmt.Println(styles.InfoIcon + " " + styles.Info.Render("No upstream branch configured"))

		err := runWithSpinner(styles.InfoIcon+" "+styles.Info.Render("Setting upstream to origin/HEAD"), func() error {
			pushCmd := exec.Command("git", "push", "--set-upstream", "origin", "HEAD")
			pushCmd.Stdout = os.Stdout
			pushCmd.Stderr = os.Stderr
			return pushCmd.Run()
		})

		if err != nil {
			fmt.Println(styles.ErrorIcon)
			return fmt.Errorf("failed to push and set upstream: %w", err)
		}
		fmt.Println(styles.SuccessIcon)

		// Show success message
		fmt.Println()
		fmt.Println(styles.Card.Render(
			styles.SuccessIcon + " " + styles.Success.Render("Push completed!") + "\n" +
				styles.Neutral.Render("Status: ") + styles.Success.Render("Upstream set and changes pushed"),
		))
		return nil
	}

	// Upstream exists, just push
	err := runWithSpinner(styles.InfoIcon+" "+styles.Info.Render("Pushing changes to remote..."), func() error {
		pushCmd := exec.Command("git", "push")
		pushCmd.Stdout = os.Stdout
		pushCmd.Stderr = os.Stderr
		return pushCmd.Run()
	})

	if err != nil {
		fmt.Println(styles.ErrorIcon)
		return fmt.Errorf("failed to push: %w", err)
	}
	fmt.Println(styles.SuccessIcon)

	// Show success message
	fmt.Println()
	fmt.Println(styles.Card.Render(
		styles.SuccessIcon + " " + styles.Success.Render("Push completed!") + "\n" +
			styles.Neutral.Render("Status: ") + styles.Success.Render("Changes pushed to remote"),
	))
	return nil
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
