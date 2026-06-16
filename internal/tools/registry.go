package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// Tool represents a capability that can be executed by an AI agent or a system.
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema Schema      `json:"inputSchema"`
	Handler     HandlerFunc `json:"-"`
}

// Schema defines the input parameters schema (JSON Schema format).
type Schema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required,omitempty"`
}

// Property defines a single parameter property in the schema.
type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

// HandlerFunc is the function signature for executing a tool.
type HandlerFunc func(ctx context.Context, args json.RawMessage) (any, error)

// Registry manages the collection of available tools.
type Registry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

// NewRegistry creates a new tool registry.
func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

// Register adds a tool to the registry.
func (r *Registry) Register(tool Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if tool.Name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}
	if tool.Handler == nil {
		return fmt.Errorf("tool handler cannot be nil for %s", tool.Name)
	}
	if _, exists := r.tools[tool.Name]; exists {
		return fmt.Errorf("tool already registered: %s", tool.Name)
	}

	r.tools[tool.Name] = tool
	return nil
}

// Get retrieves a tool by name.
func (r *Registry) Get(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.tools[name]
	return t, ok
}

// List returns all registered tools.
func (r *Registry) List() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]Tool, 0, len(r.tools))
	for _, t := range r.tools {
		list = append(list, t)
	}
	return list
}

// Execute runs a tool by name with the given JSON arguments.
func (r *Registry) Execute(ctx context.Context, name string, args json.RawMessage) (any, error) {
	tool, ok := r.Get(name)
	if !ok {
		return nil, fmt.Errorf("tool not found: %s", name)
	}
	return tool.Handler(ctx, args)
}

// DefaultRegistry is a shared package-level registry.
var DefaultRegistry = NewRegistry()

// Register registers a tool to the DefaultRegistry.
func Register(tool Tool) error {
	return DefaultRegistry.Register(tool)
}

// Execute runs a tool in the DefaultRegistry.
func Execute(ctx context.Context, name string, args json.RawMessage) (any, error) {
	return DefaultRegistry.Execute(ctx, name, args)
}

// List returns all registered tools in the DefaultRegistry.
func List() []Tool {
	return DefaultRegistry.List()
}
