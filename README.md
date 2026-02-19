# 💰 ExpenseVault — Privacy-Focused Personal Finance CLI

A comprehensive Go CLI application for managing personal finances with interactive TUI, local SQLite storage, CSV/JSON import/export, and optional backend sync.

## 🚀 Quick Start

```bash
# Build
cd expenseVault
go build -o expensevault .

# Add transactions
./expensevault add -t expense -a 250 -d "Lunch at canteen" -c Food
./expensevault add -t income -a 50000 -d "Monthly salary" -c Salary
./expensevault add -t expense -a 800 -d "Uber ride" --auto-cat

# List & manage
./expensevault list
./expensevault list --type expense --category Food
./expensevault edit --id 1 --amount 300
./expensevault delete --id 2

# Reports
./expensevault report --type monthly
./expensevault report --type category
./expensevault report --all

# Import/Export
./expensevault export -f csv -o data.csv
./expensevault export -f json -o data.json
./expensevault import -f bank_statement.csv

# Backup/Restore
./expensevault backup
./expensevault restore -f expensevault_backup.json

# Interactive TUI Dashboard
./expensevault tui

# Demo all Go concepts
./expensevault demo

# Backend Sync Server
./expensevault server --port 8080
./expensevault sync --server http://localhost:8080 --token <jwt>
```

## 📁 Project Structure

```
expenseVault/
├── main.go                     # Entry point
├── cmd/                        # Cobra CLI commands
│   ├── root.go                 # Root command + DB init
│   ├── add.go, list.go         # CRUD commands
│   ├── edit.go, delete.go
│   ├── report.go               # Interface-based reports
│   ├── importcmd.go, exportcmd.go
│   ├── backup.go               # JSON backup/restore
│   ├── tui.go, server.go, sync.go
│   └── demo.go                 # Syllabus demo
├── models/                     # Types, structs, errors
│   ├── types.go                # Custom types, zero values
│   ├── transaction.go          # Structs, embedding
│   └── errors.go               # Custom error types
├── db/sqlite.go                # SQLite CRUD + defer
├── services/                   # Business logic
│   ├── transaction.go          # Variadic, closures, recursion
│   ├── categorizer.go          # Callbacks, rules
│   └── reporter.go             # Interface polymorphism
├── export/                     # Import/Export
│   ├── exporter.go             # Exporter/Importer interfaces
│   ├── csv.go                  # CSV handler
│   └── json_export.go          # JSON handler
├── tui/                        # BubbleTea TUI
│   ├── app.go                  # Interactive dashboard
│   └── styles.go               # Lipgloss styles
├── api/                        # REST API
│   ├── server.go               # HTTP server + handlers
│   └── auth.go                 # JWT authentication
└── utils/helpers.go            # Utilities
```

## 🎓 Syllabus Concepts Mapping

### Unit 1 — Programming Fundamentals
| Concept | Location |
|---|---|
| Custom types, `type` keyword | `models/types.go` — `Rupees`, `Category` |
| Short declaration `:=` | Throughout all `cmd/*.go` files |
| `var` keyword & zero values | `models/types.go` — `ZeroValueDemo()` |
| `fmt` package | Used everywhere for output |
| Type conversion (not casting) | `models/types.go` — `ToFloat64()` |
| `for` loop, `for range` | `services/transaction.go`, `db/sqlite.go` |
| `if-else`, `switch` | `cmd/add.go`, `services/reporter.go` |
| Packages & imports | Every file demonstrates package system |

### Unit 2 — Grouping Data
| Concept | Location |
|---|---|
| Arrays (fixed-size) | `models/types.go` — `DefaultExpenseCategories` |
| Slices, `append`, `make` | `services/transaction.go` |
| Delete from slice | `services/transaction.go` — `DeleteFromSlice()` |
| Composite literals | `models/transaction.go` — `SampleTransactions()` |
| Multi-dimensional slices | `services/transaction.go` — `MonthlyTotals()` |
| Maps (add, range, delete) | `services/reporter.go`, `cmd/demo.go` |
| Structs with tags | `models/transaction.go` — `Transaction` |
| Embedded structs | `models/transaction.go` — `DateInfo` embedded |
| Anonymous structs | `models/transaction.go` — `QuickSummary()` |

### Unit 3 — Functions
| Concept | Location |
|---|---|
| Variadic parameters | `services/transaction.go` — `AddMultiple()` |
| Unfurling a slice | `services/transaction.go` — `AddFromSlice()` |
| `defer` | `db/sqlite.go`, `export/csv.go` |
| `panic` / `recover` | `cmd/demo.go` — `deferPanicDemo()` |
| Methods | All struct types have methods |
| Interfaces & polymorphism | `services/reporter.go` — `Reporter` interface |
| Anonymous functions | `services/transaction.go` — `GetFilterFunc()` |
| Function expressions | `services/transaction.go` |
| Returning a function | `services/categorizer.go` — `MakeCategorizer()` |
| Callbacks | `services/transaction.go` — `ForEachTransaction()` |
| Closures | `services/transaction.go` — `MakeBalanceTracker()` |
| Recursion | `services/transaction.go` — `RecursiveCategorySum()` |
| Custom error types | `models/errors.go` |
| Error checking & logging | `db/sqlite.go`, all `cmd/*.go` |

## 🛠️ Tech Stack

| Technology | Purpose |
|---|---|
| Go (Golang) | Core language |
| Cobra | CLI framework |
| BubbleTea + Lipgloss | Interactive TUI |
| SQLite (modernc.org) | Local database |
| JWT | API authentication |
| bcrypt | Password hashing |
| encoding/csv, encoding/json | Data import/export |

## 🧪 Running Tests

```bash
go test ./models/ ./services/ -v
```
