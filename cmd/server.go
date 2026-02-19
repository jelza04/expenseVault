package cmd

import (
	"expenseVault/api"
	"expenseVault/db"
	"expenseVault/utils"

	"github.com/spf13/cobra"
)

// serverCmd starts the REST API server for multi-device sync.
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the self-hosted REST API server for sync",
	Long: `Start the ExpenseVault backend server with JWT authentication
for syncing transactions across multiple devices.

Examples:
  expensevault server
  expensevault server --port 9090
  expensevault server --secret my-jwt-secret`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetInt("port")
		secret, _ := cmd.Flags().GetString("secret")

		// Server uses its own DB (separate from client)
		serverStore, err := db.NewStore(utils.GetServerDBPath())
		if err != nil {
			return err
		}
		defer serverStore.Close()

		srv := api.NewServer(serverStore, secret)
		return srv.Start(port)
	},
}

func init() {
	serverCmd.Flags().IntP("port", "p", 8080, "Server port")
	serverCmd.Flags().StringP("secret", "s", "expensevault-secret-key", "JWT secret key")
}
