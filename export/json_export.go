package export

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"expenseVault/models"
)

// ============================================================
// JSON Handling using encoding/json (standard library)
// Demonstrates: Data serialization with JSON, defer, error handling
// ============================================================

// JSONHandler implements both Exporter and Importer for JSON data.
// Demonstrates: Struct implementing interfaces (Unit 3: polymorphism)
type JSONHandler struct{}

func (j *JSONHandler) Format() string {
	return "JSON"
}

// Export writes transactions to a JSON file.
// Demonstrates: encoding/json, defer (Unit 3)
func (j *JSONHandler) Export(transactions []models.Transaction, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	// Wrap in a backup envelope with metadata
	// Demonstrates: Anonymous struct (Unit 2) for one-off serialization
	backup := struct {
		Version      string               `json:"version"`
		ExportedAt   time.Time            `json:"exported_at"`
		Count        int                  `json:"count"`
		Transactions []models.Transaction `json:"transactions"`
	}{
		Version:      "1.0",
		ExportedAt:   time.Now(),
		Count:        len(transactions),
		Transactions: transactions,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(backup); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// Import reads transactions from a JSON file.
// Demonstrates: JSON deserialization, error handling (Unit 3)
func (j *JSONHandler) Import(filePath string) ([]models.Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	var backup struct {
		Transactions []models.Transaction `json:"transactions"`
	}

	if err := json.NewDecoder(file).Decode(&backup); err != nil {
		return nil, &models.ParseError{
			Input:    filePath,
			Expected: "valid JSON with transactions array",
		}
	}

	return backup.Transactions, nil
}

// ExportToString serializes transactions to a JSON string (for API responses).
func ExportToString(transactions []models.Transaction) (string, error) {
	data, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
