package main

import "expenseVault/cmd"

// ============================================================
// ExpenseVault — Privacy-Focused Personal Finance CLI Tool
//
// Entry point: delegates to Cobra command framework
// ============================================================

func main() {
	cmd.Execute()
}
