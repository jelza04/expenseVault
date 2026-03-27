package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"expenseVault/models"
)

// APIClient handles HTTP communication between the TUI and the backend API server.
// This enables the UI to interact with backend endpoints instead of direct DB access.
type APIClient struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewAPIClient creates a new API client pointing to the given server URL.
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ── Auth Methods ────────────────────────────────────────────

// Signup creates a new account via POST /api/signup.
func (c *APIClient) Signup(username, password string) error {
	body := map[string]string{"username": username, "password": password}
	resp, err := c.post("/api/signup", body, false)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return c.parseError(resp)
	}
	return nil
}

// Login authenticates via POST /api/login and stores the JWT token.
func (c *APIClient) Login(username, password string) (*models.User, error) {
	body := map[string]string{"username": username, "password": password}
	resp, err := c.post("/api/login", body, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var result struct {
		Token    string `json:"token"`
		Username string `json:"username"`
		UserID   int64  `json:"user_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode login response: %w", err)
	}

	c.Token = result.Token

	return &models.User{
		ID:       result.UserID,
		Username: result.Username,
	}, nil
}

// ── Transaction Methods ─────────────────────────────────────

// GetTransactions fetches all transactions via GET /api/transactions.
func (c *APIClient) GetTransactions() ([]models.Transaction, error) {
	resp, err := c.get("/api/transactions")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var result struct {
		Transactions []models.Transaction `json:"transactions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Transactions, nil
}

// AddTransaction adds a new transaction via POST /api/transactions.
func (c *APIClient) AddTransaction(txType string, amount float64, category, description, date, notes string) (int64, error) {
	body := map[string]any{
		"type":        txType,
		"amount":      amount,
		"category":    category,
		"description": description,
		"date":        date,
		"notes":       notes,
	}
	resp, err := c.post("/api/transactions", body, true)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return 0, c.parseError(resp)
	}

	var result struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	return result.ID, nil
}

// ── AI Methods ──────────────────────────────────────────────

// Ask sends a natural language query via POST /api/ask.
func (c *APIClient) Ask(query string) (string, error) {
	body := map[string]string{"query": query}
	resp, err := c.post("/api/ask", body, true)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", c.parseError(resp)
	}

	var result struct {
		Answer string `json:"answer"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Answer, nil
}

// ── Internal Helpers ────────────────────────────────────────

func (c *APIClient) get(path string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	return c.HTTPClient.Do(req)
}

func (c *APIClient) post(path string, body any, auth bool) (*http.Response, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.BaseURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if auth && c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	return c.HTTPClient.Do(req)
}

func (c *APIClient) parseError(resp *http.Response) error {
	bodyBytes, _ := io.ReadAll(resp.Body)
	var errResp struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(bodyBytes, &errResp) == nil && errResp.Error != "" {
		return fmt.Errorf("%s", errResp.Error)
	}
	return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
}
