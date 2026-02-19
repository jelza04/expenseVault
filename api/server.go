package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"expenseVault/db"
	"expenseVault/models"

	"golang.org/x/crypto/bcrypt"
)

// ============================================================
// REST API Server using net/http (standard library)
// Demonstrates: Goroutines (concurrent server), methods,
//               error handling, JSON serialization
// ============================================================

// Server holds the API server state.
type Server struct {
	store  *db.Store
	jwtKey []byte
	mux    *http.ServeMux
}

// NewServer creates a new API server.
func NewServer(store *db.Store, jwtSecret string) *Server {
	s := &Server{
		store:  store,
		jwtKey: []byte(jwtSecret),
		mux:    http.NewServeMux(),
	}
	s.registerRoutes()
	return s
}

// registerRoutes sets up all HTTP routes.
func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/api/health", s.healthHandler)
	s.mux.HandleFunc("/api/register", s.registerHandler)
	s.mux.HandleFunc("/api/login", s.loginHandler)
	s.mux.HandleFunc("/api/transactions", s.authMiddleware(s.transactionsHandler))
	s.mux.HandleFunc("/api/transactions/", s.authMiddleware(s.transactionByIDHandler))
	s.mux.HandleFunc("/api/sync", s.authMiddleware(s.syncHandler))
}

// Start runs the HTTP server.
// Demonstrates: Goroutines for concurrent server (Go concurrency)
func (s *Server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("🚀 ExpenseVault API Server starting on http://localhost%s\n", addr)
	fmt.Println("   Endpoints:")
	fmt.Println("   POST   /api/register       — Create account")
	fmt.Println("   POST   /api/login           — Get JWT token")
	fmt.Println("   GET    /api/transactions     — List transactions")
	fmt.Println("   POST   /api/transactions     — Add transaction")
	fmt.Println("   GET    /api/transactions/:id — Get transaction")
	fmt.Println("   DELETE /api/transactions/:id — Delete transaction")
	fmt.Println("   POST   /api/sync             — Sync transactions")
	fmt.Println("   GET    /api/health           — Health check")
	fmt.Println()
	return http.ListenAndServe(addr, s.mux)
}

// ============================================================
// Handlers
// ============================================================

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "ExpenseVault API",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// registerHandler creates a new user account.
// Demonstrates: Password hashing with bcrypt, error handling
func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		jsonError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	// Demonstrates: Password hashing with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	id, err := s.store.CreateUser(req.Username, string(hash))
	if err != nil {
		jsonError(w, http.StatusConflict, "Username already exists")
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"id":       id,
		"username": req.Username,
		"message":  "User created successfully",
	})
}

// loginHandler authenticates and returns a JWT token.
func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := s.store.GetUserByUsername(req.Username)
	if err != nil {
		jsonError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Demonstrates: bcrypt password comparison
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		jsonError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := GenerateToken(user.Username, s.jwtKey)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{
		"token":   token,
		"message": "Login successful",
	})
}

// transactionsHandler handles GET (list) and POST (create) for transactions.
func (s *Server) transactionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		transactions, err := s.store.GetAllTransactions()
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to get transactions")
			return
		}
		jsonResponse(w, http.StatusOK, transactions)

	case http.MethodPost:
		var tx models.Transaction
		if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		id, err := s.store.AddTransaction(tx)
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		jsonResponse(w, http.StatusCreated, map[string]interface{}{
			"id":      id,
			"message": "Transaction added",
		})

	default:
		jsonError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// transactionByIDHandler handles GET and DELETE for a specific transaction.
func (s *Server) transactionByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path: /api/transactions/123
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		jsonError(w, http.StatusBadRequest, "Missing transaction ID")
		return
	}

	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		jsonError(w, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		tx, err := s.store.GetTransaction(id)
		if err != nil {
			jsonError(w, http.StatusNotFound, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, tx)

	case http.MethodDelete:
		if err := s.store.DeleteTransaction(id); err != nil {
			jsonError(w, http.StatusNotFound, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, map[string]string{"message": "Transaction deleted"})

	default:
		jsonError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// syncHandler handles transaction syncing from client.
func (s *Server) syncHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var payload models.SyncPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		jsonError(w, http.StatusBadRequest, "Invalid sync payload")
		return
	}

	count, err := s.store.BulkInsert(payload.Transactions)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Sync failed")
		return
	}

	// Return server's transactions for the client
	serverTxs, _ := s.store.GetAllTransactions()

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"synced":       count,
		"message":      "Sync complete",
		"transactions": serverTxs,
	})
}

// ============================================================
// Auth Middleware using JWT
// ============================================================

func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			jsonError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			jsonError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		_, err := ValidateToken(parts[1], s.jwtKey)
		if err != nil {
			jsonError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		next(w, r)
	}
}

// ============================================================
// JSON helpers
// ============================================================

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, status int, message string) {
	jsonResponse(w, status, map[string]string{"error": message})
}
