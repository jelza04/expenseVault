package cmd

import (
	"fmt"

	"expenseVault/tui"

	"github.com/spf13/cobra"
)

// tuiAPICmd launches the TUI in API mode — communicates with the HTTP server.
var tuiAPICmd = &cobra.Command{
	Use:   "tui-api",
	Short: "Launch the TUI connected to the API server (requires server to be running)",
	RunE: func(cmd *cobra.Command, args []string) error {
		serverURL, _ := cmd.Flags().GetString("server")
		if serverURL == "" {
			serverURL = "http://localhost:8080"
		}
		fmt.Printf("Connecting TUI to API server at %s\n", serverURL)
		return tui.RunTUIWithAPI(serverURL)
	},
}

func init() {
	tuiAPICmd.Flags().String("server", "http://localhost:8080", "API server URL")
}
