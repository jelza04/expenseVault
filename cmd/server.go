package cmd

import (
	"fmt"

	"expenseVault/api"
	"expenseVault/db"
	"expenseVault/utils"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the ExpenseVault HTTP API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := utils.LoadConfig()
		if err != nil {
			return err
		}

		// Initialize the database store for the server
		var serverStore *db.Store
		if cfg.DBType == "mysql" {
			serverStore, err = db.NewMySQLStore(cfg.MySQLDSN)
		} else {
			serverStore, err = db.NewStore(cfg.SQLitePath)
		}
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer serverStore.Close()

		port, _ := cmd.Flags().GetInt("port")
		if port == 0 {
			port = cfg.ServerPort
		}
		addr := fmt.Sprintf(":%d", port)
		fmt.Printf("Starting ExpenseVault API server on %s\n", addr)
		fmt.Println("Routes:")
		fmt.Println("  POST /api/signup         — Create account")
		fmt.Println("  POST /api/login          — Login & get JWT")
		fmt.Println("  GET  /api/transactions   — List transactions (auth)")
		fmt.Println("  POST /api/transactions   — Add transaction (auth)")
		fmt.Println("  POST /api/ask            — AI query (auth)")
		fmt.Println("  GET  /health             — Health check")

		return api.StartServer(addr, serverStore, cfg)
	},
}

func init() {
	serverCmd.Flags().Int("port", 0, "Server port")
}
