package cmd

import (
	"fmt"

	"expenseVault/tui"

	"github.com/spf13/cobra"
)

// tuiCmd launches the interactive BubbleTea TUI dashboard.
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive terminal dashboard",
	Long:  `Open the beautiful interactive TUI dashboard powered by BubbleTea.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🚀 Launching ExpenseVault TUI...")
		return tui.RunTUI(store)
	},
}
