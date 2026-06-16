package main

import (
	"encoding/json"
	"fmt"

	"finduo-ai/internal/tools"
)

func main() {
	fmt.Println("Finduo AI Tools Structure Demo")
	fmt.Println("=============================")

	availableTools := tools.List()
	fmt.Printf("Registered %d tools:\n", len(availableTools))

	for _, t := range availableTools {
		schemaJSON, _ := json.MarshalIndent(t.InputSchema, "  ", "  ")
		fmt.Printf("- %s: %s\n  Schema:\n  %s\n\n", t.Name, t.Description, string(schemaJSON))
	}
}
