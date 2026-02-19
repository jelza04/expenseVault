package cmd

import (
	"fmt"

	"expenseVault/models"

	"github.com/spf13/cobra"
)

// ============================================================
// UNIT 1: Control flow — conditionals for field updates
// UNIT 3: Error handling — NotFoundError, validation
// ============================================================

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an existing transaction",
	Long: `Edit a transaction by its ID. Only specified fields will be updated.

Examples:
  expensevault edit --id 1 --amount 300
  expensevault edit --id 2 --category Travel --notes "Updated"
  expensevault edit --id 3 --desc "New description" --date 2026-02-20`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetInt("id")
		if id <= 0 {
			return fmt.Errorf("please provide a valid transaction ID with --id")
		}

		// Fetch existing transaction
		// Demonstrates: Error handling with custom error types (Unit 3)
		tx, err := store.GetTransaction(id)
		if err != nil {
			return fmt.Errorf("failed to get transaction: %w", err)
		}

		// Update only the fields that were explicitly set
		// Demonstrates: Control flow — conditional field updates (Unit 1)
		if cmd.Flags().Changed("amount") {
			amount, _ := cmd.Flags().GetFloat64("amount")
			tx.Amount = models.Rupees(amount)
		}
		if cmd.Flags().Changed("category") {
			category, _ := cmd.Flags().GetString("category")
			tx.Category = models.Category(category)
		}
		if cmd.Flags().Changed("desc") {
			desc, _ := cmd.Flags().GetString("desc")
			tx.Description = desc
		}
		if cmd.Flags().Changed("date") {
			date, _ := cmd.Flags().GetString("date")
			tx.Date = date
		}
		if cmd.Flags().Changed("notes") {
			notes, _ := cmd.Flags().GetString("notes")
			tx.Notes = notes
		}
		if cmd.Flags().Changed("type") {
			txType, _ := cmd.Flags().GetString("type")
			tx.Type = models.TransactionType(txType)
		}

		// Demonstrates: Error checking (Unit 3)
		if err := store.UpdateTransaction(*tx); err != nil {
			return fmt.Errorf("failed to update transaction: %w", err)
		}

		fmt.Printf("\n✅ Transaction #%d updated successfully!\n", id)
		fmt.Println(tx)
		fmt.Println()

		return nil
	},
}

func init() {
	editCmd.Flags().Int("id", 0, "Transaction ID to edit (required)")
	editCmd.Flags().StringP("type", "t", "", "New type (income/expense)")
	editCmd.Flags().Float64P("amount", "a", 0, "New amount")
	editCmd.Flags().StringP("category", "c", "", "New category")
	editCmd.Flags().StringP("desc", "d", "", "New description")
	editCmd.Flags().String("date", "", "New date (YYYY-MM-DD)")
	editCmd.Flags().StringP("notes", "n", "", "New notes")
	editCmd.MarkFlagRequired("id")
}
