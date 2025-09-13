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
