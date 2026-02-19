package services

import (
	"strings"

	"expenseVault/models"
)

// ============================================================
// UNIT 3: Callbacks, Closures, Function Expressions
// UNIT 1: Control flow — switch, if-else
// Rule-based categorization with keyword matching
// ============================================================

// CategorizationRule defines a rule for auto-categorizing transactions.
type CategorizationRule struct {
	Keywords []string
	Category models.Category
}

// DefaultRules returns the default categorization rules.
// Demonstrates: Composite literals, slice of structs (Unit 2)
func DefaultRules() []CategorizationRule {
	return []CategorizationRule{
		{Keywords: []string{"food", "lunch", "dinner", "breakfast", "restaurant", "cafe", "canteen", "zomato", "swiggy"}, Category: models.CategoryFood},
		{Keywords: []string{"uber", "ola", "bus", "train", "metro", "fuel", "petrol", "flight", "taxi"}, Category: models.CategoryTravel},
		{Keywords: []string{"amazon", "flipkart", "shop", "mall", "store", "purchase", "buy"}, Category: models.CategoryShopping},
		{Keywords: []string{"electricity", "water", "rent", "wifi", "internet", "phone", "recharge", "gas"}, Category: models.CategoryBills},
		{Keywords: []string{"hospital", "doctor", "medicine", "pharmacy", "medical", "gym", "health"}, Category: models.CategoryHealth},
		{Keywords: []string{"book", "course", "tuition", "school", "college", "udemy", "tutorial"}, Category: models.CategoryEducation},
		{Keywords: []string{"movie", "netflix", "spotify", "game", "concert", "party", "fun"}, Category: models.CategoryEntertainment},
		{Keywords: []string{"salary", "stipend", "pay", "wages"}, Category: models.CategorySalary},
		{Keywords: []string{"freelance", "gig", "project", "contract", "consulting"}, Category: models.CategoryFreelance},
	}
}

// Categorizer provides rule-based and callback-based categorization.
type Categorizer struct {
	rules []CategorizationRule
}

// NewCategorizer creates a categorizer with default rules.
func NewCategorizer() *Categorizer {
	return &Categorizer{rules: DefaultRules()}
}

// AddRule adds a custom categorization rule.
func (c *Categorizer) AddRule(keywords []string, category models.Category) {
	c.rules = append(c.rules, CategorizationRule{
		Keywords: keywords,
		Category: category,
	})
}

// AutoCategorize categorizes a description based on keyword rules.
// Demonstrates: for range, strings operations, control flow (Unit 1 + Unit 2)
func (c *Categorizer) AutoCategorize(description string) models.Category {
	lower := strings.ToLower(description)

	for _, rule := range c.rules {
		for _, keyword := range rule.Keywords {
			if strings.Contains(lower, keyword) {
				return rule.Category
			}
		}
	}

	return models.CategoryOther
}

// ============================================================
// UNIT 3: Callback pattern for categorization
// ============================================================

// CategorizationCallback is a function that categorizes a description.
// Demonstrates: Function type for callback (Unit 3)
type CategorizationCallback func(description string) models.Category

// CategorizeWithCallback applies a custom categorization function.
// Demonstrates: Callback pattern (Unit 3)
func CategorizeWithCallback(transactions []models.Transaction, cb CategorizationCallback) []models.Transaction {
	result := make([]models.Transaction, len(transactions))
	copy(result, transactions)

	for i := range result {
		if result[i].Category == "" || result[i].Category == models.CategoryOther {
			result[i].Category = cb(result[i].Description)
		}
	}
	return result
}

// ============================================================
// UNIT 3: Returning a Function / Closure
// ============================================================

// MakeKeywordMatcher returns a function that checks if a string
// contains any of the given keywords.
// Demonstrates: Returning a function, closure over 'keywords' (Unit 3)
func MakeKeywordMatcher(keywords ...string) func(string) bool {
	// The returned function closes over 'keywords'
	return func(input string) bool {
		lower := strings.ToLower(input)
		for _, kw := range keywords {
			if strings.Contains(lower, strings.ToLower(kw)) {
				return true
			}
		}
		return false
	}
}

// MakeCategorizer returns a categorization closure for a specific category.
// Demonstrates: Returning a function (Unit 3)
func MakeCategorizer(category models.Category, keywords ...string) CategorizationCallback {
	matcher := MakeKeywordMatcher(keywords...)
	return func(description string) models.Category {
		if matcher(description) {
			return category
		}
		return models.CategoryOther
	}
}
