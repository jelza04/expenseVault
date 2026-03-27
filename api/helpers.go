package api

import (
	"encoding/json"
	"net/http"

	"expenseVault/models"
)

// ── Utility Helpers ─────────────────────────────────────────

// writeJSON encodes the given payload as JSON and writes it to the response.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// handleHealth is a simple health-check endpoint.
func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleSync handles POST /api/sync (existing functionality preserved).
func handleSync(w http.ResponseWriter, r *http.Request) {
	payload := &models.SyncPayload{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	_, marshalErr := models.MarshalTransactions(payload.Transactions)

	writeJSON(w, http.StatusOK, map[string]any{
		"status":     "ok",
		"received":   len(payload.Transactions),
		"marshal_ok": marshalErr == nil,
	})
}
