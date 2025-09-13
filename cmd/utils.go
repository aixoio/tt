package cmd

import (
	"fmt"
	"sync"
	"time"
)

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

// runWithSpinnerForMessage runs a function with a spinner and returns a string result
func runWithSpinnerForMessage(title string, action func() (string, error)) (string, error) {
	fmt.Print("\033[?25l") // hide cursor
	stop := make(chan struct{})
	var wg sync.WaitGroup
	var result string
	var err error
	var once sync.Once

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

	// Run the action in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		var innerErr error
		result, innerErr = action()
		if innerErr != nil {
			once.Do(func() {
				err = innerErr
			})
		}
		close(stop)
	}()

	wg.Wait()
	fmt.Println()
	fmt.Print("\033[?25h") // show cursor
	return result, err
}
