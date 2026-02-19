package models

import (
	"fmt"
	"time"
)

// ============================================================
// UNIT 2: Structs — Introduction, Embedded Structs, Anonymous Structs
// ============================================================

// DateInfo is embedded into Transaction.
// Demonstrates: Struct that will be embedded (embedded structs)
type DateInfo struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FormattedCreatedAt returns a human-readable created date.
func (d DateInfo) FormattedCreatedAt() string {
	return d.CreatedAt.Format("2006-01-02 15:04")
}

// FormattedUpdatedAt returns a human-readable updated date.
func (d DateInfo) FormattedUpdatedAt() string {
	return d.UpdatedAt.Format("2006-01-02 15:04")
}

// Transaction represents a single income/expense entry.
// Demonstrates: Struct with multiple field types, struct tags, embedded struct
type Transaction struct {
	ID          int             `json:"id"`
	Type        TransactionType `json:"type"`
	Amount      Rupees          `json:"amount"`
	Category    Category        `json:"category"`
	Description string          `json:"description"`
	Date        string          `json:"date"` // YYYY-MM-DD format
	Notes       string          `json:"notes"`
	DateInfo                    // Embedded struct (Unit 2: embedded structs)
}

// String provides a formatted representation of a Transaction.
// Demonstrates: fmt.Sprintf, method on struct
func (t Transaction) String() string {
	symbol := "💰"
	if t.Type == Expense {
		symbol = "💸"
	}
	return fmt.Sprintf("%s #%d | %s | %-15s | %-12s | %s | %s",
		symbol, t.ID, t.Date, t.Category, t.Amount, t.Description, t.Notes)
}

// IsExpense returns true if this is an expense transaction.
func (t Transaction) IsExpense() bool {
	return t.Type == Expense
}

// IsIncome returns true if this is an income transaction.
func (t Transaction) IsIncome() bool {
	return t.Type == Income
}

// ParsedDate converts the Date string to time.Time.
func (t Transaction) ParsedDate() (time.Time, error) {
	return time.Parse("2006-01-02", t.Date)
}

// MonthYear returns the month and year of this transaction.
func (t Transaction) MonthYear() string {
	parsed, err := t.ParsedDate()
	if err != nil {
		return "Unknown"
	}
	return parsed.Format("January 2006")
}

// ============================================================
// UNIT 2: Anonymous Structs — Used for quick, one-off groupings
// ============================================================

// QuickSummary returns an anonymous struct with transaction summary.
// Demonstrates: Anonymous structs for one-off data grouping
func QuickSummary(transactions []Transaction) struct {
	TotalIncome  Rupees
	TotalExpense Rupees
	Balance      Rupees
	Count        int
} {
	// Demonstrates: Anonymous struct (Unit 2)
	result := struct {
		TotalIncome  Rupees
		TotalExpense Rupees
		Balance      Rupees
		Count        int
	}{}

	// Demonstrates: for range loop (Unit 1)
	for _, tx := range transactions {
		if tx.Type == Income {
			result.TotalIncome += tx.Amount
		} else {
			result.TotalExpense += tx.Amount
		}
		result.Count++
	}

	result.Balance = result.TotalIncome - result.TotalExpense
	return result
}

// ============================================================
// UNIT 2: Composite Literals for creating test/sample data
// ============================================================

// SampleTransactions returns test data using composite literals.
// Demonstrates: Composite literals, slice literal
func SampleTransactions() []Transaction {
	now := time.Now()
	return []Transaction{
		{
			ID: 1, Type: Expense, Amount: 250.00,
			Category: CategoryFood, Description: "Lunch at canteen",
			Date: now.Format("2006-01-02"), Notes: "With friends",
			DateInfo: DateInfo{CreatedAt: now, UpdatedAt: now},
		},
		{
			ID: 2, Type: Expense, Amount: 1500.00,
			Category: CategoryShopping, Description: "Amazon order",
			Date: now.AddDate(0, 0, -1).Format("2006-01-02"), Notes: "Books",
			DateInfo: DateInfo{CreatedAt: now, UpdatedAt: now},
		},
		{
			ID: 3, Type: Income, Amount: 50000.00,
			Category: CategorySalary, Description: "Monthly salary",
			Date: now.AddDate(0, 0, -5).Format("2006-01-02"), Notes: "",
			DateInfo: DateInfo{CreatedAt: now, UpdatedAt: now},
		},
		{
			ID: 4, Type: Expense, Amount: 800.00,
			Category: CategoryTravel, Description: "Uber ride",
			Date: now.AddDate(0, 0, -2).Format("2006-01-02"), Notes: "Office commute",
			DateInfo: DateInfo{CreatedAt: now, UpdatedAt: now},
		},
		{
			ID: 5, Type: Expense, Amount: 2000.00,
			Category: CategoryBills, Description: "Electricity bill",
			Date: now.AddDate(0, 0, -3).Format("2006-01-02"), Notes: "February",
			DateInfo: DateInfo{CreatedAt: now, UpdatedAt: now},
		},
	}
}

// ============================================================
// Report types for reporting
// ============================================================

// ReportEntry represents a single line in a report.
type ReportEntry struct {
	Label  string
	Amount Rupees
	Count  int
}

// MonthlySummary holds data for a monthly report.
type MonthlySummary struct {
	Month        string
	Income       Rupees
	Expenses     Rupees
	Balance      Rupees
	ByCategory   map[Category]Rupees
	Transactions []Transaction
}

// User represents a user for backend sync.
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never serialize password hash
	CreatedAt    time.Time `json:"created_at"`
}

// SyncPayload is used for syncing transactions with the backend.
type SyncPayload struct {
	Transactions []Transaction `json:"transactions"`
	LastSyncAt   time.Time     `json:"last_sync_at"`
}
