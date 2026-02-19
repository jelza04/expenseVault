package db

import (
	"database/sql"
	"fmt"
	"time"

	"expenseVault/models"
)

// ============================================================
// UNIT 3: Defer for resource cleanup, Error Handling
// UNIT 1: Control flow, conditionals
// ============================================================

// Store manages all database operations with SQLite.
type Store struct {
	db *sql.DB
}

// NewStore opens/creates the SQLite database and initializes tables.
// Demonstrates: Error handling, defer for cleanup on failure
func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, &models.DatabaseError{Operation: "open", Err: err}
	}

	store := &Store{db: db}

	if err := store.createTables(); err != nil {
		// Demonstrates: defer-like cleanup on error path
		db.Close()
		return nil, err
	}

	return store, nil
}

// Close closes the database connection.
// Demonstrates: Cleanup method (typically called with defer)
func (s *Store) Close() error {
	return s.db.Close()
}

// createTables sets up the schema.
func (s *Store) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS transactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL,
			amount REAL NOT NULL,
			category TEXT NOT NULL,
			description TEXT NOT NULL,
			date TEXT NOT NULL,
			notes TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_tx_date ON transactions(date)`,
		`CREATE INDEX IF NOT EXISTS idx_tx_category ON transactions(category)`,
		`CREATE INDEX IF NOT EXISTS idx_tx_type ON transactions(type)`,
	}

	// Demonstrates: for range over slice (Unit 1: looping)
	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return &models.DatabaseError{Operation: "create_tables", Err: err}
		}
	}
	return nil
}

// ============================================================
// CRUD Operations
// ============================================================

// AddTransaction inserts a new transaction.
// Demonstrates: Error checking, returning values
func (s *Store) AddTransaction(t models.Transaction) (int64, error) {
	// Demonstrates: Input validation (Unit 3: checking errors)
	if err := models.ValidateTransaction(t); err != nil {
		return 0, err
	}

	result, err := s.db.Exec(
		`INSERT INTO transactions (type, amount, category, description, date, notes)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		t.Type, t.Amount.ToFloat64(), t.Category, t.Description, t.Date, t.Notes,
	)
	if err != nil {
		return 0, &models.DatabaseError{Operation: "insert", Err: err}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, &models.DatabaseError{Operation: "last_insert_id", Err: err}
	}
	return id, nil
}

// GetTransaction retrieves a single transaction by ID.
// Demonstrates: Error handling with NotFoundError
func (s *Store) GetTransaction(id int) (*models.Transaction, error) {
	row := s.db.QueryRow(
		`SELECT id, type, amount, category, description, date, notes, created_at, updated_at
		 FROM transactions WHERE id = ?`, id,
	)

	var t models.Transaction
	var amount float64
	err := row.Scan(&t.ID, &t.Type, &amount, &t.Category,
		&t.Description, &t.Date, &t.Notes, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, &models.NotFoundError{Resource: "Transaction", ID: id}
	}
	if err != nil {
		return nil, &models.DatabaseError{Operation: "get", Err: err}
	}
	t.Amount = models.Rupees(amount)
	return &t, nil
}

// ListTransactions retrieves transactions with optional filters.
// Demonstrates: Slice building, conditional query construction (Unit 1: control flow)
func (s *Store) ListTransactions(txType string, category string, startDate string, endDate string, limit int) ([]models.Transaction, error) {
	query := "SELECT id, type, amount, category, description, date, notes, created_at, updated_at FROM transactions WHERE 1=1"
	// Demonstrates: Slice (Unit 2), append (Unit 2)
	args := make([]interface{}, 0)

	// Demonstrates: Conditional logic (Unit 1: if-else)
	if txType != "" {
		query += " AND type = ?"
		args = append(args, txType)
	}
	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}
	if startDate != "" {
		query += " AND date >= ?"
		args = append(args, startDate)
	}
	if endDate != "" {
		query += " AND date <= ?"
		args = append(args, endDate)
	}

	query += " ORDER BY date DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, &models.DatabaseError{Operation: "list", Err: err}
	}
	// Demonstrates: defer for resource cleanup (Unit 3)
	defer rows.Close()

	// Demonstrates: make to pre-allocate slice (Unit 2)
	transactions := make([]models.Transaction, 0)

	// Demonstrates: for loop with conditional (Unit 1)
	for rows.Next() {
		var t models.Transaction
		var amount float64
		if err := rows.Scan(&t.ID, &t.Type, &amount, &t.Category,
			&t.Description, &t.Date, &t.Notes, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, &models.DatabaseError{Operation: "scan", Err: err}
		}
		t.Amount = models.Rupees(amount)
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, &models.DatabaseError{Operation: "rows_iteration", Err: err}
	}

	return transactions, nil
}

