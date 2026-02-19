package cmd

import (
	"fmt"

	"expenseVault/models"
	"expenseVault/utils"

	"github.com/spf13/cobra"
)

// ============================================================
// UNIT 1: Looping structures — for range
// UNIT 2: Slice operations — filtering, slicing
// ============================================================

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List transactions with optional filters",
	Long: `List all transactions or filter by type, category, date range, etc.

Examples:
  expensevault list
  expensevault list --type expense
  expensevault list --category Food --limit 10
  expensevault list --from 2026-01-01 --to 2026-01-31`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Demonstrates: Short declaration operator (Unit 1)
		txType, _ := cmd.Flags().GetString("type")
		category, _ := cmd.Flags().GetString("category")
		startDate, _ := cmd.Flags().GetString("from")
		endDate, _ := cmd.Flags().GetString("to")
		limit, _ := cmd.Flags().GetInt("limit")

		// Demonstrates: Error checking (Unit 3)
		transactions, err := store.ListTransactions(txType, category, startDate, endDate, limit)
		if err != nil {
			return fmt.Errorf("failed to list transactions: %w", err)
		}

		// Demonstrates: Control flow — conditional (Unit 1)
		if len(transactions) == 0 {
			fmt.Println("\n📭 No transactions found.")
			fmt.Println("   Use 'expensevault add' to add your first transaction.")
			return nil
		}

		// Build table headers and rows
		headers := []string{"ID", "Type", "Amount", "Category", "Description", "Date", "Notes"}

		// Demonstrates: make for pre-allocation (Unit 2)
		rows := make([][]string, 0, len(transactions))

		// Demonstrates: for range over slice (Unit 1 & Unit 2)
		for _, tx := range transactions {
			symbol := "💰"
			if tx.Type == models.Expense {
				symbol = "💸"
			}
			rows = append(rows, []string{
				fmt.Sprintf("%d", tx.ID),
				fmt.Sprintf("%s %s", symbol, tx.Type),
				tx.Amount.String(),
				string(tx.Category),
				utils.TruncateString(tx.Description, 25),
				tx.Date,
				utils.TruncateString(tx.Notes, 20),
			})
		}

		fmt.Printf("\n📋 Transactions (%d found)\n\n", len(transactions))
		utils.PrintTable(headers, rows)

		// Show quick summary using anonymous struct
		// Demonstrates: Anonymous struct usage (Unit 2)
		summary := models.QuickSummary(transactions)
		fmt.Printf("\n📊 Summary: Income: %s | Expenses: %s | Balance: %s\n\n",
			summary.TotalIncome, summary.TotalExpense, summary.Balance)

		return nil
	},
}

func init() {
	listCmd.Flags().StringP("type", "t", "", "Filter by type (income/expense)")
	listCmd.Flags().StringP("category", "c", "", "Filter by category")
	listCmd.Flags().String("from", "", "Start date (YYYY-MM-DD)")
	listCmd.Flags().String("to", "", "End date (YYYY-MM-DD)")
	listCmd.Flags().IntP("limit", "l", 0, "Limit number of results")
}
