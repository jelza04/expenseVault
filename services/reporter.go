package services

import (
	"fmt"
	"strings"

	"expenseVault/models"
)

// ============================================================
// UNIT 3: Interfaces & Polymorphism
// Demonstrates: Interface definition, multiple implementations,
//               polymorphic behavior
// UNIT 2: Maps for aggregation
// ============================================================

// Reporter is an interface for generating reports.
// Demonstrates: Interface definition (Unit 3)
type Reporter interface {
	Generate(transactions []models.Transaction) string
	Name() string
}

// ============================================================
// MonthlyReporter — Implementation 1
// ============================================================

// MonthlyReporter generates monthly income vs. expense reports.
// Demonstrates: Struct implementing an interface (Unit 3: polymorphism)
type MonthlyReporter struct{}

func (r *MonthlyReporter) Name() string {
	return "Monthly Report"
}

// Generate creates a monthly summary report.
// Demonstrates: Maps for aggregation (Unit 2), for range (Unit 2)
func (r *MonthlyReporter) Generate(transactions []models.Transaction) string {
	var sb strings.Builder

	sb.WriteString("\n╔══════════════════════════════════════════════════════════╗\n")
	sb.WriteString("║              📊 MONTHLY REPORT                          ║\n")
	sb.WriteString("╠══════════════════════════════════════════════════════════╣\n")

	// Demonstrates: Map for aggregation — map[string]*models.MonthlySummary (Unit 2)
	monthlyData := make(map[string]*models.MonthlySummary)

	// Demonstrates: for range over slice (Unit 2)
	for _, tx := range transactions {
		month := tx.MonthYear()
		if _, exists := monthlyData[month]; !exists {
			monthlyData[month] = &models.MonthlySummary{
				Month:      month,
				ByCategory: make(map[models.Category]models.Rupees),
			}
		}
		ms := monthlyData[month]
		ms.Transactions = append(ms.Transactions, tx)

		if tx.Type == models.Income {
			ms.Income += tx.Amount
		} else {
			ms.Expenses += tx.Amount
			// Demonstrates: Map — add/update element (Unit 2)
			ms.ByCategory[tx.Category] += tx.Amount
		}
		ms.Balance = ms.Income - ms.Expenses
	}

	// Demonstrates: for range over map (Unit 2)
	for month, data := range monthlyData {
		sb.WriteString(fmt.Sprintf("║ 📅 %s\n", month))
		sb.WriteString(fmt.Sprintf("║   Income:   %s\n", data.Income))
		sb.WriteString(fmt.Sprintf("║   Expenses: %s\n", data.Expenses))
		sb.WriteString(fmt.Sprintf("║   Balance:  %s\n", data.Balance))

		if len(data.ByCategory) > 0 {
			sb.WriteString("║   By Category:\n")
			for cat, amt := range data.ByCategory {
				sb.WriteString(fmt.Sprintf("║     %-15s %s\n", cat, amt))
			}
		}
		sb.WriteString("║\n")
		_ = month // suppress unused warning
	}

	sb.WriteString("╚══════════════════════════════════════════════════════════╝\n")
	return sb.String()
}

// ============================================================
// CategoryReporter — Implementation 2
// ============================================================

// CategoryReporter generates category-wise spending breakdown.
// Demonstrates: Another struct implementing the same Reporter interface (polymorphism)
type CategoryReporter struct{}

func (r *CategoryReporter) Name() string {
	return "Category Report"
}

