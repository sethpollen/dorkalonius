package dorkalonius

import (
	"fmt"
	"time"
)

// Sleeps for 'duration', printing regular updates to stdout. 'duration'
// may be rounded up or down by 1 second.
func VerboseSleep(duration time.Duration, bell bool) {
	deadline := time.Now().Add(duration)
	remaining := duration

	// Print updates every second.
	var ticker = time.NewTicker(time.Second)
	for time.Now().Before(deadline) {
		// Home the cursor, clear to end of line, and print the remaining duration.
		fmt.Printf("\033[9999D\033[KSleeping for %s", remaining.String())
		remaining = remaining - time.Second

		// Await a tick.
		<-ticker.C
	}

	// Clear the remaining indicator before returning.
	fmt.Print("\033[9999D\033[K")

	// Print the final output line.
	fmt.Printf("Slept for %s\n", duration.String())

	if bell {
		fmt.Print("\007")
	}
}
