package analytics_test

import (
	"testing"
	"time"

	"expenseVault/analytics"
)

// ─────────────────────────────────────────────────────────────
// Test Helpers
// ─────────────────────────────────────────────────────────────

func mkTx(amount float64, dateStr, txType string) analytics.Transaction {
	t, _ := time.Parse("2006-01-02", dateStr)
	return analytics.Transaction{Amount: amount, Date: t, Type: txType, Category: "Test"}
}

var sampleTxs = []analytics.Transaction{
	mkTx(5000, "2026-01-05", "expense"),
	mkTx(3000, "2026-01-20", "expense"),
	mkTx(7000, "2026-02-10", "expense"),
	mkTx(2000, "2026-02-25", "expense"),
	mkTx(6000, "2026-03-15", "expense"),
	mkTx(10000, "2026-01-01", "income"), // Should be excluded from expense analytics
}

// ─────────────────────────────────────────────────────────────
// PART B.1: Spending Trends
// ─────────────────────────────────────────────────────────────

func TestComputeMonthlyTrends(t *testing.T) {
	tests := []struct {
		name     string
		txs      []analytics.Transaction
		expected map[string]float64
	}{
		{
			name:     "empty input",
			txs:      []analytics.Transaction{},
			expected: map[string]float64{},
		},
		{
			name:     "income only — no expense trends",
			txs:      []analytics.Transaction{mkTx(5000, "2026-01-01", "income")},
			expected: map[string]float64{},
		},
		{
			name: "single month",
			txs:  []analytics.Transaction{mkTx(3000, "2026-03-10", "expense"), mkTx(2000, "2026-03-20", "expense")},
			expected: map[string]float64{
				"2026-03": 5000,
			},
		},
		{
			name: "multi-month",
			txs:  sampleTxs,
			expected: map[string]float64{
				"2026-01": 8000,
				"2026-02": 9000,
				"2026-03": 6000,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := analytics.ComputeMonthlyTrends(tc.txs)
			if len(got) != len(tc.expected) {
				t.Errorf("got %d months, want %d", len(got), len(tc.expected))
			}
			for k, wantV := range tc.expected {
				if gotV, ok := got[k]; !ok || gotV != wantV {
					t.Errorf("month %s: got %.2f, want %.2f", k, gotV, wantV)
				}
			}
		})
	}
}

// ─────────────────────────────────────────────────────────────
// PART B.2: Average Spending
// ─────────────────────────────────────────────────────────────

func TestComputeAveragePerTransaction(t *testing.T) {
	tests := []struct {
		name     string
		txs      []analytics.Transaction
		expected float64
	}{
		{"empty", []analytics.Transaction{}, 0},
		{"income only", []analytics.Transaction{mkTx(5000, "2026-01-01", "income")}, 0},
		{"single expense", []analytics.Transaction{mkTx(3000, "2026-01-01", "expense")}, 3000},
		{"multi expenses", sampleTxs, 4600}, // (5000+3000+7000+2000+6000)/5
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := analytics.ComputeAveragePerTransaction(tc.txs)
			if got != tc.expected {
				t.Errorf("got %.2f, want %.2f", got, tc.expected)
			}
		})
	}
}

func TestComputeMonthlyAverage(t *testing.T) {
	tests := []struct {
		name     string
		trends   map[string]float64
		expected float64
	}{
		{"empty map", map[string]float64{}, 0},
		{"single month", map[string]float64{"2026-01": 8000}, 8000},
		{"multi month", map[string]float64{"2026-01": 8000, "2026-02": 9000, "2026-03": 6000}, 7666.666666666667},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := analytics.ComputeMonthlyAverage(tc.trends)
			// Use tolerance for float comparison
			if diff := got - tc.expected; diff > 0.01 || diff < -0.01 {
				t.Errorf("got %.6f, want %.6f", got, tc.expected)
			}
		})
	}
}

// ─────────────────────────────────────────────────────────────
// PART B.3: Prediction
// ─────────────────────────────────────────────────────────────

