package rules

import (
	"github.com/richardbowden/valforge/internal/builder"
	"github.com/richardbowden/valforge/internal/vtypes"
)

type Registry struct {
	rules map[string]Rule
}

func NewRegistry() *Registry {
	return &Registry{
		rules: make(map[string]Rule),
	}
}

func (r *Registry) Register(rule Rule) {
	a := rule.Aliases()
	for _, alias := range a {
		r.rules[alias] = rule
	}
	r.rules[rule.Name()] = rule
}

func (r *Registry) Get(name string) (Rule, bool) {
	rule, exists := r.rules[name]
	return rule, exists
}

func (r *Registry) GetForTypeCheck(name string) (interface{ SupportsType(vtypes.FieldType) bool }, bool) {
	rule, exists := r.rules[name]
	if !exists {
		return nil, false
	}
	return rule, true
}

func (r *Registry) GetAll() map[string]Rule {
	return r.rules
}

func (r *Registry) GetAllForGeneration() map[string]interface {
	Generate(*builder.CodeBuilder, vtypes.ValidationField, string) error
	Priority() int
} {
	result := make(map[string]interface {
		Generate(*builder.CodeBuilder, vtypes.ValidationField, string) error
		Priority() int
	})

	for name, rule := range r.rules {
		result[name] = rule
	}
	return result
}

func (r *Registry) GetRequiredImports(fields []vtypes.ValidationField) []string {
	imports := make(map[string]bool)

	for _, field := range fields {
		for ruleName := range field.Rules {
			if rule, exists := r.rules[ruleName]; exists {
				for _, imp := range rule.RequiredImports() {
					imports[imp] = true
				}
			}
		}
	}

	var result []string
	for imp := range imports {
		result = append(result, imp)
	}
	return result
}

func (r *Registry) HasRule(name string) bool {
	_, exists := r.rules[name]
	return exists
}
