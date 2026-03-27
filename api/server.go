package api

import (
	"log"
	"net/http"
	"time"

	"expenseVault/db"
	"expenseVault/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// StartServer starts the HTTP server with Chi router and middleware pipeline.
// UNIT 5: Demonstrates efficient routing and request management.
func StartServer(addr string, store *db.Store, cfg *utils.Config) error {
	// Initialize the shared JWT secret from config
	JwtSecret = []byte(cfg.JWTSecret)

	r := chi.NewRouter()

	// ── Middleware Pipeline ──────────────────────────────────
	// Request management: logging, panic recovery, timeouts.
	r.Use(middleware.Logger)      // Logs every request (method, path, duration)
	r.Use(middleware.Recoverer)   // Recovers from panics and returns 500
	r.Use(middleware.Timeout(30 * time.Second)) // Request timeout
	r.Use(CORSMiddleware)         // Allow cross-origin requests

	// ── Public Routes (no auth required) ────────────────────
	r.Get("/health", handleHealth)

	r.Post("/api/signup", makeSignupHandler(store))
	r.Post("/api/login", makeLoginHandler(store))

	// ── Protected Routes (JWT required) ─────────────────────
	r.Group(func(r chi.Router) {
		r.Use(JWTAuthMiddleware(store))

		// Transaction CRUD
		r.Get("/api/transactions", makeListTransactionsHandler(store))
		r.Post("/api/transactions", makeAddTransactionHandler(store))

		// AI Ask
		r.Post("/api/ask", makeAskHandler(store, cfg))

		// Sync (existing)
		r.Post("/api/sync", handleSync)
	})

	log.Printf("Server listening on %s", addr)
	return http.ListenAndServe(addr, r)
}
