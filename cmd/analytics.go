package cmd

import (
	"fmt"
	"sort"

	"expenseVault/analytics"

	"github.com/spf13/cobra"
)

var analyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "Run the Financial Analytics Engine (trends, averages, predictions)",
	Long: `Analyzes your transaction history using a concurrent worker pool.

Computes:
  • Monthly spending trends
  • Average per-transaction and monthly spending
  • Weighted prediction for next month's spending`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workers, _ := cmd.Flags().GetInt("workers")
		chunk, _ := cmd.Flags().GetInt("chunk-size")

		// Retrieve all transactions for the logged-in user
		userID, err := getCurrentUserID()
		if err != nil {
			return fmt.Errorf("please log in first: %w", err)
		}

		txs, err := store.GetAllTransactions(userID)
		if err != nil {
			return fmt.Errorf("failed to load transactions: %w", err)
		}

		// Configure and run the analytics engine
		svc := analytics.NewAnalyticsService()
		if workers > 0 {
			svc.NumWorkers = workers
		}
		if chunk > 0 {
			svc.ChunkSize = chunk
		}

		fmt.Printf("\n🔍 Running analytics with %d workers (chunk size: %d)...\n\n", svc.NumWorkers, svc.ChunkSize)

		report, err := svc.Run(txs)
		if err != nil {
			return fmt.Errorf("analytics error: %w", err)
		}

		// ── Display Results ─────────────────────────────────────
		fmt.Println("┌─────────────────────────────────────────────┐")
		fmt.Println("│          📊 Financial Analytics Report       │")
		fmt.Println("├─────────────────────────────────────────────┤")

		// Spending Trends
		fmt.Println("│  📅 Monthly Spending Trends                  │")
		fmt.Println("├─────────────────────────────────────────────┤")

		keys := sortedKeys(report.Trends)
		for _, k := range keys {
			bar := progressBar(report.Trends[k], maxVal(report.Trends))
			fmt.Printf("│  %s  %s ₹%-10.2f            │\n", k, bar, report.Trends[k])
		}

		fmt.Println("├─────────────────────────────────────────────┤")
		fmt.Printf("│  💰 Avg Per Transaction : ₹%-10.2f        │\n", report.Average)
		fmt.Printf("│  📆 Avg Monthly Spend   : ₹%-10.2f        │\n", report.MonthlyAvg)
		fmt.Printf("│  🔮 Predicted Next Month: ₹%-10.2f        │\n", report.Prediction)
		fmt.Println("├─────────────────────────────────────────────┤")
		fmt.Printf("│  ⏱  Generated: %-30s │\n", report.GeneratedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("└─────────────────────────────────────────────┘")

		return nil
	},
}

func init() {
	analyticsCmd.Flags().Int("workers", 3, "Number of worker goroutines (default: 3)")
	analyticsCmd.Flags().Int("chunk-size", 50, "Transactions per worker chunk (default: 50)")
}

// ─────────────────────────────────────────────────────────────
// Display helpers
// ─────────────────────────────────────────────────────────────

func sortedKeys(m map[string]float64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func maxVal(m map[string]float64) float64 {
	var max float64
	for _, v := range m {
		if v > max {
			max = v
		}
	}
	return max
}

func progressBar(val, max float64) string {
	const barWidth = 10
	if max == 0 {
		return "[          ]"
	}
	filled := int((val / max) * barWidth)
	bar := "["
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += " "
		}
	}
	return bar + "]"
}
