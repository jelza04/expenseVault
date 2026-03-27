
package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"expenseVault/db"
	"expenseVault/services"
	"expenseVault/utils"

	"golang.org/x/sync/singleflight"
)

// aiGroup is a singleflight group for coalescing concurrent AI requests.
var aiGroup singleflight.Group

// ── Request / Response types ────────────────────────────────

type askRequest struct {
	Query string `json:"query"`
}

type askResponse struct {
	Answer string `json:"answer"`
}

// ── Handler ─────────────────────────────────────────────────

// makeAskHandler returns a handler for POST /api/ask.
// Wraps the Gemini LLM integration with caching and singleflight support.
func makeAskHandler(store *db.Store, cfg *utils.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserIDFromContext(r)

		var req askRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		if req.Query == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "query is required"})
			return
		}

		// 1. Check Backend Cache First
		if answer, ok := services.DefaultAICache.Get(userID, req.Query); ok {
			writeJSON(w, http.StatusOK, askResponse{Answer: answer})
			return
		}

		// 2. Use Singleflight to coalesce concurrent requests
		// Keys for singleflight should be user-specific to ensure data privacy.
		sfKey := fmt.Sprintf("%d:%s", userID, req.Query)
		result, err, _ := aiGroup.Do(sfKey, func() (interface{}, error) {
			dbType := "sqlite"
			if cfg != nil {
				dbType = cfg.DBType
			}

			// Step 1: Use Gemini to convert natural language → SQL
			sqlQuery, err := services.GenerateSQL(req.Query, userID, dbType)
			if err != nil {
				return "", fmt.Errorf("LLM Error: %w", err)
			}

			// Step 2: Execute the generated SQL safely
			results, err := store.ExecuteReadQuery(sqlQuery)
			if err != nil {
				return "", fmt.Errorf("DB Error: %w", err)
			}

			// Step 3: Use Gemini to summarize the raw data
			summary, err := services.SummarizeData(req.Query, results)
			if err != nil {
				return "", fmt.Errorf("Summarize Error: %w", err)
			}

			// Cache the final summary
			services.DefaultAICache.Set(userID, req.Query, summary)

			return summary, nil
		})

		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, askResponse{Answer: result.(string)})
	}
}
