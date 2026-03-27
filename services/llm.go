package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Gemini API endpoint for flash model
// Gemini API endpoint for flash model
const geminiBaseURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key="

// GeminiRequest represents the payload structure for the Gemini API
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

// GeminiResponse represents the expected response from the Gemini API
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

const dbSchemaContext = `
The database schema consists of two tables. ALL currency values are in Indian Rupees (₹).

CREATE TABLE users (
	id INTEGER PRIMARY KEY,
	username TEXT,
	created_at DATETIME
);

CREATE TABLE transactions (
	id INTEGER PRIMARY KEY,
	user_id INTEGER REFERENCES users(id),
	type TEXT NOT NULL,          -- 'INCOME' or 'EXPENSE'
	amount REAL NOT NULL,        -- Amount in Indian Rupees (₹). 
	category TEXT NOT NULL,      -- Example: 'Food', 'Transport', 'Entertainment'
	description TEXT NOT NULL,
	date TEXT NOT NULL,          -- Format: 'YYYY-MM-DD'
	notes TEXT,
	created_at DATETIME
);

IMPORTANT RULES:
1. ONLY return a single standard SQL SELECT statement block starting with "SELECT " and ending with ";".
2. DO NOT wrap the SQL in markdown blocks (e.g., no `+"```"+`sql or `+"```"+`).
3. DO NOT include any explanations or extra text.
4. Filter by user_id = <USER_ID> to ensure users only query their own data.
5. If the user asks for "month-wise" or "monthly" grouping:
   - For SQLite: Use strftime('%Y-%m', date) as month
   - For MySQL: Use DATE_FORMAT(date, '%Y-%m') as month
6. Dates are stored as strings in 'YYYY-MM-DD' format.
7. Treat all currency as Rupees (₹).
`

// GenerateSQL uses the Gemini API to translate a natural language query into a SQL statement.
func GenerateSQL(query string, userID int64, dbFlavor string) (string, error) {
	apiKey := strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable is not set. Please add it to your .env file")
	}

	prompt := fmt.Sprintf("%s\n\nDATABASE TYPE: %s\n\nUSER QUERY: %s\n\nNote: The user's user_id is %d. Replace <USER_ID> with %d in your query.", dbSchemaContext, dbFlavor, query, userID, userID)

	return callGeminiAPI(prompt, apiKey)
}

// SummarizeData takes the raw JSON result from the database query and formats it into a conversational summary.
func SummarizeData(query string, data []map[string]interface{}) (string, error) {
	apiKey := strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable is not set. Please add it to your .env file")
	}

	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal data for LLM: %w", err)
	}

	prompt := fmt.Sprintf(`
You are an AI assistant for a CLI expense tracker app. You are helping an Indian user manage their finances in Indian Rupees (₹).

USER'S ORIGINAL QUESTION: "%s"

DATABASE RESULTS (JSON format):
%s

INSTRUCTIONS:
1. Based *only* on the database results provided, answer the user's question.
2. Provide a concise, friendly, conversational response.
3. If the results are empty, tell the user no data was found that matched their query.
4. MANDATORY: You MUST use the Rupee symbol (₹) for all currency amounts. 
   - CORRECT: "Your balance is ₹5,000.00"
   - INCORRECT: "Your balance is $5,000.00"
   - NEVER use the dollar sign ($) under any circumstances.
5. You can use slight terminal markdown (like **bolding** key numbers). Do not be overly verbose.
`, query, string(dataJSON))

	return callGeminiAPI(prompt, apiKey)
}

// callGeminiAPI is a helper to perform the HTTP request to Google's Gemini API.
func callGeminiAPI(prompt, apiKey string) (string, error) {
	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := geminiBaseURL + apiKey
	
	var lastErr error
	maxRetries := 3
	backoff := 2 * time.Second

	for i := 0; i <= maxRetries; i++ {
		if i > 0 {
			time.Sleep(backoff)
			backoff *= 2
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			lastErr = fmt.Errorf("failed to connect to Gemini API: %w", err)
			continue
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read API response body: %w", err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(bodyBytes))
			if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
				continue // Retry on rate limit or server errors
			}
			return "", lastErr // Return immediately on auth or client errors
		}

		var geminiResp GeminiResponse
		if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
			return "", fmt.Errorf("failed to decode API response: %w", err)
		}

		if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
			return "", fmt.Errorf("unexpected empty response from Gemini API")
		}

		// Clean up any potential markdown code blocks if the LLM hallucinated them
		response := geminiResp.Candidates[0].Content.Parts[0].Text
		response = strings.TrimPrefix(response, "```sql\n")
		response = strings.TrimPrefix(response, "```\n")
		response = strings.TrimSuffix(response, "\n```")
		response = strings.TrimSpace(response)

		return response, nil
	}

	return "", fmt.Errorf("after %d retries: %w", maxRetries, lastErr)
}