// Generate creates a category-wise breakdown report.
// Demonstrates: map[Category]Rupees for category totals (Unit 2)
func (r *CategoryReporter) Generate(transactions []models.Transaction) string {
	var sb strings.Builder

	sb.WriteString("\n╔══════════════════════════════════════════════════════════╗\n")
	sb.WriteString("║            📂 CATEGORY-WISE REPORT                      ║\n")
	sb.WriteString("╠══════════════════════════════════════════════════════════╣\n")

	// Demonstrates: map for fast lookups and aggregation (Unit 2)
	categoryTotals := make(map[models.Category]models.Rupees)
	categoryCounts := make(map[models.Category]int)
	var totalExpense models.Rupees

	for _, tx := range transactions {
		if tx.Type == models.Expense {
			categoryTotals[tx.Category] += tx.Amount
			categoryCounts[tx.Category]++
			totalExpense += tx.Amount
		}
	}

	// Demonstrates: for range over map (Unit 2)
	for cat, total := range categoryTotals {
		count := categoryCounts[cat]
		percentage := 0.0
		if totalExpense > 0 {
			percentage = (total.ToFloat64() / totalExpense.ToFloat64()) * 100
		}
		bar := strings.Repeat("█", int(percentage/5))
		sb.WriteString(fmt.Sprintf("║ %-15s %s (%d txns, %.1f%%)\n", cat, total, count, percentage))
		sb.WriteString(fmt.Sprintf("║   %s\n", bar))
	}

	sb.WriteString(fmt.Sprintf("║\n║ Total Expenses: %s\n", totalExpense))
	sb.WriteString("╚══════════════════════════════════════════════════════════╝\n")
	return sb.String()
}

// ============================================================
// YearlyReporter — Implementation 3
// ============================================================

// YearlyReporter generates yearly summary.
// Demonstrates: Third implementation of the same interface (polymorphism)
type YearlyReporter struct{}

func (r *YearlyReporter) Name() string {
	return "Yearly Report"
}

func (r *YearlyReporter) Generate(transactions []models.Transaction) string {
	var sb strings.Builder

	sb.WriteString("\n╔══════════════════════════════════════════════════════════╗\n")
	sb.WriteString("║              📅 YEARLY REPORT                           ║\n")
	sb.WriteString("╠══════════════════════════════════════════════════════════╣\n")

	yearData := make(map[string]struct {
		income  models.Rupees
		expense models.Rupees
	})

	for _, tx := range transactions {
		parsed, err := tx.ParsedDate()
		if err != nil {
			continue
		}
		year := fmt.Sprintf("%d", parsed.Year())
		data := yearData[year]
		if tx.Type == models.Income {
			data.income += tx.Amount
		} else {
			data.expense += tx.Amount
		}
		yearData[year] = data
	}

	for year, data := range yearData {
		balance := data.income - data.expense
		sb.WriteString(fmt.Sprintf("║ 📅 Year %s\n", year))
		sb.WriteString(fmt.Sprintf("║   Total Income:   %s\n", data.income))
		sb.WriteString(fmt.Sprintf("║   Total Expenses: %s\n", data.expense))
		sb.WriteString(fmt.Sprintf("║   Net Balance:    %s\n", balance))
		sb.WriteString("║\n")
	}

	sb.WriteString("╚══════════════════════════════════════════════════════════╝\n")
	return sb.String()
}

// ============================================================
// UNIT 3: Polymorphism in action — using the Reporter interface
// ============================================================

// GenerateReport demonstrates polymorphism: a single function
// that works with any Reporter implementation.
// Demonstrates: Interface polymorphism (Unit 3)
func GenerateReport(reporter Reporter, transactions []models.Transaction) string {
	header := fmt.Sprintf("\n📋 Generating: %s\n", reporter.Name())
	return header + reporter.Generate(transactions)
}

// GetReporter returns the appropriate reporter based on report type.
// Also demonstrates: switch control flow (Unit 1)
func GetReporter(reportType string) Reporter {
	switch strings.ToLower(reportType) {
	case "monthly":
		return &MonthlyReporter{}
	case "category":
		return &CategoryReporter{}
	case "yearly":
		return &YearlyReporter{}
	default:
		return &MonthlyReporter{}
	}
}

// GenerateAllReports runs all reporter types.
// Demonstrates: Slice of interfaces (Unit 3: polymorphism)
func GenerateAllReports(transactions []models.Transaction) string {
	reporters := []Reporter{
		&MonthlyReporter{},
		&CategoryReporter{},
		&YearlyReporter{},
	}

	var sb strings.Builder
	for _, r := range reporters {
		sb.WriteString(GenerateReport(r, transactions))
		sb.WriteString("\n")
	}
	return sb.String()
}
