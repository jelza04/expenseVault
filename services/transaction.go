package services

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"expenseVault/db"
	"expenseVault/models"
)

// ============================================================
// UNIT 3: Functions — Variadic Parameters, Closures, Callbacks,
//         Returning Functions, Function Expressions, Recursion
// UNIT 2: Slice operations — append, delete, filter, make,
//         multi-dimensional slices
// ============================================================

// TransactionService handles business logic for transactions.
type TransactionService struct {
	store *db.Store
}

// NewTransactionService creates a new service instance.
func NewTransactionService(store *db.Store) *TransactionService {
	return &TransactionService{store: store}
}

// ============================================================
// UNIT 3: Variadic Parameters — Accept multiple transactions
// ============================================================

// AddMultiple accepts a variadic number of transactions and adds them all.
// Demonstrates: Variadic parameter (Unit 3)
func (s *TransactionService) AddMultiple(transactions ...models.Transaction) ([]int64, error) {
	ids := make([]int64, 0, len(transactions))
	for _, t := range transactions {
		id, err := s.store.AddTransaction(t)
		if err != nil {
			return ids, fmt.Errorf("failed to add transaction '%s': %w", t.Description, err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// ============================================================
// UNIT 3: Unfurling a Slice — pass slice as variadic args
// ============================================================

// AddFromSlice demonstrates unfurling a slice into variadic params.
// Demonstrates: Unfurling a slice with ... (Unit 3)
func (s *TransactionService) AddFromSlice(transactions []models.Transaction) ([]int64, error) {
	// Unfurl the slice into variadic function call
	return s.AddMultiple(transactions...)
}

// ============================================================
// UNIT 2: Slice Operations — Delete from slice, slicing a slice
// ============================================================

// FilterTransactions retrieves transactions and applies in-memory filters.
// Demonstrates: Slice operations, for range, slicing a slice (Unit 2)
func (s *TransactionService) FilterTransactions(
	txType string, category string, minAmount float64, maxAmount float64,
	keyword string, startDate string, endDate string,
) ([]models.Transaction, error) {
	// Get from DB with basic filters
	transactions, err := s.store.ListTransactions(txType, category, startDate, endDate, 0)
	if err != nil {
		return nil, err
	}

	// Demonstrates: make to create result slice (Unit 2)
	result := make([]models.Transaction, 0, len(transactions))

	// Demonstrates: for range over slice (Unit 2)
	for _, tx := range transactions {
		// Demonstrates: Control flow — conditional filtering (Unit 1)
		if minAmount > 0 && tx.Amount.ToFloat64() < minAmount {
			continue
		}
		if maxAmount > 0 && tx.Amount.ToFloat64() > maxAmount {
			continue
		}
		if keyword != "" {
			lower := strings.ToLower(keyword)
			if !strings.Contains(strings.ToLower(tx.Description), lower) &&
				!strings.Contains(strings.ToLower(tx.Notes), lower) {
				continue
			}
		}
		result = append(result, tx)
	}

	return result, nil
}

// DeleteAndCompact demonstrates deleting from a slice and re-indexing.
// Demonstrates: Delete from a slice (Unit 2)
func DeleteFromSlice(transactions []models.Transaction, index int) []models.Transaction {
	if index < 0 || index >= len(transactions) {
		return transactions
	}
	// Demonstrates: Delete from slice by re-slicing (Unit 2)
	// transactions = append(transactions[:i], transactions[i+1:]...)
	return append(transactions[:index], transactions[index+1:]...)
}

// ============================================================
// UNIT 2: Multi-dimensional Slices
// ============================================================

// GroupByMonth groups transactions into a 2D slice by month.
// Demonstrates: Multi-dimensional slice [][]models.Transaction (Unit 2)
func GroupByMonth(transactions []models.Transaction) [][]models.Transaction {
	// Demonstrates: Map for grouping (Unit 2)
	monthMap := make(map[string][]models.Transaction)

	for _, tx := range transactions {
		key := tx.MonthYear()
		monthMap[key] = append(monthMap[key], tx)
	}

	// Convert map to 2D slice
	// Demonstrates: Multi-dimensional slice (Unit 2)
	result := make([][]models.Transaction, 0, len(monthMap))
	for _, txs := range monthMap {
		result = append(result, txs)
	}

	return result
}

// MonthlyTotals returns a 2D slice of [month, income, expense] values.
// Demonstrates: Multi-dimensional slice [][]float64 (Unit 2)
func MonthlyTotals(transactions []models.Transaction) ([]string, [][]float64) {
	type monthData struct {
		income  float64
		expense float64
	}

	monthMap := make(map[string]*monthData)
	monthOrder := make([]string, 0)

	for _, tx := range transactions {
		key := tx.MonthYear()
		if _, exists := monthMap[key]; !exists {
			monthMap[key] = &monthData{}
			monthOrder = append(monthOrder, key)
		}
		if tx.Type == models.Income {
			monthMap[key].income += tx.Amount.ToFloat64()
		} else {
			monthMap[key].expense += tx.Amount.ToFloat64()
		}
	}

	// Build 2D float slice
	months := make([]string, len(monthOrder))
	totals := make([][]float64, len(monthOrder))
	for i, m := range monthOrder {
		months[i] = m
		totals[i] = []float64{monthMap[m].income, monthMap[m].expense}
	}

	return months, totals
}

// ============================================================
// UNIT 3: Closures — Functions that capture outer variables
// ============================================================

// MakeBalanceTracker returns a closure that tracks running balance.
// Demonstrates: Closure (Unit 3) — returned function captures 'balance'
func MakeBalanceTracker(initialBalance models.Rupees) func(models.Transaction) models.Rupees {
	balance := initialBalance
	// This anonymous function closes over 'balance'
	return func(tx models.Transaction) models.Rupees {
		if tx.Type == models.Income {
			balance += tx.Amount
		} else {
			balance -= tx.Amount
		}
		return balance
	}
}

// ============================================================
// UNIT 3: Callback — Passing functions as arguments
// ============================================================

// TransactionCallback is a function type that operates on a transaction.
// Demonstrates: Function type definition
type TransactionCallback func(models.Transaction) bool

// ForEachTransaction applies a callback to each transaction.
// Demonstrates: Callback pattern (Unit 3)
func ForEachTransaction(transactions []models.Transaction, cb TransactionCallback) []models.Transaction {
	result := make([]models.Transaction, 0)
	for _, tx := range transactions {
		if cb(tx) {
			result = append(result, tx)
		}
	}
	return result
}

// ============================================================
// UNIT 3: Function Expression & Anonymous Function
// ============================================================

// GetFilterFunc returns a filter function based on criteria.
// Demonstrates: Function expression, anonymous function, returning a function (Unit 3)
func GetFilterFunc(filterType string) TransactionCallback {
	// Demonstrates: Switch with function expressions (Unit 1 + Unit 3)
	switch filterType {
	case "expenses_only":
		// Anonymous function assigned to return
		return func(tx models.Transaction) bool {
			return tx.Type == models.Expense
		}
	case "income_only":
		return func(tx models.Transaction) bool {
			return tx.Type == models.Income
		}
	case "large":
		return func(tx models.Transaction) bool {
			return tx.Amount > 5000
		}
	case "today":
		today := time.Now().Format("2006-01-02")
		// Demonstrates: Closure capturing 'today' variable
		return func(tx models.Transaction) bool {
			return tx.Date == today
		}
	default:
		return func(tx models.Transaction) bool {
			return true // pass-through
		}
	}
}

// ============================================================
// UNIT 3: Recursion
// ============================================================

// RecursiveCategorySum recursively sums amounts for a category tree.
// Demonstrates: Recursion (Unit 3)
// categoryTree maps parent categories to subcategories
func RecursiveCategorySum(
	transactions []models.Transaction,
	categoryTree map[models.Category][]models.Category,
	rootCategory models.Category,
) models.Rupees {
	var total models.Rupees

	// Sum direct transactions for this category
	for _, tx := range transactions {
		if tx.Category == rootCategory {
			total += tx.Amount
		}
	}

	// Recursively sum subcategories
	if subCats, ok := categoryTree[rootCategory]; ok {
		for _, sub := range subCats {
			total += RecursiveCategorySum(transactions, categoryTree, sub)
		}
	}

	return total
}

// ============================================================
// Sorting helpers
// ============================================================

// SortByDate sorts transactions by date (newest first).
func SortByDate(transactions []models.Transaction) {
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Date > transactions[j].Date
	})
}

// SortByAmount sorts transactions by amount (highest first).
func SortByAmount(transactions []models.Transaction) {
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Amount > transactions[j].Amount
	})
}
