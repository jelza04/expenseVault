package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"expenseVault/db"
	"expenseVault/models"
	"expenseVault/services"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ============================================================
// Interactive TUI using BubbleTea
// Now includes: Dashboard, Transactions, Add Form, Reports
// ============================================================

// View represents the current TUI screen
type View int

const (
	ViewDashboard View = iota
	ViewTransactions
	ViewAddForm
	ViewReports
)

// Form field indices
const (
	fieldType = iota
	fieldAmount
	fieldCategory
	fieldDescription
	fieldDate
	fieldNotes
	fieldCount
)

// Model is the main BubbleTea model.
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

	// Form fields for adding transactions
	inputs      []textinput.Model
	focusIndex  int
	formMessage string
}

// NewModel creates a new TUI model.
func NewModel(store *db.Store) Model {
	// Create text input fields for the Add Transaction form
	inputs := make([]textinput.Model, fieldCount)

	inputs[fieldType] = textinput.New()
	inputs[fieldType].Placeholder = "expense or income"
	inputs[fieldType].CharLimit = 10
	inputs[fieldType].Width = 30
	inputs[fieldType].Prompt = "  Type ▸ "

	inputs[fieldAmount] = textinput.New()
	inputs[fieldAmount].Placeholder = "e.g. 250.00"
	inputs[fieldAmount].CharLimit = 15
	inputs[fieldAmount].Width = 30
	inputs[fieldAmount].Prompt = "  Amount ▸ "

	inputs[fieldCategory] = textinput.New()
	inputs[fieldCategory].Placeholder = "Food, Travel, Shopping, Bills..."
	inputs[fieldCategory].CharLimit = 20
	inputs[fieldCategory].Width = 30
	inputs[fieldCategory].Prompt = "  Category ▸ "

	inputs[fieldDescription] = textinput.New()
	inputs[fieldDescription].Placeholder = "e.g. Lunch at canteen"
	inputs[fieldDescription].CharLimit = 50
	inputs[fieldDescription].Width = 40
	inputs[fieldDescription].Prompt = "  Description ▸ "

	inputs[fieldDate] = textinput.New()
	inputs[fieldDate].Placeholder = time.Now().Format("2006-01-02")
	inputs[fieldDate].CharLimit = 10
	inputs[fieldDate].Width = 30
	inputs[fieldDate].Prompt = "  Date ▸ "

	inputs[fieldNotes] = textinput.New()
	inputs[fieldNotes].Placeholder = "(optional)"
	inputs[fieldNotes].CharLimit = 50
	inputs[fieldNotes].Width = 40
	inputs[fieldNotes].Prompt = "  Notes ▸ "

	// Focus the first input
	inputs[fieldType].Focus()

	return Model{
		store: store,
		view:  ViewDashboard,
		menuItems: []string{
			"📊 Dashboard",
			"📋 View Transactions",
			"➕ Add Transaction",
			"📈 Reports",
			"🚪 Exit",
		},
		inputs:     inputs,
		focusIndex: 0,
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
type txAddedMsg struct{ id int64 }

// Update implements tea.Model — handles input and state changes.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle form view separately since it captures key input
	if m.view == ViewAddForm {
		return m.updateForm(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case txsLoadedMsg:
		m.transactions = msg.transactions
		m.message = ""
		return m, nil

	case txAddedMsg:
		m.message = fmt.Sprintf("✅ Transaction #%d added!", msg.id)
		m.view = ViewDashboard
		m.cursor = 0
		return m, m.loadTransactions

	case errMsg:
		m.message = fmt.Sprintf("❌ Error: %v", msg.err)
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
			m.message = ""
			return m, nil

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			max := len(m.menuItems) - 1
			if m.view == ViewTransactions {
				max = len(m.transactions) - 1
				if max > 19 {
					max = 19
				}
			}
			if m.cursor < max {
				m.cursor++
			}

		case "enter":
			return m.handleEnter()

		case "1":
			m.view = ViewDashboard
			m.cursor = 0
		case "2":
			m.view = ViewTransactions
			m.cursor = 0
		case "3":
			m.view = ViewAddForm
			m.resetForm()
			return m, m.inputs[0].Focus()
		case "4":
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
			m.message = ""
		}
	}

	return m, nil
}

// updateForm handles input when the Add Form is visible.
func (m Model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case txsLoadedMsg:
		m.transactions = msg.transactions
		return m, nil

	case txAddedMsg:
		m.message = fmt.Sprintf("✅ Transaction #%d added successfully!", msg.id)
		m.view = ViewDashboard
		m.cursor = 0
		return m, m.loadTransactions

	case errMsg:
		m.formMessage = fmt.Sprintf("❌ %v", msg.err)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "esc":
			m.view = ViewDashboard
			m.cursor = 0
			m.formMessage = ""
			return m, nil

		case "tab", "down":
			// Move to next field
			m.focusIndex = (m.focusIndex + 1) % fieldCount
			return m, m.updateFocus()

		case "shift+tab", "up":
			// Move to previous field
			m.focusIndex = (m.focusIndex - 1 + fieldCount) % fieldCount
			return m, m.updateFocus()

		case "enter":
			// If on last field or pressing enter, try to submit
			if m.focusIndex == fieldCount-1 {
				return m.submitForm()
			}
			// Otherwise move to next field
			m.focusIndex = (m.focusIndex + 1) % fieldCount
			return m, m.updateFocus()

		case "ctrl+s":
			// Submit form with Ctrl+S from any field
			return m.submitForm()
		}
	}

	// Update the focused text input
	var cmd tea.Cmd
	m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
	return m, cmd
}

