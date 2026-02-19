package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ============================================================
// Utility helpers used throughout ExpenseVault
// ============================================================

// GetDBPath returns the path to the SQLite database file.
func GetDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "expensevault.db"
	}
	dir := filepath.Join(home, ".expensevault")
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "expensevault.db")
}

// GetServerDBPath returns the DB path for the sync server.
func GetServerDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "expensevault_server.db"
	}
	dir := filepath.Join(home, ".expensevault")
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "server.db")
}

// FormatDate formats a time.Time to YYYY-MM-DD string.
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// Today returns today's date as YYYY-MM-DD.
func Today() string {
	return FormatDate(time.Now())
}

// ParseDate parses a YYYY-MM-DD string to time.Time.
func ParseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// TruncateString truncates a string to maxLen and adds "..." if needed.
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// PadRight pads a string with spaces to reach the desired width.
func PadRight(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}

// PrintTable prints a formatted ASCII table.
func PrintTable(headers []string, rows [][]string) {
	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Build separator
	sep := "+"
	for _, w := range widths {
		sep += strings.Repeat("-", w+2) + "+"
	}

	// Print header
	fmt.Println(sep)
	fmt.Print("|")
	for i, h := range headers {
		fmt.Printf(" %-*s |", widths[i], h)
	}
	fmt.Println()
	fmt.Println(sep)

	// Print rows
	for _, row := range rows {
		fmt.Print("|")
		for i, cell := range row {
			if i < len(widths) {
				fmt.Printf(" %-*s |", widths[i], cell)
			}
		}
		fmt.Println()
	}
	fmt.Println(sep)
}

// ColorText returns ANSI-colored text for terminal output.
func ColorText(text string, colorCode string) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", colorCode, text)
}

// Color codes for terminal
const (
	ColorRed    = "31"
	ColorGreen  = "32"
	ColorYellow = "33"
	ColorBlue   = "34"
	ColorCyan   = "36"
	ColorBold   = "1"
)
