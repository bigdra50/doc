package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bigdra50/doc/internal/utils"
)

var verbose bool

// log outputs debug messages when verbose mode is enabled
func log(format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stderr, "[DEBUG] %s: %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(format, args...))
	}
}

// progress outputs informational messages
func progress(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[INFO] %s\n", fmt.Sprintf(format, args...))
}

// Spinner represents a loading spinner with elapsed time display
type Spinner struct {
	message   string
	frames    []string
	interval  time.Duration
	startTime time.Time
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

// NewSpinner creates a new spinner with the given message
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message:  message,
		frames:   []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		interval: 100 * time.Millisecond,
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	if !isTerminal() {
		fmt.Fprintf(os.Stderr, "[INFO] %s\n", s.message)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.startTime = time.Now()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		frame := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(s.interval):
				elapsed := time.Since(s.startTime)
				fmt.Fprintf(os.Stderr, "\r%s %s (%s)", s.frames[frame], s.message, formatDuration(elapsed))
				frame = (frame + 1) % len(s.frames)
			}
		}
	}()
}

// Stop ends the spinner animation and displays a final message
func (s *Spinner) Stop(finalMessage string) {
	if s.cancel == nil {
		return
	}

	s.cancel()
	s.wg.Wait()

	if isTerminal() {
		elapsed := time.Since(s.startTime)
		fmt.Fprintf(os.Stderr, "\r✓ %s (%s)\n", finalMessage, formatDuration(elapsed))
	} else {
		fmt.Fprintf(os.Stderr, "[INFO] %s\n", finalMessage)
	}
}

// isTerminal checks if stderr is connected to a terminal
func isTerminal() bool {
	fileInfo, _ := os.Stderr.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// formatDuration formats a duration for display
func formatDuration(d time.Duration) string {
	return utils.FormatDuration(d)
}

// maskAPIKey masks an API key for safe logging
func maskAPIKey(key string) string {
	return utils.MaskAPIKey(key)
}
