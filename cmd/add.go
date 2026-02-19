package cmd

import (
	"fmt"
	"time"

	"expenseVault/models"
	"expenseVault/services"

	"github.com/spf13/cobra"
)

// ============================================================
// UNIT 1: Variables (var, :=), fmt package, control flow (if-else)
// UNIT 3: Error handling — checking errors, printing errors
// ============================================================

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new transaction (income or expense)",
	Long: `Add a new transaction to your ExpenseVault.

Examples:
  expensevault add --type expense --amount 250 --category Food --desc "Lunch" --date 2026-02-19
  expensevault add -t income -a 50000 -c Salary -d "Monthly salary"
  expensevault add --type expense --amount 800 --desc "Uber ride" --auto-cat`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Demonstrates: Short declaration operator := (Unit 1)
		txType, _ := cmd.Flags().GetString("type")
		amount, _ := cmd.Flags().GetFloat64("amount")
		category, _ := cmd.Flags().GetString("category")
		desc, _ := cmd.Flags().GetString("desc")
		date, _ := cmd.Flags().GetString("date")
		notes, _ := cmd.Flags().GetString("notes")
		autoCat, _ := cmd.Flags().GetBool("auto-cat")

		// Demonstrates: Control flow — if-else validation (Unit 1)
		if txType != "income" && txType != "expense" {
			return fmt.Errorf("type must be 'income' or 'expense', got '%s'", txType)
		}

		if amount <= 0 {
			return fmt.Errorf("amount must be positive, got %.2f", amount)
		}

		// Demonstrates: var keyword with zero value (Unit 1)
		var cat models.Category

		// Demonstrates: if-else control flow (Unit 1)
		if autoCat {
			// Use auto-categorization from the categorizer service
			categorizer := services.NewCategorizer()
			cat = categorizer.AutoCategorize(desc)
			fmt.Printf("🏷️  Auto-categorized as: %s\n", cat)
		} else if category != "" {
			cat = models.Category(category)
		} else {
			cat = models.CategoryOther
		}

		// Default date to today if not provided
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}

		// Demonstrates: Composite literal (Unit 2)
		tx := models.Transaction{
			Type:        models.TransactionType(txType),
			Amount:      models.Rupees(amount),
			Category:    cat,
			Description: desc,
			Date:        date,
			Notes:       notes,
		}

		// Demonstrates: Error checking (Unit 3)
		id, err := store.AddTransaction(tx)
		if err != nil {
			return fmt.Errorf("failed to add transaction: %w", err)
		}

		// Demonstrates: fmt package for output (Unit 1)
		symbol := "💰"
		if txType == "expense" {
			symbol = "💸"
		}
		fmt.Printf("\n%s Transaction added successfully!\n", symbol)
		fmt.Printf("   ID:       %d\n", id)
		fmt.Printf("   Type:     %s\n", txType)
		fmt.Printf("   Amount:   %s\n", models.Rupees(amount))
		fmt.Printf("   Category: %s\n", cat)
		fmt.Printf("   Date:     %s\n", date)
		fmt.Printf("   Desc:     %s\n", desc)
		if notes != "" {
			fmt.Printf("   Notes:    %s\n", notes)
		}
		fmt.Println()

		return nil
	},
}

func init() {
	addCmd.Flags().StringP("type", "t", "", "Transaction type: income or expense (required)")
	addCmd.Flags().Float64P("amount", "a", 0, "Amount in rupees (required)")
	addCmd.Flags().StringP("category", "c", "", "Category (Food, Travel, Shopping, etc.)")
	addCmd.Flags().StringP("desc", "d", "", "Description (required)")
	addCmd.Flags().String("date", "", "Date in YYYY-MM-DD format (default: today)")
	addCmd.Flags().StringP("notes", "n", "", "Additional notes")
	addCmd.Flags().Bool("auto-cat", false, "Auto-categorize based on description")
	addCmd.MarkFlagRequired("type")
	addCmd.MarkFlagRequired("amount")
	addCmd.MarkFlagRequired("desc")
}
