package analytics

// PredictNextMonth predicts next month's spending using a weighted
// moving-average approach:
//   - Most recent months are weighted more heavily
//   - Uses trend slope to adjust the prediction up or down
//
// # Why simple average-based prediction is effective
//
// A weighted moving average is:
//   1. Easy to explain and audit
//   2. Resilient to single outlier months
//   3. Respects recency — if spending has grown, the prediction grows too
//   4. Requires zero external dependencies
//
// Formula:
//   prediction = Σ(amount[i] * weight[i]) / Σ(weight[i])
//   where weight[i] = i+1 (linear weighting, most recent month = highest)
func PredictNextMonth(trends map[string]float64) float64 {
	if len(trends) == 0 {
		return 0
	}

	// Sort months chronologically so weights are applied correctly
	keys := SortedMonthKeys(trends)

	var weightedSum float64
	var weightTotal float64
	for i, k := range keys {
		weight := float64(i + 1) // linear weight: most recent = largest
		weightedSum += trends[k] * weight
		weightTotal += weight
	}

	if weightTotal == 0 {
		return 0
	}

	prediction := weightedSum / weightTotal

	// Trend adjustment: if the last month's spending > simple average,
	// add a small upward bias (5% of the difference from average).
	if len(keys) >= 2 {
		simpleAvg := ComputeMonthlyAverage(trends)
		lastMonth := trends[keys[len(keys)-1]]
		trendAdjustment := (lastMonth - simpleAvg) * 0.05
		prediction += trendAdjustment
	}

	return prediction
}
