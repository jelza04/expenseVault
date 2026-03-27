package analytics

import (
	"sort"
	"time"
)

// Transaction is the local analytics representation.
// Mirrors models.Transaction but is self-contained for the analytics module.
type Transaction struct {
	Amount   float64
	Date     time.Time
	Category string
	Type     string // "income" or "expense"
}

// ComputeMonthlyTrends groups EXPENSE transactions by month (YYYY-MM)
// and returns the total spending per month.
//
// Design Note: We use map[string]float64 because:
//   - O(1) lookup by month key
//   - Month string (YYYY-MM) is a natural, human-readable key
//   - Easy to merge partial results from multiple workers
func ComputeMonthlyTrends(transactions []Transaction) map[string]float64 {
	trends := make(map[string]float64)
	for _, tx := range transactions {
		if tx.Type != "expense" {
			continue
		}
		key := tx.Date.Format("2006-01")
		trends[key] += tx.Amount
	}
	return trends
}

// SortedMonthKeys returns trend keys in ascending chronological order.
// Useful for display and prediction.
func SortedMonthKeys(trends map[string]float64) []string {
	keys := make([]string, 0, len(trends))
	for k := range trends {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
