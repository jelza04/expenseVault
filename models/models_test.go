package models

import (
	"testing"
	"time"
)

// ============================================================
// Unit Tests for Models Package
// ============================================================

func TestRupeesString(t *testing.T) {
	tests := []struct {
		amount   Rupees
		expected string
	}{
		{0, "₹0.00"},
		{100.50, "₹100.50"},
		{1500, "₹1500.00"},
		{-250.75, "₹-250.75"},
	}

	for _, tt := range tests {
		result := tt.amount.String()
		if result != tt.expected {
			t.Errorf("Rupees(%f).String() = %q, want %q", float64(tt.amount), result, tt.expected)
		}
	}
}

func TestRupeesOperations(t *testing.T) {
	a := Rupees(1000)
	b := Rupees(300)

	if got := a.Add(b); got != 1300 {
		t.Errorf("Add: got %v, want 1300", got)
	}
	if got := a.Subtract(b); got != 700 {
		t.Errorf("Subtract: got %v, want 700", got)
	}
	if got := a.ToFloat64(); got != 1000.0 {
		t.Errorf("ToFloat64: got %v, want 1000.0", got)
	}
	if a.IsNegative() {
		t.Error("1000 should not be negative")
	}
	if !Rupees(-1).IsNegative() {
		t.Error("-1 should be negative")
	}
}

func TestCategoriesForType(t *testing.T) {
	expCats := CategoriesForType(Expense)
	if len(expCats) != 8 {
		t.Errorf("Expected 8 expense categories, got %d", len(expCats))
	}

	incCats := CategoriesForType(Income)
	if len(incCats) != 4 {
		t.Errorf("Expected 4 income categories, got %d", len(incCats))
	}
}

func TestTransactionMethods(t *testing.T) {
	tx := Transaction{
		ID: 1, Type: Expense, Amount: 500,
		Category: CategoryFood, Description: "Test",
		Date:     "2026-02-19",
		DateInfo: DateInfo{CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	if !tx.IsExpense() {
		t.Error("Expected IsExpense() to be true")
	}
	if tx.IsIncome() {
		t.Error("Expected IsIncome() to be false")
	}

	parsed, err := tx.ParsedDate()
	if err != nil {
		t.Errorf("ParsedDate error: %v", err)
	}
	if parsed.Year() != 2026 {
		t.Errorf("Expected year 2026, got %d", parsed.Year())
	}

	my := tx.MonthYear()
	if my != "February 2026" {
		t.Errorf("MonthYear = %q, want 'February 2026'", my)
	}
}

func TestQuickSummary(t *testing.T) {
	txs := []Transaction{
		{Type: Income, Amount: 50000},
		{Type: Expense, Amount: 1000},
		{Type: Expense, Amount: 500},
	}

	summary := QuickSummary(txs)
	if summary.TotalIncome != 50000 {
		t.Errorf("TotalIncome = %v, want 50000", summary.TotalIncome)
	}
	if summary.TotalExpense != 1500 {
		t.Errorf("TotalExpense = %v, want 1500", summary.TotalExpense)
	}
	if summary.Balance != 48500 {
		t.Errorf("Balance = %v, want 48500", summary.Balance)
	}
	if summary.Count != 3 {
		t.Errorf("Count = %d, want 3", summary.Count)
	}
}

func TestValidateTransaction(t *testing.T) {
	// Valid transaction
	valid := Transaction{
		Type: Expense, Amount: 100, Category: CategoryFood,
		Description: "Test", Date: "2026-01-15",
	}
	if err := ValidateTransaction(valid); err != nil {
		t.Errorf("Valid transaction should not error: %v", err)
	}

	// Invalid amount
	invalid := Transaction{Type: Expense, Amount: -100, Description: "Test", Date: "2026-01-15"}
	if err := ValidateTransaction(invalid); err == nil {
		t.Error("Negative amount should fail validation")
	}

	// Empty description
	noDesc := Transaction{Type: Expense, Amount: 100, Description: "", Date: "2026-01-15"}
	if err := ValidateTransaction(noDesc); err == nil {
		t.Error("Empty description should fail validation")
	}

	// Invalid type
	badType := Transaction{Type: "invalid", Amount: 100, Description: "Test", Date: "2026-01-15"}
	if err := ValidateTransaction(badType); err == nil {
		t.Error("Invalid type should fail validation")
	}

	// Invalid date
	badDate := Transaction{Type: Expense, Amount: 100, Description: "Test", Date: "not-a-date", Category: CategoryFood}
	if err := ValidateTransaction(badDate); err == nil {
		t.Error("Invalid date should fail validation")
	}
}

func TestCustomErrors(t *testing.T) {
	// ValidationError
	ve := &ValidationError{Field: "amount", Message: "must be positive"}
	if ve.Error() != "validation error [amount]: must be positive" {
		t.Errorf("Unexpected ValidationError: %s", ve.Error())
	}

	// NotFoundError
	nf := &NotFoundError{Resource: "Transaction", ID: 42}
	if nf.Error() != "Transaction with ID 42 not found" {
		t.Errorf("Unexpected NotFoundError: %s", nf.Error())
	}

	// ParseError
	pe := &ParseError{Input: "abc", Expected: "number", Line: 5}
	expected := `parse error at line 5: expected number, got "abc"`
	if pe.Error() != expected {
		t.Errorf("Unexpected ParseError: %s", pe.Error())
	}
}

func TestSampleTransactions(t *testing.T) {
	samples := SampleTransactions()
	if len(samples) != 5 {
		t.Errorf("Expected 5 sample transactions, got %d", len(samples))
	}
}