// updateFocus sets focus to the current focus index field.
func (m *Model) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, fieldCount)
	for i := range m.inputs {
		if i == m.focusIndex {
			cmds[i] = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
	return tea.Batch(cmds...)
}

// resetForm clears all form fields and resets focus.
func (m *Model) resetForm() {
	for i := range m.inputs {
		m.inputs[i].SetValue("")
		m.inputs[i].Blur()
	}
	m.focusIndex = 0
	m.formMessage = ""
	m.inputs[fieldDate].SetValue(time.Now().Format("2006-01-02"))
}

// submitForm validates and saves the transaction.
func (m Model) submitForm() (tea.Model, tea.Cmd) {
	txType := strings.TrimSpace(strings.ToLower(m.inputs[fieldType].Value()))
	amountStr := strings.TrimSpace(m.inputs[fieldAmount].Value())
	category := strings.TrimSpace(m.inputs[fieldCategory].Value())
	desc := strings.TrimSpace(m.inputs[fieldDescription].Value())
	date := strings.TrimSpace(m.inputs[fieldDate].Value())
	notes := strings.TrimSpace(m.inputs[fieldNotes].Value())

	// Validate type
	if txType != "income" && txType != "expense" {
		m.formMessage = "⚠️  Type must be 'income' or 'expense'"
		return m, nil
	}

	// Validate amount
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		m.formMessage = "⚠️  Amount must be a positive number"
		return m, nil
	}

	// Validate description
	if desc == "" {
		m.formMessage = "⚠️  Description is required"
		return m, nil
	}

	// Default date to today
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	// Auto-categorize if category is empty
	var cat models.Category
	if category == "" {
		categorizer := services.NewCategorizer()
		cat = categorizer.AutoCategorize(desc)
	} else {
		cat = models.Category(category)
	}

	// Create and save transaction
	tx := models.Transaction{
		Type:        models.TransactionType(txType),
		Amount:      models.Rupees(amount),
		Category:    cat,
		Description: desc,
		Date:        date,
		Notes:       notes,
	}

	store := m.store
	return m, func() tea.Msg {
		id, err := store.AddTransaction(tx)
		if err != nil {
			return errMsg{err}
		}
		return txAddedMsg{id}
	}
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	if m.view == ViewDashboard {
		switch m.cursor {
		case 0: // Dashboard (already there)
			return m, nil
		case 1: // Transactions
			m.view = ViewTransactions
			m.cursor = 0
		case 2: // Add Transaction
			m.view = ViewAddForm
			m.resetForm()
			return m, m.inputs[0].Focus()
		case 3: // Reports
			m.view = ViewReports
			m.cursor = 0
		case 4: // Exit
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// View implements tea.Model — renders the TUI.
func (m Model) View() string {
	if m.quitting {
		return "\n  👋 Goodbye from ExpenseVault!\n\n"
	}

	var sb strings.Builder

	// Header
	sb.WriteString(TitleStyle.Render("💰 ExpenseVault — Personal Finance Dashboard"))
	sb.WriteString("\n\n")

	switch m.view {
	case ViewDashboard:
		sb.WriteString(m.renderDashboard())
	case ViewTransactions:
		sb.WriteString(m.renderTransactions())
	case ViewAddForm:
		sb.WriteString(m.renderAddForm())
	case ViewReports:
		sb.WriteString(m.renderReports())
	}

	// Status bar
	sb.WriteString("\n")
	if m.view == ViewAddForm {
		sb.WriteString(HelpStyle.Render("  [Tab/↓]Next field  [Shift+Tab/↑]Prev field  [Ctrl+S]Save  [Esc]Cancel"))
	} else {
		sb.WriteString(HelpStyle.Render("  [1]Dashboard [2]Transactions [3]Add [4]Reports [q]Quit [esc]Back"))
	}
	sb.WriteString("\n")

	if m.message != "" {
		sb.WriteString(WarningStyle.Render("  " + m.message))
		sb.WriteString("\n")
	}

	return sb.String()
}

// ============================================================
// Add Transaction Form View
// ============================================================
func (m Model) renderAddForm() string {
	var sb strings.Builder

	sb.WriteString(HeaderStyle.Render("  ➕ Add New Transaction"))
	sb.WriteString("\n\n")

	sb.WriteString(BoxStyle.Render(
		fmt.Sprintf(
			"Fill in the details below and press Ctrl+S to save.\n\n"+
				"%s\n\n"+
				"%s\n\n"+
				"%s\n\n"+
				"%s\n\n"+
				"%s\n\n"+
				"%s",
			m.fieldView(fieldType, "Type", "income / expense"),
			m.fieldView(fieldAmount, "Amount", "in ₹"),
			m.fieldView(fieldCategory, "Category", "leave empty for auto-detect"),
			m.fieldView(fieldDescription, "Description", "what was it for?"),
			m.fieldView(fieldDate, "Date", "YYYY-MM-DD"),
			m.fieldView(fieldNotes, "Notes", "optional extra info"),
		),
	))

	// Available categories hint
	sb.WriteString("\n\n")
	sb.WriteString(MutedStyle.Render("  📂 Categories: Food, Travel, Shopping, Bills, Health, Education, Entertainment, Salary, Freelance, Other"))
	sb.WriteString("\n")

	if m.formMessage != "" {
		sb.WriteString("\n")
		sb.WriteString(WarningStyle.Render("  " + m.formMessage))
		sb.WriteString("\n")
	}

	return sb.String()
}

// fieldView renders a single form field with its label.
func (m Model) fieldView(index int, label string, hint string) string {
	focusIndicator := "  "
	style := MutedStyle
	if m.focusIndex == index {
		focusIndicator = "▸ "
		style = HeaderStyle
	}

	return fmt.Sprintf("%s%s %s\n%s",
		focusIndicator,
		style.Render(label),
		MutedStyle.Render("("+hint+")"),
		m.inputs[index].View(),
	)
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
		sb.WriteString(MutedStyle.Render("  No transactions yet. Select '➕ Add Transaction' to get started!"))
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
