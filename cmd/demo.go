package cmd

import (
	"fmt"
	"time"

	"expenseVault/models"
	"expenseVault/services"

	"github.com/spf13/cobra"
)

// ============================================================
// Demo command — showcases ALL syllabus concepts in action
// ============================================================

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run a demonstration of all Go concepts used in ExpenseVault",
	Long: `Runs a live demo showing: custom types, zero values, type conversion,
arrays, slices, maps, structs, functions, closures, callbacks,
interfaces, polymorphism, recursion, and error handling.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n🎓 ExpenseVault — Go Concepts Demonstration")
		fmt.Println("============================================\n")

		// ── UNIT 1: Zero Values ──
		fmt.Println("📘 UNIT 1: Zero Values Demo")
		models.ZeroValueDemo()

		// ── UNIT 1: Custom Types & Type Conversion ──
		fmt.Println("\n📘 UNIT 1: Custom Types & Type Conversion")
		var amount models.Rupees = 1500.50
		fmt.Printf("  Rupees value: %s\n", amount)
		fmt.Printf("  As float64 (type conversion): %f\n", amount.ToFloat64())
		fmt.Printf("  Addition: %s + ₹500 = %s\n", amount, amount.Add(500))
		fmt.Printf("  Is Negative? %v\n", amount.IsNegative())

		// ── UNIT 1: Short declaration ──
		fmt.Println("\n📘 UNIT 1: Short Declaration Operator (:=)")
		name := "ExpenseVault"
		version := 1.0
		active := true
		fmt.Printf("  name := %q, version := %.1f, active := %v\n", name, version, active)

		// ── UNIT 2: Arrays ──
		fmt.Println("\n📘 UNIT 2: Arrays (Fixed Size)")
		fmt.Printf("  Expense categories (array[8]): %v\n", models.DefaultExpenseCategories)
		fmt.Printf("  Income categories (array[4]):  %v\n", models.DefaultIncomeCategories)

		// ── UNIT 2: Slices ──
		fmt.Println("\n📘 UNIT 2: Slices")
		transactions := models.SampleTransactions()
		fmt.Printf("  Sample transactions (slice, len=%d): \n", len(transactions))
		for i, tx := range transactions {
			fmt.Printf("    [%d] %s\n", i, tx)
		}

		// Demonstrates: append
		newTx := models.Transaction{
			ID: 6, Type: models.Expense, Amount: 350,
			Category: models.CategoryFood, Description: "Coffee",
			Date:     time.Now().Format("2006-01-02"),
			DateInfo: models.DateInfo{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		transactions = append(transactions, newTx)
		fmt.Printf("\n  After append (len=%d): added '%s'\n", len(transactions), newTx.Description)

		// Demonstrates: delete from slice
		fmt.Printf("  Before delete: len=%d\n", len(transactions))
		transactions = services.DeleteFromSlice(transactions, 2)
		fmt.Printf("  After delete index 2: len=%d\n", len(transactions))

		// Demonstrates: slicing a slice
		recent := transactions[:3]
		fmt.Printf("  First 3 (slicing): %d items\n", len(recent))

		// Demonstrates: make
		filtered := make([]models.Transaction, 0, 10)
		fmt.Printf("  make([]Transaction, 0, 10): len=%d, cap=%d\n", len(filtered), cap(filtered))

		// ── UNIT 2: Multi-dimensional Slices ──
		fmt.Println("\n📘 UNIT 2: Multi-dimensional Slices")
		grouped := services.GroupByMonth(transactions)
		fmt.Printf("  Grouped by month: %d groups\n", len(grouped))
		for i, group := range grouped {
			fmt.Printf("    Group %d: %d transactions\n", i, len(group))
		}

		months, totals := services.MonthlyTotals(transactions)
		fmt.Println("  Monthly totals ([][]float64):")
		for i, m := range months {
			fmt.Printf("    %s: income=%.2f, expense=%.2f\n", m, totals[i][0], totals[i][1])
		}

		// ── UNIT 2: Maps ──
		fmt.Println("\n📘 UNIT 2: Maps")
		catTotals := make(map[models.Category]models.Rupees)
		for _, tx := range transactions {
			if tx.Type == models.Expense {
				catTotals[tx.Category] += tx.Amount
			}
		}
		fmt.Println("  Category totals (map[Category]Rupees):")
		for cat, total := range catTotals {
			fmt.Printf("    %-15s → %s\n", cat, total)
		}

		// Demonstrates: delete from map
		delete(catTotals, models.CategoryOther)
		fmt.Printf("  After delete(map, Other): %d entries\n", len(catTotals))

		// ── UNIT 2: Structs, Embedded Structs, Anonymous Structs ──
		fmt.Println("\n📘 UNIT 2: Structs & Embedded Structs")
		tx := transactions[0]
		fmt.Printf("  Transaction struct: %s\n", tx)
		fmt.Printf("  Embedded DateInfo.CreatedAt: %s\n", tx.FormattedCreatedAt())
		fmt.Printf("  Method: IsExpense() = %v\n", tx.IsExpense())

		fmt.Println("\n📘 UNIT 2: Anonymous Struct")
		summary := models.QuickSummary(transactions)
		fmt.Printf("  QuickSummary: Income=%s, Expense=%s, Balance=%s, Count=%d\n",
			summary.TotalIncome, summary.TotalExpense, summary.Balance, summary.Count)

		// ── UNIT 3: Closures ──
		fmt.Println("\n📘 UNIT 3: Closures")
		tracker := services.MakeBalanceTracker(50000)
		fmt.Println("  Balance tracker (closure captures 'balance' variable):")
		for _, tx := range transactions {
			bal := tracker(tx)
			fmt.Printf("    After %s %s: Balance = %s\n", tx.Type, tx.Amount, bal)
		}

		// ── UNIT 3: Callbacks ──
		fmt.Println("\n📘 UNIT 3: Callbacks")
		expensesOnly := services.ForEachTransaction(transactions, func(tx models.Transaction) bool {
			return tx.Type == models.Expense
		})
		fmt.Printf("  Filtered expenses via callback: %d items\n", len(expensesOnly))

		// ── UNIT 3: Function Expressions & Returning Functions ──
		fmt.Println("\n📘 UNIT 3: Function Expressions & Returning Functions")
		largeFilter := services.GetFilterFunc("large")
		largeItems := services.ForEachTransaction(transactions, largeFilter)
		fmt.Printf("  Large transactions (>5000): %d items\n", len(largeItems))

		// ── UNIT 3: Keyword matcher (returning a function) ──
		matcher := services.MakeKeywordMatcher("food", "lunch", "coffee")
		fmt.Printf("  Keyword matcher for 'Lunch at canteen': %v\n", matcher("Lunch at canteen"))
		fmt.Printf("  Keyword matcher for 'Uber ride': %v\n", matcher("Uber ride"))

		// ── UNIT 3: Auto-categorization ──
		fmt.Println("\n📘 UNIT 3: Auto-Categorization (Callbacks + Rules)")
		categorizer := services.NewCategorizer()
		testDescs := []string{"Lunch at canteen", "Uber to office", "Amazon order", "Electricity bill", "Movie tickets"}
		for _, d := range testDescs {
			fmt.Printf("  %-25s → %s\n", d, categorizer.AutoCategorize(d))
		}

		// ── UNIT 3: Interfaces & Polymorphism ──
		fmt.Println("\n📘 UNIT 3: Interfaces & Polymorphism")
		reporters := []services.Reporter{
			&services.MonthlyReporter{},
			&services.CategoryReporter{},
			&services.YearlyReporter{},
		}
		for _, r := range reporters {
			fmt.Printf("  Reporter: %-20s — implements Reporter interface\n", r.Name())
		}

		// ── UNIT 3: Recursion ──
		fmt.Println("\n📘 UNIT 3: Recursion")
		categoryTree := map[models.Category][]models.Category{
			"All Expenses": {models.CategoryFood, models.CategoryTravel, "Living"},
			"Living":       {models.CategoryBills, models.CategoryHealth},
		}
		total := services.RecursiveCategorySum(transactions, categoryTree, "All Expenses")
		fmt.Printf("  Recursive sum for 'All Expenses' tree: %s\n", total)

		// ── UNIT 3: Error Handling ──
		fmt.Println("\n📘 UNIT 3: Error Handling")
		invalidTx := models.Transaction{Amount: -100, Description: "", Type: "invalid"}
		if err := models.ValidateTransaction(invalidTx); err != nil {
			fmt.Printf("  Validation error: %v\n", err)
		}

		notFound := &models.NotFoundError{Resource: "Transaction", ID: 999}
		fmt.Printf("  NotFoundError: %v\n", notFound)

		dbErr := &models.DatabaseError{Operation: "insert", Err: fmt.Errorf("disk full")}
		fmt.Printf("  DatabaseError: %v\n", dbErr)
		fmt.Printf("  Unwrapped: %v\n", dbErr.Unwrap())

		parseErr := &models.ParseError{Input: "abc", Expected: "number", Line: 5}
		fmt.Printf("  ParseError: %v\n", parseErr)

		// ── UNIT 3: Defer, Panic, Recover ──
		fmt.Println("\n📘 UNIT 3: Defer, Panic, Recover")
		fmt.Println("  Running panic/recover demo...")
		deferPanicDemo()

		fmt.Println("\n✅ All Go concepts demonstrated successfully!")
		fmt.Println("============================================\n")
	},
}

// deferPanicDemo demonstrates defer, panic, and recover.
// Demonstrates: defer ordering, panic, recover (Unit 3)
func deferPanicDemo() {
	// Demonstrates: defer runs in LIFO order
	defer fmt.Println("  [defer 1] First deferred — runs last")
	defer fmt.Println("  [defer 2] Second deferred — runs second")
	defer fmt.Println("  [defer 3] Third deferred — runs first")

	// Demonstrates: recover from panic
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("  [recover] Recovered from panic: %v\n", r)
		}
	}()

	fmt.Println("  [normal] About to panic...")
	panic("simulated crash in ExpenseVault!")
}
