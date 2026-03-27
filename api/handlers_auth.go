package api

import (
	"encoding/json"
	"net/http"

	"expenseVault/db"

	"golang.org/x/crypto/bcrypt"
)

// ── Request / Response types ────────────────────────────────

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
}

// ── Handlers ────────────────────────────────────────────────

// makeSignupHandler returns a handler for POST /api/signup.
// Demonstrates bcrypt password hashing over HTTP.
func makeSignupHandler(store *db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req authRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		if req.Username == "" || req.Password == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username and password are required"})
			return
		}

		// Bcrypt hashing — secure credential storage
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
			return
		}

		userID, err := store.CreateUser(req.Username, string(hash))
		if err != nil {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "username already exists"})
			return
		}

		writeJSON(w, http.StatusCreated, map[string]any{
			"message": "signup successful",
			"user_id": userID,
		})
	}
}

// makeLoginHandler returns a handler for POST /api/login.
// Demonstrates bcrypt verification and JWT token generation over HTTP.
func makeLoginHandler(store *db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req authRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		if req.Username == "" || req.Password == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username and password are required"})
			return
		}

		user, err := store.GetUserByUsername(req.Username)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid username or password"})
			return
		}

		// Bcrypt comparison — verifies hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid username or password"})
			return
		}

		_ = store.UpdateLastLogin(user.ID)

		// Generate JWT for the authenticated session
		token, err := GenerateToken(req.Username, JwtSecret)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
			return
		}

		writeJSON(w, http.StatusOK, loginResponse{
			Token:    token,
			Username: user.Username,
			UserID:   user.ID,
		})
	}
}