// UpdateTransaction modifies an existing transaction.
func (s *Store) UpdateTransaction(t models.Transaction) error {
	if err := models.ValidateTransaction(t); err != nil {
		return err
	}

	result, err := s.db.Exec(
		`UPDATE transactions SET type=?, amount=?, category=?, description=?, date=?, notes=?, updated_at=?
		 WHERE id=?`,
		t.Type, t.Amount.ToFloat64(), t.Category, t.Description,
		t.Date, t.Notes, time.Now(), t.ID,
	)
	if err != nil {
		return &models.DatabaseError{Operation: "update", Err: err}
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return &models.NotFoundError{Resource: "Transaction", ID: t.ID}
	}
	return nil
}

// DeleteTransaction removes a transaction by ID.
func (s *Store) DeleteTransaction(id int) error {
	result, err := s.db.Exec("DELETE FROM transactions WHERE id = ?", id)
	if err != nil {
		return &models.DatabaseError{Operation: "delete", Err: err}
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return &models.NotFoundError{Resource: "Transaction", ID: id}
	}
	return nil
}

// GetAllTransactions returns all transactions (used for export/backup).
func (s *Store) GetAllTransactions() ([]models.Transaction, error) {
	return s.ListTransactions("", "", "", "", 0)
}

// BulkInsert adds multiple transactions in a single DB transaction.
// Demonstrates: Database transactions, defer for rollback, error handling
func (s *Store) BulkInsert(transactions []models.Transaction) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, &models.DatabaseError{Operation: "begin_tx", Err: err}
	}
	// Demonstrates: defer with recover pattern (Unit 3)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(
		`INSERT INTO transactions (type, amount, category, description, date, notes)
		 VALUES (?, ?, ?, ?, ?, ?)`,
	)
	if err != nil {
		tx.Rollback()
		return 0, &models.DatabaseError{Operation: "prepare", Err: err}
	}
	defer stmt.Close()

	count := 0
	for _, t := range transactions {
		if err := models.ValidateTransaction(t); err != nil {
			continue // Skip invalid transactions
		}
		if _, err := stmt.Exec(t.Type, t.Amount.ToFloat64(), t.Category,
			t.Description, t.Date, t.Notes); err != nil {
			tx.Rollback()
			return count, &models.DatabaseError{Operation: "bulk_insert", Err: err}
		}
		count++
	}

	if err := tx.Commit(); err != nil {
		return count, &models.DatabaseError{Operation: "commit", Err: err}
	}
	return count, nil
}

// ============================================================
// User operations (for backend sync/auth)
// ============================================================

// CreateUser adds a new user.
func (s *Store) CreateUser(username, passwordHash string) (int64, error) {
	result, err := s.db.Exec(
		"INSERT INTO users (username, password_hash) VALUES (?, ?)",
		username, passwordHash,
	)
	if err != nil {
		return 0, &models.DatabaseError{Operation: "create_user", Err: err}
	}
	return result.LastInsertId()
}

// GetUserByUsername retrieves a user by username.
func (s *Store) GetUserByUsername(username string) (*models.User, error) {
	row := s.db.QueryRow(
		"SELECT id, username, password_hash, created_at FROM users WHERE username = ?",
		username,
	)
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, &models.AuthError{Reason: "user not found"}
	}
	if err != nil {
		return nil, &models.DatabaseError{Operation: "get_user", Err: err}
	}
	return &u, nil
}
