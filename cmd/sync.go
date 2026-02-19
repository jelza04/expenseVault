package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"expenseVault/models"

	"github.com/spf13/cobra"
)

// ============================================================
// Sync command — syncs local transactions with remote server
// Demonstrates: HTTP client, JSON serialization, error handling
// ============================================================

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync transactions with the remote server",
	Long: `Sync local transactions with a self-hosted ExpenseVault server.

Examples:
  expensevault sync --server http://localhost:8080 --token <jwt-token>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		serverURL, _ := cmd.Flags().GetString("server")
		token, _ := cmd.Flags().GetString("token")

		if serverURL == "" || token == "" {
			return fmt.Errorf("both --server and --token are required")
		}

		fmt.Println("🔄 Syncing with server...")

		// Get local transactions
		transactions, err := store.GetAllTransactions()
		if err != nil {
			return fmt.Errorf("failed to get local transactions: %w", err)
		}

		// Build sync payload
		payload := models.SyncPayload{
			Transactions: transactions,
		}

		data, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal sync data: %w", err)
		}

		// Send to server
		req, err := http.NewRequest("POST", serverURL+"/api/sync", bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("sync failed: %w", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body))
		}

		var result struct {
			Synced       int                  `json:"synced"`
			Message      string               `json:"message"`
			Transactions []models.Transaction `json:"transactions"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}

		// Import server transactions locally
		if len(result.Transactions) > 0 {
			count, _ := store.BulkInsert(result.Transactions)
			fmt.Printf("   📥 Imported %d transactions from server\n", count)
		}

		fmt.Printf("   📤 Pushed %d transactions to server\n", result.Synced)
		fmt.Println("✅ Sync complete!")
		return nil
	},
}

func init() {
	syncCmd.Flags().String("server", "http://localhost:8080", "Server URL")
	syncCmd.Flags().String("token", "", "JWT auth token")
}
