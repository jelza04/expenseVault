package cmd

import (
	"fmt"

	"expenseVault/export"

	"github.com/spf13/cobra"
)

// ============================================================
// Export to CSV/JSON using Exporter interface
// Demonstrates: Interface polymorphism (Unit 3) — choose format at runtime
// ============================================================

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export transactions to CSV or JSON",
	Long: `Export all transactions to a file in CSV or JSON format.

Examples:
  expensevault export --format csv --output transactions.csv
  expensevault export --format json --output data.json
  expensevault export -f csv -o my_data.csv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		output, _ := cmd.Flags().GetString("output")

		transactions, err := store.GetAllTransactions()
		if err != nil {
			return fmt.Errorf("failed to get transactions: %w", err)
		}

		if len(transactions) == 0 {
			fmt.Println("📭 No transactions to export.")
			return nil
		}

		// Demonstrates: Interface polymorphism — choose implementation at runtime (Unit 3)
		var exporter export.Exporter

		// Demonstrates: switch control flow (Unit 1)
		switch format {
		case "csv":
			exporter = &export.CSVHandler{}
		case "json":
			exporter = &export.JSONHandler{}
		default:
			return fmt.Errorf("unsupported format '%s'. Use 'csv' or 'json'", format)
		}

		// Default output filename
		if output == "" {
			output = fmt.Sprintf("expensevault_export.%s", format)
		}

		// Demonstrates: Interface method call — works with any Exporter (Unit 3)
		if err := exporter.Export(transactions, output); err != nil {
			return fmt.Errorf("failed to export: %w", err)
		}

		fmt.Printf("\n✅ Exported %d transactions to %s (%s format)\n\n",
			len(transactions), output, exporter.Format())
		return nil
	},
}

func init() {
	exportCmd.Flags().StringP("format", "f", "csv", "Export format: csv or json")
	exportCmd.Flags().StringP("output", "o", "", "Output file path")
}
