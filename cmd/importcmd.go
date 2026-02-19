package cmd

import (
	"fmt"

	"expenseVault/export"

	"github.com/spf13/cobra"
)

// ============================================================
// CSV Import using encoding/csv
// Demonstrates: Exporter/Importer interface (Unit 3), error handling
// ============================================================

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import transactions from a CSV file",
	Long: `Import transactions from a CSV file (e.g., bank statement export).

CSV format: ID,Type,Amount,Category,Description,Date,Notes

Examples:
  expensevault import --file bank_statement.csv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")
		if filePath == "" {
			return fmt.Errorf("please provide a CSV file with --file")
		}

		// Demonstrates: Interface usage — CSVHandler implements Importer (Unit 3)
		handler := &export.CSVHandler{}
		transactions, err := handler.Import(filePath)
		if err != nil {
			return fmt.Errorf("failed to import CSV: %w", err)
		}

		if len(transactions) == 0 {
			fmt.Println("📭 No transactions found in the CSV file.")
			return nil
		}

		// Demonstrates: Variadic usage via BulkInsert (Unit 3)
		count, err := store.BulkInsert(transactions)
		if err != nil {
			return fmt.Errorf("failed to insert transactions: %w", err)
		}

		fmt.Printf("\n✅ Successfully imported %d transactions from %s\n\n", count, filePath)
		return nil
	},
}

func init() {
	importCmd.Flags().StringP("file", "f", "", "CSV file to import (required)")
	importCmd.MarkFlagRequired("file")
}
