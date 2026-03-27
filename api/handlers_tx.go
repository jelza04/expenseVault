package api

import (
	"encoding/json"
	"net/http"
	"time"

	"expenseVault/db"
	"expenseVault/models"
	"expenseVault/services"
)

// ── Request types ───────────────────────────────────────────

type addTransactionRequest struct {
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	Notes       string  `json:"notes"`
}

// ── Handlers ────────────────────────────────────────────────

// makeListTransactionsHandler returns a handler for GET /api/transactions.
// Protected by JWT — uses user_id from context.
func makeListTransactionsHandler(store *db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserIDFromContext(r)

		txs, err := store.GetAllTransactions(userID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"transactions": txs,
			"count":        len(txs),
		})
	}
}

// makeAddTransactionHandler returns a handler for POST /api/transactions.
// Protected by JWT — uses user_id from context.
func makeAddTransactionHandler(store *db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserIDFromContext(r)

		var req addTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		if req.Type != "income" && req.Type != "expense" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "type must be 'income' or 'expense'"})
			return
		}
		if req.Amount <= 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "amount must be positive"})
			return
		}
		if req.Description == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "description is required"})
			return
		}
		if req.Date == "" {
			req.Date = time.Now().Format("2006-01-02")
		}

		cat := models.Category(req.Category)
		if req.Category == "" {
			cat = services.NewCategorizer().AutoCategorize(req.Description)
		}

		tx := models.NewTransaction(
			userID,
			models.TransactionType(req.Type),
			req.Amount,
			cat,
			req.Description,
			req.Date,
		)
		if req.Notes != "" {
			tx.SetNotes(req.Notes)
		}

		id, err := store.AddTransaction(tx)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusCreated, map[string]any{
			"message": "transaction added",
			"id":      id,
		})
	}
}
