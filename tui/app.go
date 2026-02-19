package tui

import (
	"fmt"
	"strings"

	"expenseVault/db"
	"expenseVault/models"
	"expenseVault/services"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ============================================================
// Interactive TUI using BubbleTea
// Demonstrates: Interfaces (tea.Model), methods, structs,
//               switch control flow, closures
// ============================================================

// View represents the current TUI screen
type View int

const (
	ViewDashboard View = iota
	ViewTransactions
	ViewAddForm
	ViewReports
)

// Model is the main BubbleTea model.
// Demonstrates: Struct with multiple fields (Unit 2)
type Model struct {
	store        *db.Store
	view         View
	cursor       int
	transactions []models.Transaction
	menuItems    []string
	message      string
	width        int
	height       int
	reportType   int
	quitting     bool
}

// NewModel creates a new TUI model.
func NewModel(store *db.Store) Model {
	return Model{
		store: store,
		view:  ViewDashboard,
		menuItems: []string{
			"📊 Dashboard",
			"📋 View Transactions",
			"📈 Reports",
			"🚪 Exit",
		},
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return m.loadTransactions
}

// loadTransactions fetches all transactions from DB.
func (m Model) loadTransactions() tea.Msg {
	txs, err := m.store.GetAllTransactions()
	if err != nil {
		return errMsg{err}
	}
	return txsLoadedMsg{txs}
}

// Message types
type txsLoadedMsg struct{ transactions []models.Transaction }
type errMsg struct{ err error }

// Update implements tea.Model — handles input and state changes.
// Demonstrates: switch control flow (Unit 1), methods (Unit 3)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case txsLoadedMsg:
		m.transactions = msg.transactions
		return m, nil

	case errMsg:
		m.message = fmt.Sprintf("Error: %v", msg.err)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.view == ViewDashboard {
				m.quitting = true
				return m, tea.Quit
			}
			m.view = ViewDashboard
			m.cursor = 0
			return m, nil

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			m.cursor++

		case "enter":
			return m.handleEnter()

		case "1":
			m.view = ViewDashboard
			m.cursor = 0
		case "2":
			m.view = ViewTransactions
			m.cursor = 0
		case "3":
			m.view = ViewReports
			m.cursor = 0
			m.reportType = 0

		case "tab":
			if m.view == ViewReports {
				m.reportType = (m.reportType + 1) % 3
			}

		case "esc":
			m.view = ViewDashboard
			m.cursor = 0
		}
	}

	return m, nil
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	if m.view == ViewDashboard {
		switch m.cursor {
		case 0: // Dashboard (already there)
			return m, nil
		case 1: // Transactions
			m.view = ViewTransactions
			m.cursor = 0
		case 2: // Reports
			m.view = ViewReports
			m.cursor = 0
		case 3: // Exit
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// View implements tea.Model — renders the TUI.
// Demonstrates: switch control flow, string building (Unit 1)
func (m Model) View() string {
	if m.quitting {
		return "\n  👋 Goodbye from ExpenseVault!\n\n"
	}

	var sb strings.Builder

	// Header
	sb.WriteString(TitleStyle.Render("💰 ExpenseVault — Personal Finance Dashboard"))
	sb.WriteString("\n\n")

	// Demonstrates: switch for view routing (Unit 1)
	switch m.view {
	case ViewDashboard:
		sb.WriteString(m.renderDashboard())
	case ViewTransactions:
		sb.WriteString(m.renderTransactions())
	case ViewReports:
		sb.WriteString(m.renderReports())
	}

	// Status bar
	sb.WriteString("\n")
	sb.WriteString(HelpStyle.Render("  [1]Dashboard [2]Transactions [3]Reports [q]Quit [esc]Back"))
	sb.WriteString("\n")

	if m.message != "" {
		sb.WriteString(WarningStyle.Render("  " + m.message))
		sb.WriteString("\n")
	}

	return sb.String()
}

// ============================================================
// Dashboard View
// ============================================================
func (m Model) renderDashboard() string {
	var sb strings.Builder

	// Quick summary using anonymous struct
	summary := models.QuickSummary(m.transactions)

	// Summary boxes
	incomeBox := IncomeBoxStyle.Render(fmt.Sprintf("💰 Income\n%s", summary.TotalIncome))
	expenseBox := ExpenseBoxStyle.Render(fmt.Sprintf("💸 Expenses\n%s", summary.TotalExpense))
	balanceBox := BalanceBoxStyle.Render(fmt.Sprintf("📊 Balance\n%s", summary.Balance))

	sb.WriteString("  " + lipgloss.JoinHorizontal(lipgloss.Top, incomeBox, "  ", expenseBox, "  ", balanceBox))
	sb.WriteString("\n\n")

	sb.WriteString(HeaderStyle.Render("  📌 Menu"))
	sb.WriteString("\n\n")

	// Menu items with cursor
	for i, item := range m.menuItems {
		if i == m.cursor {
			sb.WriteString(SelectedMenuStyle.Render("▸ " + item))
		} else {
			sb.WriteString(MenuItemStyle.Render("  " + item))
		}
		sb.WriteString("\n")
	}

	// Recent transactions
	sb.WriteString("\n")
	sb.WriteString(HeaderStyle.Render("  📝 Recent Transactions"))
	sb.WriteString("\n\n")

	count := 5
	if len(m.transactions) < count {
		count = len(m.transactions)
	}

	if count == 0 {
		sb.WriteString(MutedStyle.Render("  No transactions yet. Use 'expensevault add' to add some!"))
		sb.WriteString("\n")
	} else {
		for i := 0; i < count; i++ {
			tx := m.transactions[i]
			style := IncomeStyle
			if tx.Type == models.Expense {
				style = ExpenseStyle
			}
			sb.WriteString(fmt.Sprintf("  %s  %s\n",
				style.Render(fmt.Sprintf("%-8s %s", tx.Amount, tx.Type)),
				MutedStyle.Render(fmt.Sprintf("%-15s %s", tx.Category, tx.Date)),
			))
		}
	}

	return sb.String()
}

// ============================================================
// Transactions List View
// ============================================================
func (m Model) renderTransactions() string {
	var sb strings.Builder

	sb.WriteString(HeaderStyle.Render("  📋 All Transactions"))
	sb.WriteString("\n\n")

	if len(m.transactions) == 0 {
		sb.WriteString(MutedStyle.Render("  No transactions found."))
		sb.WriteString("\n")
		return sb.String()
	}

	// Header row
	sb.WriteString(TableHeaderStyle.Render(
		fmt.Sprintf("  %-4s %-8s %-12s %-15s %-25s %-12s", "ID", "Type", "Amount", "Category", "Description", "Date"),
	))
	sb.WriteString("\n")

	// Demonstrates: for range with index (Unit 2)
	displayCount := len(m.transactions)
	if displayCount > 20 {
		displayCount = 20
	}

	for i := 0; i < displayCount; i++ {
		tx := m.transactions[i]
		style := TableRowStyle
		if i == m.cursor {
			style = style.Bold(true).Foreground(lipgloss.Color("#FFFFFF"))
		}

		symbol := "💰"
		if tx.Type == models.Expense {
			symbol = "💸"
		}

		row := fmt.Sprintf("  %-4d %s%-6s %-12s %-15s %-25s %-12s",
			tx.ID, symbol, tx.Type, tx.Amount, tx.Category, tx.Description, tx.Date)
		sb.WriteString(style.Render(row))
		sb.WriteString("\n")
	}

	if len(m.transactions) > 20 {
		sb.WriteString(MutedStyle.Render(fmt.Sprintf("\n  ... and %d more transactions", len(m.transactions)-20)))
		sb.WriteString("\n")
	}

	return sb.String()
}

// ============================================================
// Reports View — uses Reporter interface polymorphism
// ============================================================
func (m Model) renderReports() string {
	var sb strings.Builder

	reportTypes := []string{"Monthly", "Category", "Yearly"}
	sb.WriteString(HeaderStyle.Render("  📈 Reports"))
	sb.WriteString("  ")
	for i, rt := range reportTypes {
		if i == m.reportType {
			sb.WriteString(SelectedMenuStyle.Render("[" + rt + "]"))
		} else {
			sb.WriteString(MutedStyle.Render(" " + rt + " "))
		}
		sb.WriteString(" ")
	}
	sb.WriteString(MutedStyle.Render("  (press Tab to switch)"))
	sb.WriteString("\n\n")

	if len(m.transactions) == 0 {
		sb.WriteString(MutedStyle.Render("  No data for reports."))
		sb.WriteString("\n")
		return sb.String()
	}

	// Demonstrates: Interface polymorphism — select reporter at runtime (Unit 3)
	var reporter services.Reporter
	switch m.reportType {
	case 0:
		reporter = &services.MonthlyReporter{}
	case 1:
		reporter = &services.CategoryReporter{}
	case 2:
		reporter = &services.YearlyReporter{}
	}

	sb.WriteString(reporter.Generate(m.transactions))
	return sb.String()
}

// RunTUI starts the BubbleTea TUI application.
func RunTUI(store *db.Store) error {
	model := NewModel(store)
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
