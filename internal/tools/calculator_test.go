package tools

import (
	"context"
	"encoding/json"
	"testing"
)

func TestSplitExpense(t *testing.T) {
	ctx := context.Background()

	input := SplitExpenseInput{
		Amount:       150.00,
		Payer:        "Alice",
		Participants: "Alice, Bob, Charlie",
	}

	args, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal input: %v", err)
	}

	result, err := Execute(ctx, "split_expense", args)
	if err != nil {
		t.Fatalf("failed to execute split_expense: %v", err)
	}

	res, ok := result.(SplitExpenseResult)
	if !ok {
		t.Fatalf("unexpected result type: %T", result)
	}

	if res.Payer != "Alice" {
		t.Errorf("expected payer Alice, got %s", res.Payer)
	}
	if res.Total != 150.00 {
		t.Errorf("expected total 150.00, got %f", res.Total)
	}
	if res.SharePerPerson != 50.00 {
		t.Errorf("expected share per person 50.00, got %f", res.SharePerPerson)
	}

	if val, ok := res.OwedShares["Bob"]; !ok || val != 50.00 {
		t.Errorf("expected Bob to owe 50.00, got %f", val)
	}
	if val, ok := res.OwedShares["Charlie"]; !ok || val != 50.00 {
		t.Errorf("expected Charlie to owe 50.00, got %f", val)
	}
	if _, ok := res.OwedShares["Alice"]; ok {
		t.Error("expected Alice to not be in OwedShares")
	}
}

func TestCalculateBalances(t *testing.T) {
	ctx := context.Background()

	input := CalculateBalancesInput{
		ExpensesJSON: `[
			{"payer": "Alice", "amount": 90.00, "participants": ["Alice", "Bob", "Charlie"]},
			{"payer": "Bob", "amount": 30.00, "participants": ["Bob", "Charlie"]}
		]`,
	}

	args, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal input: %v", err)
	}

	result, err := Execute(ctx, "calculate_balances", args)
	if err != nil {
		t.Fatalf("failed to execute calculate_balances: %v", err)
	}

	res, ok := result.(BalancesResult)
	if !ok {
		t.Fatalf("unexpected result type: %T", result)
	}

	if val := res.Balances["Alice"]; val != 60.00 {
		t.Errorf("expected Alice net balance to be 60.00, got %f", val)
	}
	if val := res.Balances["Bob"]; val != -15.00 {
		t.Errorf("expected Bob net balance to be -15.00, got %f", val)
	}
	if val := res.Balances["Charlie"]; val != -45.00 {
		t.Errorf("expected Charlie net balance to be -45.00, got %f", val)
	}

	if len(res.Transactions) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(res.Transactions))
	}

	foundCharlieToAlice := false
	foundBobToAlice := false

	for _, tx := range res.Transactions {
		if tx.From == "Charlie" && tx.To == "Alice" && tx.Amount == 45.00 {
			foundCharlieToAlice = true
		}
		if tx.From == "Bob" && tx.To == "Alice" && tx.Amount == 15.00 {
			foundBobToAlice = true
		}
	}

	if !foundCharlieToAlice {
		t.Error("missing suggested transaction: Charlie pays Alice 45.00")
	}
	if !foundBobToAlice {
		t.Error("missing suggested transaction: Bob pays Alice 15.00")
	}
}
