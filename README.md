# ExpenseVault — AI-Powered Financial Command Center

ExpenseVault is a professional-grade Go application for personal finance management. It uniquely combines the speed of a **CLI**, the interactivity of a **TUI**, and the intelligence of **Google Gemini AI** to provide a comprehensive financial overview.

Whether you're tracking daily spends, setting monthly budgets, or querying your data with natural language, ExpenseVault is built for performance and scale.

---

## 🌟 Key Features

### 🤖 AI Smart Query (`ask`)
Tired of complex filters? Just ask. ExpenseVault uses **Google Gemini 2.0 Flash** to translate your natural language questions into secure, optimized SQL queries.
- *"How much did I spend on coffee last week?"*
- *"What is my biggest expense category this month?"*
- *"Did I spend more on transport in March than in April?"*

### 📊 Concurrent Analytics Engine
Built with Go's powerful concurrency primitives, the analytics engine uses a **worker pool pattern** to process large transaction histories.
- **Spending Trends**: Visual bar charts for monthly flow.
- **Predictive Budgeting**: Weighted averages to predict next month's spending.
- **Concurrent Processing**: Configurable worker counts for high-speed data crunching.

### 🎯 Budgeting & Targets
Stay disciplined by setting monthly targets across categories. The system tracks your progress and warns you when you're nearing your limits.

### 🖥️ Interactive TUI Dashboard
A beautiful, responsive terminal interface built with **BubbleTea** and **Lipgloss**. Navigate your finances without ever leaving your console.

### ☁️ Hybrid Persistence & Sync
- **SQLite**: Zero-config, pure-Go local storage for portability.
- **MySQL**: Full multi-user support with Workbench compatibility.
- **Remote Sync**: Synchronize your local data with a central server via a built-in REST API.

---

## 🚀 Quick Start

### Build from Source
```bash
# Recommended build command
go build -o expensevault .

# Launch the interactive experience
./expensevault tui
```

### Basic CLI Usage
```bash
# Set up your account
./expensevault signup -u "John"
./expensevault login -u "John"

# Add a quick expense
./expensevault add -t expense -a 1200 -c "Shopping" -d "New Sneakers"

# Natural Language Query
./expensevault ask "What were my top categories this month?"
```

## 🧪 Testing & Benchmarks

```bash
# Run the full test suite
go test ./... -v

# Run performance benchmarks for analytics & services
go test ./services/ ./analytics/ -bench=. -v
```

## 📂 Project Structure

```
expenseVault/
├── main.go                  # Entry point (Cobra execution)
├── analytics/               # Concurrent analytics engine & worker pools
├── api/                     # REST API handlers, JWT auth, and server logic
├── cmd/                     # CLI command definitions (add, ask, analytics, etc.)
├── db/                      # Persistence layer (SQLite & MySQL drivers)
├── export/                  # CSV/JSON import and export logic
├── models/                  # Core data structures and custom types
├── services/                # Business logic: Gemini AI, insights, reporting
├── tui/                     # BubbleTea-based Terminal User Interface
└── utils/                   # Shared utilities: config, helpers, logging
```

---

## 🛠️ Commands Reference

| Command | Category | Description |
|---------|----------|-------------|
| `tui` | Interface | Launch the full interactive terminal dashboard |
| `ask` | **AI** | Query your data in plain English (Gemini 2.0) |
| `analytics` | **Data** | Run trends, averages, and spending predictions |
| `add` | Core | Create a new transaction (Income/Expense) |
| `list` | Core | View and filter your transaction history |
| `edit` | Core | Modify an existing transaction |
| `delete` | Core | Remove a transaction permanently |
| `budget` | Data | Set monthly targets for specific categories |
| `report` | Data | Generate summaries (monthly, yearly, category) |
| `export` | Data | Save your data to CSV or JSON |
| `import` | Data | Bulk load data from CSV/JSON files |
| `backup` | Security | Create a comprehensive JSON snapshot |
| `restore` | Security | Reset the database from a backup file |
| `server` | DevOps | Spin up the REST API server |
| `sync` | DevOps | Synchronize data with a remote server |
| `signup` | Auth | Create a local or remote user account |
| `login` | Auth | Authenticate and obtain a JWT session |
| `demo` | Utils | Populate the vault with sample data |

---

## 🛠️ Tech Stack

- **Go 1.25.3** — High-performance systems language
- **Google Gemini 2.0 Flash** — AI-powered natural language processing
- **Cobra** — Industry-standard CLI framework
- **BubbleTea & Lipgloss** — Modern TUI development kit
- **SQLite** (`modernc.org/sqlite`) — Pure Go local storage
- **MySQL** — Multi-user enterprise-grade database support
- **bcrypt** — Secure password hashing
- **JWT** — Stateless session management
- **Goroutines & Channels** — Powering the concurrent analytics engine

---

## 🔐 TUI Auth Flow

Running `./expensevault tui` starts an integrated authentication flow:
1. **Sign up**: Locally secured with bcrypt hashing.
2. **Login**: Session validation via JWT for both local and API use.
3. **Dashboard**: Live updates and category-safe views.

---

## 🗄️ Database Setup (MySQL)

If you prefer MySQL over the default SQLite:

1. **Configure `.env`**:
   ```env
   DB_TYPE=mysql
   MYSQL_HOST=localhost
   MYSQL_PORT=3306
   MYSQL_DATABASE=expensevault
   ...
   ```
2. **Permissions**: Ensure your user has `ALL PRIVILEGES` on the `expensevault` schema.
3. **Auto-Migration**: The app handles table creation on the first run.

---

## 🎓 Go Concepts Implemented

This project serves as a comprehensive reference for modern Go development:

### Unit 1 — Basics
- `var` vs `:=`, Zero values, Custom types (`Rupees`, `Category`), Constants with `iota`, Control flow.

### Unit 2 — Composite Types
- Fixed-size Arrays, Slices (`append`, `make`, `delete`), Maps, Structs (Embedded & Anonymous).

### Unit 3 — Functions & Error Handling
- Variadic parameters, Unfurling slices, `defer`, `panic/recover`, Method sets (Value/Pointer), Interfaces & Polymorphism, Closures, Recursion, Error wrapping (`errors.Is/As`).

### Unit 4 — Advanced Patterns
- Factory functions, Pointer receivers for mutation, JSON Unmarshal/Marshal, bcrypt encryption, Table-driven testing, High-performance Benchmarking.

---

## 💻 Cross-Platform

The project runs seamlessly on **Linux, macOS, and Windows**:
- Native Windows terminal support via BubbleTea.
- Pure Go SQLite driver (no CGO/GCC toolchain required).
- Content-agnostic path handling (`filepath.Join`).
