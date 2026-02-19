package services

import (
	"testing"

	"expenseVault/models"
)

// ============================================================
// Unit Tests for Services Package
// ============================================================

func TestDeleteFromSlice(t *testing.T) {
	txs := []models.Transaction{
		{ID: 1, Description: "A"},
		{ID: 2, Description: "B"},
		{ID: 3, Description: "C"},
	}

	// Delete middle element
	result := DeleteFromSlice(txs, 1)
	if len(result) != 2 {
		t.Errorf("Expected 2 items after delete, got %d", len(result))
	}
	if result[0].ID != 1 || result[1].ID != 3 {
		t.Error("Wrong items remaining after delete")
	}

	// Delete out of bounds
	result2 := DeleteFromSlice(txs, 10)
	if len(result2) != 3 {
		t.Error("Out of bounds delete should return original slice")
	}
}

func TestGroupByMonth(t *testing.T) {
	txs := []models.Transaction{
		{ID: 1, Date: "2026-01-15"},
		{ID: 2, Date: "2026-01-20"},
		{ID: 3, Date: "2026-02-10"},
	}

	groups := GroupByMonth(txs)
	if len(groups) != 2 {
		t.Errorf("Expected 2 month groups, got %d", len(groups))
	}
}

func TestMakeBalanceTracker(t *testing.T) {
	tracker := MakeBalanceTracker(10000)

	bal := tracker(models.Transaction{Type: models.Expense, Amount: 500})
	if bal != 9500 {
		t.Errorf("After expense 500: expected 9500, got %v", bal)
	}

	bal = tracker(models.Transaction{Type: models.Income, Amount: 2000})
	if bal != 11500 {
		t.Errorf("After income 2000: expected 11500, got %v", bal)
	}
}

func TestForEachTransaction(t *testing.T) {
	txs := []models.Transaction{
		{Type: models.Income, Amount: 1000},
		{Type: models.Expense, Amount: 200},
		{Type: models.Expense, Amount: 300},
		{Type: models.Income, Amount: 500},
	}

	expenses := ForEachTransaction(txs, func(tx models.Transaction) bool {
		return tx.Type == models.Expense
	})

	if len(expenses) != 2 {
		t.Errorf("Expected 2 expenses, got %d", len(expenses))
	}
}

func TestGetFilterFunc(t *testing.T) {
	txs := []models.Transaction{
		{Type: models.Income, Amount: 10000},
		{Type: models.Expense, Amount: 200},
		{Type: models.Expense, Amount: 6000},
	}

	// expenses_only filter
	expFilter := GetFilterFunc("expenses_only")
	result := ForEachTransaction(txs, expFilter)
	if len(result) != 2 {
		t.Errorf("expenses_only: expected 2, got %d", len(result))
	}

	// large filter
	largeFilter := GetFilterFunc("large")
	result = ForEachTransaction(txs, largeFilter)
	if len(result) != 2 {
		t.Errorf("large: expected 2, got %d", len(result))
	}

	// default filter
	allFilter := GetFilterFunc("unknown")
	result = ForEachTransaction(txs, allFilter)
	if len(result) != 3 {
		t.Errorf("default: expected 3, got %d", len(result))
	}
}

func TestRecursiveCategorySum(t *testing.T) {
	txs := []models.Transaction{
		{Category: models.CategoryFood, Amount: 500},
		{Category: models.CategoryTravel, Amount: 300},
		{Category: models.CategoryBills, Amount: 1000},
		{Category: models.CategoryHealth, Amount: 200},
	}

	tree := map[models.Category][]models.Category{
		"All":    {models.CategoryFood, "Living"},
		"Living": {models.CategoryBills, models.CategoryHealth},
	}

	total := RecursiveCategorySum(txs, tree, "All")
	// Food(500) + Bills(1000) + Health(200) = 1700
	if total != 1700 {
		t.Errorf("RecursiveCategorySum: expected 1700, got %v", total)
	}
}

func TestAutoCategorize(t *testing.T) {
	cat := NewCategorizer()

	tests := []struct {
		desc     string
		expected models.Category
	}{
		{"Lunch at canteen", models.CategoryFood},
		{"Uber to office", models.CategoryTravel},
		{"Amazon order", models.CategoryShopping},
		{"Electricity bill payment", models.CategoryBills},
		{"Something random here", models.CategoryOther},
	}

	for _, tt := range tests {
		result := cat.AutoCategorize(tt.desc)
		if result != tt.expected {
			t.Errorf("AutoCategorize(%q) = %q, want %q", tt.desc, result, tt.expected)
		}
	}
}

func TestMakeKeywordMatcher(t *testing.T) {
	matcher := MakeKeywordMatcher("food", "lunch", "dinner")

	if !matcher("Had lunch today") {
		t.Error("Should match 'lunch'")
	}
	if !matcher("FOOD delivery") {
		t.Error("Should match 'FOOD' (case insensitive)")
	}
	if matcher("Uber ride") {
		t.Error("Should not match 'Uber ride'")
	}
}

func TestReporters(t *testing.T) {
	txs := models.SampleTransactions()

	reporters := []Reporter{
		&MonthlyReporter{},
		&CategoryReporter{},
		&YearlyReporter{},
	}

	for _, r := range reporters {
		name := r.Name()
		if name == "" {
			t.Error("Reporter name should not be empty")
		}

		output := r.Generate(txs)
		if output == "" {
			t.Errorf("%s generated empty output", name)
		}
	}
}

func TestMonthlyTotals(t *testing.T) {
	txs := []models.Transaction{
		{Type: models.Income, Amount: 50000, Date: "2026-01-15"},
		{Type: models.Expense, Amount: 1000, Date: "2026-01-20"},
	}

	months, totals := MonthlyTotals(txs)
	if len(months) != 1 {
		t.Errorf("Expected 1 month, got %d", len(months))
	}
	if len(totals) != 1 || len(totals[0]) != 2 {
		t.Error("Unexpected totals structure")
	}
	if totals[0][0] != 50000 {
		t.Errorf("Income total: expected 50000, got %v", totals[0][0])
	}
	if totals[0][1] != 1000 {
		t.Errorf("Expense total: expected 1000, got %v", totals[0][1])
	}
}
