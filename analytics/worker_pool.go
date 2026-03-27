package analytics

import (
	"log"
	"sync"
)

// Job is a unit of work sent to a worker.
// Each job contains a slice of transactions to process.
type Job struct {
	ID           int
	Transactions []Transaction
}

// Result holds the output of a single worker's computation.
type Result struct {
	JobID      int
	Trends     map[string]float64
	Average    float64
	Prediction float64
}

// ─────────────────────────────────────────────────────────────
//  WHY WORKER POOL instead of raw goroutines?
//  - Unbounded goroutines (go func() for each item) can exhaust
//    memory and cause scheduler thrashing on large datasets.
//  - A worker pool caps concurrency at N workers, ensuring
//    predictable memory usage and CPU utilization.
//  - Channels provide safe, race-free communication — no mutexes
//    needed for data exchange between goroutines.
// ─────────────────────────────────────────────────────────────

// RunWorkerPool processes a slice of jobs concurrently using a
// fixed-size pool of numWorkers goroutines.
//
// Concurrency Model:
//   - jobCh:    buffered input channel (producer → workers)
//   - resultCh: buffered output channel (workers → consumer)
//   - WaitGroup: ensures all workers finish before closing resultCh
//
// Race Condition Prevention:
//   - Each worker operates on its own Job copy (no shared mutable state)
//   - All communication happens through channels (CSP model)
//   - WaitGroup prevents premature resultCh close
func RunWorkerPool(jobs []Job, numWorkers int) []Result {
	if numWorkers <= 0 {
		numWorkers = 3
	}

	jobCh := make(chan Job, len(jobs))
	resultCh := make(chan Result, len(jobs))

	var wg sync.WaitGroup

	// ── Launch the worker pool ───────────────────────────────
	for id := 1; id <= numWorkers; id++ {
		wg.Add(1)
		workerID := id
		go func() {
			defer wg.Done()
			worker(workerID, jobCh, resultCh)
		}()
	}

	// ── Send jobs to workers ─────────────────────────────────
	for _, job := range jobs {
		jobCh <- job
	}
	close(jobCh) // Signal workers: no more jobs

	// ── Close result channel once all workers finish ─────────
	// This runs in a separate goroutine so the main thread can
	// read from resultCh without deadlocking.
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// ── Collect all results ──────────────────────────────────
	var results []Result
	for r := range resultCh {
		results = append(results, r)
	}
	return results
}

// worker is the goroutine function: reads from jobCh, computes analytics,
// sends Result to resultCh.
func worker(id int, jobCh <-chan Job, resultCh chan<- Result) {
	for job := range jobCh {
		log.Printf("[Worker %d] processing Job #%d (%d transactions)", id, job.ID, len(job.Transactions))

		trends := ComputeMonthlyTrends(job.Transactions)
		avg := ComputeAveragePerTransaction(job.Transactions)
		pred := PredictNextMonth(trends)

		resultCh <- Result{
			JobID:      job.ID,
			Trends:     trends,
			Average:    avg,
			Prediction: pred,
		}

		log.Printf("[Worker %d] Job #%d done — avg=%.2f prediction=%.2f", id, job.ID, avg, pred)
	}
	log.Printf("[Worker %d] shutting down", id)
}

// ChunkTransactions splits a flat slice into chunks of chunkSize.
// This enables parallel processing: each chunk becomes one Job.
func ChunkTransactions(transactions []Transaction, chunkSize int) [][]Transaction {
	if chunkSize <= 0 {
		chunkSize = 1
	}
	var chunks [][]Transaction
	for i := 0; i < len(transactions); i += chunkSize {
		end := i + chunkSize
		if end > len(transactions) {
			end = len(transactions)
		}
		chunks = append(chunks, transactions[i:end])
	}
	return chunks
}

// MergeResults combines partial Results from multiple workers into
// a single unified Result.
//
// Merging strategy:
//   - Trends: add values for the same month key across all results
//   - Average: weighted mean of partial averages (by transaction count)
//   - Prediction: re-run PredictNextMonth on the merged trends
func MergeResults(results []Result, allTransactions []Transaction) Result {
	merged := Result{
		Trends: make(map[string]float64),
	}

	// Merge trend maps — safe because all workers write to separate Result copies
	for _, r := range results {
		for month, amount := range r.Trends {
			merged.Trends[month] += amount
		}
	}

	// Recompute aggregate stats on the full dataset for accuracy
	merged.Average = ComputeAveragePerTransaction(allTransactions)
	merged.Prediction = PredictNextMonth(merged.Trends)
	return merged
}
