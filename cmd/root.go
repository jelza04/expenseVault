package cmd

import (
	"fmt"
	"os"

	"expenseVault/db"
	"expenseVault/utils"

	"github.com/spf13/cobra"

	// SQLite driver (pure Go, no CGO required)
	_ "modernc.org/sqlite"
)

var store *db.Store

// rootCmd is the base command for ExpenseVault CLI.
var rootCmd = &cobra.Command{
	Use:   "expensevault",
	Short: "💰 ExpenseVault — Privacy-Focused Personal Finance CLI",
	Long: `
╔══════════════════════════════════════════════════════════════╗
║          💰 ExpenseVault — Personal Finance CLI             ║
║                                                              ║
║  A privacy-focused CLI tool for managing personal finances.  ║
║  Track income & expenses, generate reports, import/export    ║
║  data, and sync across devices — all from your terminal.     ║
╚══════════════════════════════════════════════════════════════╝`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip DB init for server command (uses its own DB)
		if cmd.Name() == "server" {
			return nil
		}
		var err error
		store, err = db.NewStore(utils.GetDBPath())
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if store != nil {
			store.Close()
		}
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Add all subcommands
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(tuiCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(demoCmd)
}
