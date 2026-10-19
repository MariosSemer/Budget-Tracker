# Budget Tracker

A small single-user HTTP API, written in Go, for recording every income and expense as it happens, so I can see my balance, where my money goes, and how much I save each month.

## Overview

| | |
|---|---|
| Users | One (no accounts or login in the first version; a CLI and multi-user support may come later) |
| Currency | One (euro) |
| Storage | In memory first, behind an interface; SQL (or a CSV file) can be added later |
| Language | Go 1.22 or newer (needed for method patterns in `net/http` routes) |

### Conventions used everywhere

- **Money is an integer number of cents.** `350` means 3.50. Amounts are always positive; the category says whether money came in or went out.
- **Dates are `YYYY-MM-DD` strings**, interpreted in the `Europe/Athens` time zone.
- **Errors have one shape:** `{"error": "message"}`.
- **Wrong method on a known route** returns `405` automatically from the router.

## Data model

### Categories

Every category goes in exactly one direction, so the category alone decides whether a transaction is income or expense. There is no separate type field.

| Category | Direction |
|---|---|
| `groceries` | expense |
| `bills` | expense |
| `entertainment` | expense |
| `cat` | expense |
| `cigarettes` | expense |
| `other_expense` | expense |
| `salary` | income |
| `family_received` | income |
| `other_income` | income |

The list lives in one place in the code (a map from category to direction). Adding a category is one constant plus one map entry.

### Transaction

| Field | Type | Notes |
|---|---|---|
| `id` | int64 | Assigned by the store, never by the caller |
| `category` | Category | Must be one of the categories above |
| `amount` | int64 | Cents, always positive |
| `date` | date | The day it happened |
| `description` | string | Optional (Beer, Coffee, Supermarket) |

Responses also include a derived `kind` (`income` or `expense`) on every transaction.

## API

| Operation | Endpoint | Success |
|---|---|---|
| Add a transaction | `POST /transactions` | `201 Created` |
| Delete a transaction | `DELETE /transactions/{id}` | `204 No Content` |
| List transactions (with balance) | `GET /transactions` | `200 OK` |
| Balance | `GET /balance` | `200 OK` |
| Summary for a month | `GET /summary?month=YYYY-MM` | `200 OK` |

### Add a transaction: `POST /transactions`

Request body (the ID is never sent):

```json
{"category": "groceries", "amount": 2550, "date": "2026-10-11", "description": "Supermarket"}
```

Validation:

- The body must be valid JSON of the right types, at most **1 KB**.
- `category` is required and must be a known category.
- `amount` is required: a positive integer, in cents.
- `date` is required, in `YYYY-MM-DD` format, and **not in the future**.
- `description` is optional, trimmed, at most **100 characters** (counted as characters, not bytes, so Greek text works).

Response `201 Created`, with a `Location: /transactions/{id}` header and the created transaction:

```json
{"id": 3, "kind": "expense", "category": "groceries", "amount": 2550, "date": "2026-10-11", "description": "Supermarket"}
```

| Case | Status |
|---|---|
| Invalid JSON, wrong type, or a failed validation rule | 400 |
| Body larger than 1 KB | 413 |
| Store failed | 500 |
| Success | 201 |

There is deliberately **no rule tying an expense to the current balance**: the app records what happened, and the balance is allowed to go negative.

### Delete a transaction: `DELETE /transactions/{id}`

- The ID comes from the URL. There is no request body.
- Validation: the ID must be a positive integer that fits in `int64`.

| Case | Status |
|---|---|
| Invalid ID | 400 |
| No transaction with that ID (including deleting the same ID twice) | 404 |
| Store failed | 500 |
| Success, empty body | 204 |

### List transactions: `GET /transactions`

- The caller sends nothing (filters may come later).
- All transactions are fetched in one call, sorted **newest date first** (ties by ID, descending). The balance is computed from that same list, so the two always agree.

```json
{
  "balance": 97450,
  "transactions": [
    {"id": 3, "kind": "expense", "category": "groceries", "amount": 2550, "date": "2026-10-11", "description": "Supermarket"},
    {"id": 1, "kind": "income", "category": "salary", "amount": 100000, "date": "2026-10-01"}
  ]
}
```

| Case | Status |
|---|---|
| Store failed | 500 |
| Success, including an empty list (`{"balance": 0, "transactions": []}`) | 200 |

### Balance: `GET /balance`

- The caller sends nothing.
- Calculation: start at **0**; for every transaction, add its amount if the category is income, subtract it if it is an expense.
- The balance can be negative. The money you had on day one is entered as an income transaction.

```json
{"balance": 97450}
```

| Case | Status |
|---|---|
| Store failed | 500 |
| Success, including no transactions (`{"balance": 0}`) | 200 |

### Summary for a month: `GET /summary?month=YYYY-MM`

- The caller sends `month` as a query parameter.
- Validation: `month` is required, formatted `YYYY-MM`, and **not after the current month**.
- Calculation:
  1. Keep only the transactions dated in that month.
  2. Add each amount to its category's total.
  3. Split the totals by direction into `income` and `expenses`.
  4. `net = income - expenses`.
- Rows are sorted income first, then expenses, each by total descending. Categories with no transactions that month are left out.

```json
{
  "month": "2026-10",
  "income": 105000,
  "expenses": 12000,
  "net": 93000,
  "categories": [
    {"category": "salary", "kind": "income", "total": 100000},
    {"category": "family_received", "kind": "income", "total": 5000},
    {"category": "bills", "kind": "expense", "total": 8000},
    {"category": "groceries", "kind": "expense", "total": 4000}
  ]
}
```

`net` is that month's change, **not the balance**. The balance counts everything ever recorded.

| Case | Status |
|---|---|
| `month` missing, malformed, or in the future | 400 |
| Store failed | 500 |
| Success, including an empty month (zeros and `"categories": []`) | 200 |