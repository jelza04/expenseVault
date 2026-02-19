package models

import "fmt"

// ============================================================
// UNIT 3: Error Handling — Custom Error Types, Errors with Info
// Demonstrates: error interface, custom error types, fmt.Errorf,
//               errors with additional context/info
// ============================================================

// ValidationError is returned when input validation fails.
// Demonstrates: Custom error type implementing the error interface
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error [%s]: %s", e.Field, e.Message)
}

// NotFoundError is returned when a requested resource doesn't exist.
// Demonstrates: Custom error type with context info
type NotFoundError struct {
	Resource string
	ID       int
}

// Error implements the error interface.
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %d not found", e.Resource, e.ID)
}

// DatabaseError wraps database-related errors with additional context.
// Demonstrates: Error wrapping, errors with info
type DatabaseError struct {
	Operation string
	Err       error
}

// Error implements the error interface.
func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s: %v", e.Operation, e.Err)
}

// Unwrap returns the underlying error for errors.Is/As support.
func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// ParseError is returned when data parsing fails.
// Demonstrates: Custom error with multiple info fields
type ParseError struct {
	Input    string
	Expected string
	Line     int
}

// Error implements the error interface.
func (e *ParseError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("parse error at line %d: expected %s, got %q", e.Line, e.Expected, e.Input)
	}
	return fmt.Sprintf("parse error: expected %s, got %q", e.Expected, e.Input)
}

// AuthError is returned for authentication failures.
type AuthError struct {
	Reason string
}

// Error implements the error interface.
func (e *AuthError) Error() string {
	return fmt.Sprintf("authentication error: %s", e.Reason)
}

// ============================================================
// Sentinel Errors — common predefined errors
// ============================================================

// Common sentinel errors used throughout the application.
var (
	ErrInvalidAmount   = &ValidationError{Field: "amount", Message: "amount must be positive"}
	ErrInvalidDate     = &ValidationError{Field: "date", Message: "date must be in YYYY-MM-DD format"}
	ErrInvalidCategory = &ValidationError{Field: "category", Message: "unknown category"}
	ErrInvalidType     = &ValidationError{Field: "type", Message: "must be 'income' or 'expense'"}
	ErrEmptyDesc       = &ValidationError{Field: "description", Message: "description cannot be empty"}
)

// ValidateTransaction checks a transaction for errors.
// Demonstrates: Comprehensive error checking (Unit 3)
func ValidateTransaction(t Transaction) error {
	if t.Amount <= 0 {
		return ErrInvalidAmount
	}
	if t.Description == "" {
		return ErrEmptyDesc
	}
	if t.Type != Income && t.Type != Expense {
		return ErrInvalidType
	}
	if t.Date == "" {
		return ErrInvalidDate
	}
	// Validate date format
	if _, err := t.ParsedDate(); err != nil {
		return &ValidationError{
			Field:   "date",
			Message: fmt.Sprintf("invalid date format: %v", err),
		}
	}
	// Validate category
	valid := false
	for _, cat := range CategoriesForType(t.Type) {
		if cat == t.Category {
			valid = true
			break
		}
	}
	if !valid {
		return &ValidationError{
			Field:   "category",
			Message: fmt.Sprintf("'%s' is not a valid category for %s", t.Category, t.Type),
		}
	}
	return nil
}
