package cmd

import (
	"fmt"

	"expenseVault/services"

	"github.com/spf13/cobra"
)

// ============================================================
// UNIT 3: Interfaces & Polymorphism — using Reporter interface
// UNIT 2: Maps for aggregation in report generation
// ============================================================

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate financial reports",
	Long: `Generate various financial reports: monthly, category-wise, or yearly.

Examples:
  expensevault report --type monthly
  expensevault report --type category
  expensevault report --type yearly
  expensevault report --all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		reportType, _ := cmd.Flags().GetString("type")
		showAll, _ := cmd.Flags().GetBool("all")

		transactions, err := store.GetAllTransactions()
		if err != nil {
			return fmt.Errorf("failed to get transactions: %w", err)
		}

		if len(transactions) == 0 {
			fmt.Println("\n📭 No transactions found. Add some transactions first!")
			return nil
		}

		// Demonstrates: Interface polymorphism (Unit 3)
		if showAll {
			// Generate all report types using polymorphic interface
			fmt.Print(services.GenerateAllReports(transactions))
		} else {
			// Get the appropriate reporter via factory function
			// Demonstrates: switch-based factory, interface usage (Unit 3)
			reporter := services.GetReporter(reportType)
			fmt.Print(services.GenerateReport(reporter, transactions))
		}

		return nil
	},
}

func init() {
	reportCmd.Flags().StringP("type", "t", "monthly", "Report type: monthly, category, yearly")
	reportCmd.Flags().Bool("all", false, "Generate all report types")
}
