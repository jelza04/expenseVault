package analytics

import (
	"fmt"
	"time"

	"expenseVault/models"
)

const (
	DefaultNumWorkers = 3
	DefaultChunkSize  = 50 // transactions per worker chunk
)

// AnalyticsService orchestrates the full analytics pipeline:
//  1. Fetch transactions from the DB store
//  2. Convert to local analytics.Transaction type
//  3. Chunk data and dispatch to the worker pool
//  4. Merge partial results into a final AnalyticsReport
type AnalyticsService struct {
	NumWorkers int
	ChunkSize  int
}

// NewAnalyticsService creates an AnalyticsService with sensible defaults.
func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{
		NumWorkers: DefaultNumWorkers,
		ChunkSize:  DefaultChunkSize,
	}
}

// AnalyticsReport is the final aggregated output of the analytics pipeline.
type AnalyticsReport struct {
	Trends      map[string]float64
	Average     float64
	MonthlyAvg  float64
	Prediction  float64
	GeneratedAt time.Time
}

// Run executes the analytics pipeline on the given model transactions.
// Returns an error if the list is empty or all transactions are income-only.
func (s *AnalyticsService) Run(modelTxs []models.Transaction) (*AnalyticsReport, error) {
	if len(modelTxs) == 0 {
		return nil, fmt.Errorf("no transactions found — add some expenses first")
	}

	// Step 1: Convert models.Transaction → analytics.Transaction
	txs := ConvertFromModels(modelTxs)

	// Step 2: Check for expense data
	var hasExpense bool
	for _, tx := range txs {
		if tx.Type == "expense" {
			hasExpense = true
			break
		}
	}
	if !hasExpense {
		return nil, fmt.Errorf("no expense transactions found — analytics requires expense data")
	}

	// Step 3: Split into chunks for parallel processing
	chunks := ChunkTransactions(txs, s.ChunkSize)

	jobs := make([]Job, len(chunks))
	for i, chunk := range chunks {
		jobs[i] = Job{ID: i + 1, Transactions: chunk}
	}

	// Step 4: Run the worker pool
	results := RunWorkerPool(jobs, s.NumWorkers)

	// Step 5: Merge partial results
	merged := MergeResults(results, txs)
	merged.Average = ComputeAveragePerTransaction(txs)

	return &AnalyticsReport{
		Trends:      merged.Trends,
		Average:     merged.Average,
		MonthlyAvg:  ComputeMonthlyAverage(merged.Trends),
		Prediction:  merged.Prediction,
		GeneratedAt: time.Now(),
	}, nil
}

// ConvertFromModels converts a slice of models.Transaction to []analytics.Transaction.
func ConvertFromModels(modelTxs []models.Transaction) []Transaction {
	txs := make([]Transaction, 0, len(modelTxs))
	for _, m := range modelTxs {
		t, err := time.Parse("2006-01-02", m.Date)
		if err != nil {
			t = time.Now()
		}
		txs = append(txs, Transaction{
			Amount:   float64(m.Amount),
			Date:     t,
			Category: string(m.Category),
			Type:     string(m.Type),
		})
	}
	return txs
}
