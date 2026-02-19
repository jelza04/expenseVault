package tui

import "github.com/charmbracelet/lipgloss"

// ============================================================
// TUI Styles using lipgloss for beautiful terminal rendering
// ============================================================

var (
	// Color palette
	primaryColor   = lipgloss.Color("#7C3AED")
	secondaryColor = lipgloss.Color("#10B981")
	dangerColor    = lipgloss.Color("#EF4444")
	warningColor   = lipgloss.Color("#F59E0B")
	infoColor      = lipgloss.Color("#3B82F6")
	mutedColor     = lipgloss.Color("#6B7280")
	bgColor        = lipgloss.Color("#1F2937")

	// Title style
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(primaryColor).
			Padding(0, 2).
			MarginBottom(1)

	// Menu item styles
	MenuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	SelectedMenuStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(primaryColor).
				PaddingLeft(1).
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(primaryColor)

	// Box styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	IncomeBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(0, 2)

	ExpenseBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(dangerColor).
			Padding(0, 2)

	BalanceBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(infoColor).
			Padding(0, 2)

	// Text styles
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	IncomeStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	ExpenseStyle = lipgloss.NewStyle().
			Foreground(dangerColor)

	MutedStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	WarningStyle = lipgloss.NewStyle().
			Foreground(warningColor)

	// Help style
	HelpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	// Status bar
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#374151")).
			Padding(0, 1).
			Width(60)

	// Table styles
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(primaryColor).
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(mutedColor)

	TableRowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D5DB"))

	_ = bgColor // prevent unused warning
	_ = warningColor
)
