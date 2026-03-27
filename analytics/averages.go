package analytics

// ComputeAveragePerTransaction calculates the mean expense amount
// across all individual expense transactions.
//
// Returns 0 if there are no expense transactions (safe for empty input).
func ComputeAveragePerTransaction(transactions []Transaction) float64 {
	var total float64
	var count int
	for _, tx := range transactions {
		if tx.Type != "expense" {
			continue
		}
		total += tx.Amount
		count++
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

// ComputeMonthlyAverage calculates the mean monthly spending
// from a pre-computed trends map.
//
// Returns 0 if the trends map is empty.
func ComputeMonthlyAverage(trends map[string]float64) float64 {
	if len(trends) == 0 {
		return 0
	}
	var total float64
	for _, v := range trends {
		total += v
	}
	return total / float64(len(trends))
}
