package export

import (
	"expenseVault/models"
)

// ============================================================
// UNIT 3: Interfaces & Polymorphism
// Demonstrates: Interface for export functionality with multiple
//               implementations (CSV, JSON)
// ============================================================

// Exporter defines the interface for exporting transaction data.
// Demonstrates: Interface definition (Unit 3)
type Exporter interface {
	Export(transactions []models.Transaction, filePath string) error
	Format() string
}

// Importer defines the interface for importing transaction data.
// Demonstrates: Another interface (Unit 3)
type Importer interface {
	Import(filePath string) ([]models.Transaction, error)
	Format() string
}

// ExportImporter combines both interfaces.
// Demonstrates: Interface composition (Unit 3)
type ExportImporter interface {
	Exporter
	Importer
}
