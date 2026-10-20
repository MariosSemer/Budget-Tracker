package budget

type Category string

const (
	//expense
	Groceries     Category = "groceries"
	Bills         Category = "bills"
	Cat           Category = "cat"
	Cigarettes    Category = "cigarettes"
	Entertainment Category = "entertainment"
	OtherExpense  Category = "other_expense"
	//income
	Salary         Category = "salary"
	FamilyReceived Category = "family_received"
	OtherIncome    Category = "other_income"
)

type Kind string

const (
	Income  Kind = "income"
	Expense Kind = "expense"
)

// The single source of truth: every valid category and its direction.
var categoryKind = map[Category]Kind{
	Groceries:      Expense,
	Bills:          Expense,
	Cat:            Expense,
	Cigarettes:     Expense,
	Entertainment:  Expense,
	OtherExpense:   Expense,
	Salary:         Income,
	FamilyReceived: Income,
	OtherIncome:    Income,
}

// Kind returns the category's direction. ok is false for an unknown category.
func (c Category) Kind() (kind Kind, ok bool) {
	kind, ok = categoryKind[c]
	return kind, ok
}

type Transaction struct {
	ID          int64    `json:"id"`                    // assigned by the store, never by the caller
	Category    Category `json:"category"`              // decides whether it is income or expense
	Amount      int64    `json:"amount"`                // in cents, always positive
	Date        string   `json:"date"`                  // the day it happened
	Description string   `json:"description,omitempty"` // optional
}
