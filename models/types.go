// Package models defines all core types, structs, and error types for ExpenseVault.
//
// SYLLABUS CONCEPTS DEMONSTRATED:
// Unit 1: Custom types, type conversion (not casting), zero values, fmt package,
//
//	creating your own type, variables, values & types
//
// Unit 2: Arrays (fixed-size), structs
package models

import (
	"fmt"
	"time"
)

// ============================================================
// UNIT 1: Creating Your Own Type & Type Conversion
// ============================================================

// Rupees is a custom type wrapping float64 for currency amounts.
// Demonstrates: Creating your own type based on existing type
type Rupees float64

// String implements fmt.Stringer for pretty-printing currency.
// Demonstrates: fmt package, methods on custom types
func (r Rupees) String() string {
	return fmt.Sprintf("₹%.2f", float64(r))
}

// ToFloat64 converts Rupees to float64.
// Demonstrates: Type conversion (not casting — Go has no casting)
func (r Rupees) ToFloat64() float64 {
	return float64(r)
}

// Add returns the sum of two Rupees values.
func (r Rupees) Add(other Rupees) Rupees {
	return r + other
}

// Subtract returns the difference.
func (r Rupees) Subtract(other Rupees) Rupees {
	return r - other
}

// IsNegative checks if the amount is negative.
func (r Rupees) IsNegative() bool {
	return r < 0
}

// ============================================================
// UNIT 1: Category Type with Constants
// ============================================================

// Category represents a transaction category.
// Demonstrates: Custom string type
type Category string

// Predefined categories using typed constants.
// Demonstrates: const block, iota alternative with string constants
const (
	CategoryFood          Category = "Food"
	CategoryTravel        Category = "Travel"
	CategoryShopping      Category = "Shopping"
	CategoryBills         Category = "Bills"
	CategoryHealth        Category = "Health"
	CategoryEducation     Category = "Education"
	CategoryEntertainment Category = "Entertainment"
	CategorySalary        Category = "Salary"
	CategoryFreelance     Category = "Freelance"
	CategoryInvestment    Category = "Investment"
	CategoryOther         Category = "Other"
)

// TransactionType indicates income or expense.
// Demonstrates: Custom type for type safety
type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

// ============================================================
// UNIT 2: Arrays — Fixed-size category lists
// Demonstrates: Arrays (contrast with slices), composite literals
// ============================================================

// DefaultExpenseCategories is a fixed-size array of expense categories.
var DefaultExpenseCategories = [8]Category{
	CategoryFood, CategoryTravel, CategoryShopping, CategoryBills,
	CategoryHealth, CategoryEducation, CategoryEntertainment, CategoryOther,
}

// DefaultIncomeCategories is a fixed-size array of income categories.
var DefaultIncomeCategories = [4]Category{
	CategorySalary, CategoryFreelance, CategoryInvestment, CategoryOther,
}

// AllExpenseCategories returns a slice from the array.
// Demonstrates: Slicing an array to get a slice
func AllExpenseCategories() []Category {
	return DefaultExpenseCategories[:]
}

// AllIncomeCategories returns a slice from the array.
func AllIncomeCategories() []Category {
	return DefaultIncomeCategories[:]
}

// CategoriesForType returns categories based on transaction type.
// Demonstrates: Control flow with switch
func CategoriesForType(t TransactionType) []Category {
	switch t {
	case Income:
		return AllIncomeCategories()
	case Expense:
		return AllExpenseCategories()
	default:
		return AllExpenseCategories()
	}
}

// ============================================================
// UNIT 1: Zero Values Demonstration
// ============================================================

// ZeroValueDemo prints zero values of all types used in ExpenseVault.
// Demonstrates: Zero values concept — every Go type has a default zero value
func ZeroValueDemo() {
	var amount Rupees          // zero: ₹0.00
	var category Category      // zero: ""
	var txType TransactionType // zero: ""
	var count int              // zero: 0
	var rate float64           // zero: 0.0
	var name string            // zero: ""
	var active bool            // zero: false
	var timestamp time.Time    // zero: 0001-01-01 00:00:00

	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║       ZERO VALUES DEMONSTRATION          ║")
	fmt.Println("╠══════════════════════════════════════════╣")
	fmt.Printf("║ Rupees:          %s\n", amount)
	fmt.Printf("║ Category:        %q\n", category)
	fmt.Printf("║ TransactionType: %q\n", txType)
	fmt.Printf("║ int:             %d\n", count)
	fmt.Printf("║ float64:         %f\n", rate)
	fmt.Printf("║ string:          %q\n", name)
	fmt.Printf("║ bool:            %v\n", active)
	fmt.Printf("║ time.Time:       %v\n", timestamp)
	fmt.Println("╚══════════════════════════════════════════╝")
}
