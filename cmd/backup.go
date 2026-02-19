package cmd

import (
	"fmt"

	"expenseVault/export"

	"github.com/spf13/cobra"
)

// ============================================================
// Backup & Restore — JSON serialization
// ============================================================

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup all transactions to a JSON file",
	Long: `Create a full backup of all transactions in JSON format.

Examples:
  expensevault backup
  expensevault backup --output my_backup.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = "expensevault_backup.json"
		}

		transactions, err := store.GetAllTransactions()
		if err != nil {
			return fmt.Errorf("failed to get transactions: %w", err)
		}

		// Demonstrates: Using Exporter interface (Unit 3)
		handler := &export.JSONHandler{}
		if err := handler.Export(transactions, output); err != nil {
			return fmt.Errorf("backup failed: %w", err)
		}

		fmt.Printf("\n💾 Backup complete! %d transactions saved to %s\n\n", len(transactions), output)
		return nil
	},
}

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore transactions from a JSON backup file",
	Long: `Restore transactions from a previously created JSON backup.

Examples:
  expensevault restore --file expensevault_backup.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")
		if filePath == "" {
			return fmt.Errorf("please provide a backup file with --file")
		}

		// Demonstrates: Importer interface (Unit 3)
		handler := &export.JSONHandler{}
		transactions, err := handler.Import(filePath)
		if err != nil {
			return fmt.Errorf("restore failed: %w", err)
		}

		// Demonstrates: BulkInsert with slice (Unit 2 + Unit 3)
		count, err := store.BulkInsert(transactions)
		if err != nil {
			return fmt.Errorf("failed to restore transactions: %w", err)
		}

		fmt.Printf("\n✅ Restored %d transactions from %s\n\n", count, filePath)
		return nil
	},
}

func init() {
	backupCmd.Flags().StringP("output", "o", "", "Backup file path (default: expensevault_backup.json)")
	restoreCmd.Flags().StringP("file", "f", "", "Backup file to restore (required)")
	restoreCmd.MarkFlagRequired("file")
}
