package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"expenseVault/models"
)

// ============================================================
// UNIT 3: Interface implementation, Defer, Error Handling
// CSV Handling using encoding/csv (standard library)
// ============================================================

// CSVHandler implements both Exporter and Importer for CSV data.
// Demonstrates: Struct implementing interfaces (Unit 3: polymorphism)
type CSVHandler struct{}

func (c *CSVHandler) Format() string {
	return "CSV"
}

// Export writes transactions to a CSV file.
// Demonstrates: defer for file cleanup (Unit 3), error handling
func (c *CSVHandler) Export(transactions []models.Transaction, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	// Demonstrates: defer for resource cleanup (Unit 3)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Type", "Amount", "Category", "Description", "Date", "Notes"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Demonstrates: for range over slice (Unit 2)
	for _, tx := range transactions {
		record := []string{
			strconv.Itoa(tx.ID),
			string(tx.Type),
			fmt.Sprintf("%.2f", tx.Amount.ToFloat64()),
			string(tx.Category),
			tx.Description,
			tx.Date,
			tx.Notes,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

// Import reads transactions from a CSV file.
// Demonstrates: Error handling with line info (Unit 3), slice operations (Unit 2)
func (c *CSVHandler) Import(filePath string) ([]models.Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, &models.ParseError{
			Input:    filePath,
			Expected: "CSV with header and at least one data row",
		}
	}

	// Demonstrates: make to pre-allocate slice (Unit 2)
	transactions := make([]models.Transaction, 0, len(records)-1)

	// Skip header row (index 0), process data rows
	// Demonstrates: slicing a slice (Unit 2) — records[1:]
	for lineNum, record := range records[1:] {
		if len(record) < 6 {
			return nil, &models.ParseError{
				Input:    fmt.Sprintf("line %d", lineNum+2),
				Expected: "at least 6 fields",
				Line:     lineNum + 2,
			}
		}

		amount, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, &models.ParseError{
				Input:    record[2],
				Expected: "valid number for amount",
				Line:     lineNum + 2,
			}
		}

		notes := ""
		if len(record) > 6 {
			notes = record[6]
		}

		now := time.Now()
		tx := models.Transaction{
			Type:        models.TransactionType(record[1]),
			Amount:      models.Rupees(amount),
			Category:    models.Category(record[3]),
			Description: record[4],
			Date:        record[5],
			Notes:       notes,
			DateInfo:    models.DateInfo{CreatedAt: now, UpdatedAt: now},
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}
