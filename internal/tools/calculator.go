package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
)

func init() {
	// Register the split_expense tool
	_ = Register(Tool{
		Name:        "split_expense",
		Description: "Calculate how a total expense amount is split among multiple participants.",
		InputSchema: Schema{
			Type: "object",
			Properties: map[string]Property{
				"amount": {
					Type:        "number",
					Description: "Total expense amount to be split (e.g. 100.50)",
				},
				"payer": {
					Type:        "string",
					Description: "The name or ID of the user who paid for the expense",
				},
				"participants": {
					Type:        "array",
					Description: "Comma-separated list of participant names/IDs (e.g., 'Alice,Bob,Charlie')",
				},
			},
			Required: []string{"amount", "payer", "participants"},
		},
		Handler: handleSplitExpense,
	})

	// Register the calculate_balances tool
	_ = Register(Tool{
		Name:        "calculate_balances",
		Description: "Calculate net balances and suggest transactions to settle up a list of expenses.",
		InputSchema: Schema{
			Type: "object",
			Properties: map[string]Property{
				"expenses_json": {
					Type:        "string",
					Description: "JSON array of expenses. Each expense must have: 'payer' (string), 'amount' (number), and 'participants' (array of strings). Example: '[{\"payer\":\"Alice\",\"amount\":30,\"participants\":[\"Alice\",\"Bob\",\"Charlie\"]}]'",
				},
			},
			Required: []string{"expenses_json"},
		},
		Handler: handleCalculateBalances,
	})
}

// SplitExpenseInput defines input structure for split_expense tool.
type SplitExpenseInput struct {
	Amount       float64 `json:"amount"`
	Payer        string  `json:"payer"`
	Participants string  `json:"participants"` // comma-separated
}

// SplitExpenseResult defines the response structure for split_expense tool.
type SplitExpenseResult struct {
	Payer          string             `json:"payer"`
	Total          float64            `json:"total"`
	SharePerPerson float64            `json:"sharePerPerson"`
	OwedShares     map[string]float64 `json:"owedShares"`
}

func handleSplitExpense(ctx context.Context, args json.RawMessage) (any, error) {
	var input SplitExpenseInput
	if err := json.Unmarshal(args, &input); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if input.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	if input.Payer == "" {
		return nil, errors.New("payer is required")
	}

	rawParts := strings.Split(input.Participants, ",")
	var parts []string
	for _, p := range rawParts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}

	if len(parts) == 0 {
		return nil, errors.New("at least one participant is required")
	}

	// Calculate share
	numParticipants := float64(len(parts))
	share := math.Round((input.Amount/numParticipants)*100) / 100

	owedShares := make(map[string]float64)
	for _, p := range parts {
		if p != input.Payer {
			owedShares[p] = share
		}
	}

	return SplitExpenseResult{
		Payer:          input.Payer,
		Total:          input.Amount,
		SharePerPerson: share,
		OwedShares:     owedShares,
	}, nil
}

// ExpenseItem represents a single expense record in the calculate_balances tool.
type ExpenseItem struct {
	Payer        string   `json:"payer"`
	Amount       float64  `json:"amount"`
	Participants []string `json:"participants"`
}

type CalculateBalancesInput struct {
	ExpensesJSON string `json:"expenses_json"`
}

// SettleUpTransaction represents a suggested transfer to balance debts.
type SettleUpTransaction struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

// BalancesResult contains the final net balances and the optimized settlement list.
type BalancesResult struct {
	Balances     map[string]float64    `json:"balances"`
	Transactions []SettleUpTransaction `json:"suggestedTransactions"`
}

func handleCalculateBalances(ctx context.Context, args json.RawMessage) (any, error) {
	var input CalculateBalancesInput
	if err := json.Unmarshal(args, &input); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	var expenses []ExpenseItem
	if err := json.Unmarshal([]byte(input.ExpensesJSON), &expenses); err != nil {
		return nil, fmt.Errorf("invalid expenses_json formatting: %w", err)
	}

	// Calculate net balances for each person
	balances := make(map[string]float64)
	for _, exp := range expenses {
		if exp.Amount <= 0 {
			continue
		}
		if len(exp.Participants) == 0 {
			continue
		}

		share := math.Round((exp.Amount/float64(len(exp.Participants)))*100) / 100

		// Payer gets credited full amount
		balances[exp.Payer] += exp.Amount

		// Each participant gets debited their share
		for _, p := range exp.Participants {
			balances[p] -= share
		}
	}

	// Clean up balances to avoid float precision issues around zero
	for person, val := range balances {
		balances[person] = math.Round(val*100) / 100
		if math.Abs(balances[person]) < 0.01 {
			delete(balances, person)
		}
	}

	// Compute settlement transactions (greedy algorithm)
	type personBalance struct {
		name    string
		balance float64
	}

	var debtors []personBalance   // balance < 0
	var creditors []personBalance // balance > 0

	for name, val := range balances {
		if val < 0 {
			debtors = append(debtors, personBalance{name: name, balance: val})
		} else if val > 0 {
			creditors = append(creditors, personBalance{name: name, balance: val})
		}
	}

	var txs []SettleUpTransaction

	// greedily match debtors and creditors
	dIdx, cIdx := 0, 0
	for dIdx < len(debtors) && cIdx < len(creditors) {
		debtor := &debtors[dIdx]
		creditor := &creditors[cIdx]

		debtAmt := -debtor.balance
		credAmt := creditor.balance
		settleAmt := math.Min(debtAmt, credAmt)

		txs = append(txs, SettleUpTransaction{
			From:   debtor.name,
			To:     creditor.name,
			Amount: math.Round(settleAmt*100) / 100,
		})

		debtor.balance += settleAmt
		creditor.balance -= settleAmt

		if math.Abs(debtor.balance) < 0.01 {
			dIdx++
		}
		if math.Abs(creditor.balance) < 0.01 {
			cIdx++
		}
	}

	return BalancesResult{
		Balances:     balances,
		Transactions: txs,
	}, nil
}