func TestPredictNextMonth(t *testing.T) {
	tests := []struct {
		name        string
		trends      map[string]float64
		expectAbove float64 // prediction should be > this value
		expectBelow float64 // prediction should be < this value
	}{
		{
			name:        "empty trends",
			trends:      map[string]float64{},
			expectAbove: -1,
			expectBelow: 1,
		},
		{
			name:        "single month",
			trends:      map[string]float64{"2026-01": 5000},
			expectAbove: 4999,
			expectBelow: 5001,
		},
		{
			name: "upward trend — prediction should be near/above recent month",
			trends: map[string]float64{
				"2026-01": 3000,
				"2026-02": 5000,
				"2026-03": 8000,
			},
			expectAbove: 5000, // weighted toward recent high
			expectBelow: 10000,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := analytics.PredictNextMonth(tc.trends)
			if got <= tc.expectAbove || got >= tc.expectBelow {
				t.Errorf("PredictNextMonth()=%.2f, expected in range (%.2f, %.2f)", got, tc.expectAbove, tc.expectBelow)
			}
		})
	}
}

// ─────────────────────────────────────────────────────────────
// PART C: Worker Pool
// ─────────────────────────────────────────────────────────────

func TestRunWorkerPool(t *testing.T) {
	// Build 2 jobs from sample data
	jobs := []analytics.Job{
		{ID: 1, Transactions: sampleTxs[:3]},
		{ID: 2, Transactions: sampleTxs[3:]},
	}

	results := analytics.RunWorkerPool(jobs, 2)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestRunWorkerPool_Empty(t *testing.T) {
	results := analytics.RunWorkerPool([]analytics.Job{}, 3)
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty job list, got %d", len(results))
	}
}

func TestMergeResults(t *testing.T) {
	results := []analytics.Result{
		{JobID: 1, Trends: map[string]float64{"2026-01": 5000, "2026-02": 3000}, Average: 4000},
		{JobID: 2, Trends: map[string]float64{"2026-02": 2000, "2026-03": 6000}, Average: 4000},
	}

	merged := analytics.MergeResults(results, sampleTxs)

	// 2026-02 should be 3000 + 2000 = 5000
	if got := merged.Trends["2026-02"]; got != 5000 {
		t.Errorf("merged 2026-02 = %.2f, want 5000", got)
	}
	if merged.Prediction == 0 {
		t.Error("merged prediction should not be 0")
	}
}

func TestChunkTransactions(t *testing.T) {
	chunks := analytics.ChunkTransactions(sampleTxs, 2)
	if len(chunks) != 3 { // 6 txs, chunks of 2 → 3 chunks
		t.Errorf("expected 3 chunks, got %d", len(chunks))
	}
	if len(chunks[0]) != 2 {
		t.Errorf("expected chunk size 2, got %d", len(chunks[0]))
	}
}

// ─────────────────────────────────────────────────────────────
// PART I: Benchmarks
// ─────────────────────────────────────────────────────────────

func BenchmarkComputeMonthlyTrends(b *testing.B) {
	// Generate 10,000 transactions for a realistic benchmark
	txs := make([]analytics.Transaction, 10000)
	base := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := range txs {
		txs[i] = analytics.Transaction{
			Amount:   float64(i%5000 + 500),
			Date:     base.AddDate(0, i%24, 0),
			Type:     "expense",
			Category: "Test",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analytics.ComputeMonthlyTrends(txs)
	}
}

func BenchmarkWorkerPool(b *testing.B) {
	txs := make([]analytics.Transaction, 1000)
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := range txs {
		txs[i] = analytics.Transaction{Amount: float64(i + 100), Date: base.AddDate(0, i%12, 0), Type: "expense"}
	}

	chunks := analytics.ChunkTransactions(txs, 100)
	jobs := make([]analytics.Job, len(chunks))
	for i, c := range chunks {
		jobs[i] = analytics.Job{ID: i, Transactions: c}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analytics.RunWorkerPool(jobs, 3)
	}
}
