package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ============================================================
// UNIT 3: Error handling — NotFoundError
// ============================================================

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a transaction by ID",
	Long: `Delete a transaction from ExpenseVault by its ID.

Examples:
  expensevault delete --id 1
  expensevault delete --id 5`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetInt("id")
		if id <= 0 {
			return fmt.Errorf("please provide a valid transaction ID with --id")
		}

		// Show the transaction before deleting
		tx, err := store.GetTransaction(id)
		if err != nil {
			return fmt.Errorf("failed to find transaction: %w", err)
		}

		fmt.Printf("\n⚠️  About to delete: %s\n", tx)

		// Demonstrates: Error handling (Unit 3)
		if err := store.DeleteTransaction(id); err != nil {
			return fmt.Errorf("failed to delete transaction: %w", err)
		}

		fmt.Printf("🗑️  Transaction #%d deleted successfully!\n\n", id)
		return nil
	},
}

func init() {
	deleteCmd.Flags().Int("id", 0, "Transaction ID to delete (required)")
	deleteCmd.MarkFlagRequired("id")
}
