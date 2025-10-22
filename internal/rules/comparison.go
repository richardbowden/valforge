package rules

import (
	"github.com/richardbowden/valforge/internal/builder"
	"github.com/richardbowden/valforge/internal/vtypes"
)

type GreaterThanRule struct{}

func (r GreaterThanRule) Name() string              { return "gt" }
func (r GreaterThanRule) Priority() int             { return 3 }
func (r GreaterThanRule) RequiredImports() []string { return nil }
func (r GreaterThanRule) Aliases() []string         { return []string{"gte"} }

func (r GreaterThanRule) SupportsType(fieldType vtypes.FieldType) bool {
	return IntegerTypes.Contains(fieldType.Kind)
}

func (r GreaterThanRule) Generate(cb *builder.CodeBuilder, field vtypes.ValidationField, structName string) error {

	if gtVal, exists := field.Rules["gt"]; exists {
		cb.Printf(`if v.%s <= %s {`, field.Name, gtVal)
		cb.Indent()
		cb.Printf(`verr.AddFieldError("%s", "%s must be greater than %s", v.%s)`,
			field.JSONName, field.JSONName, gtVal, field.Name)
		cb.Dedent()
		cb.Writeln("}")
	}

	if gtVal, exists := field.Rules["gte"]; exists {
		cb.Printf(`if v.%s < %s {`, field.Name, gtVal)
		cb.Indent()
		cb.Printf(`verr.AddFieldError("%s", "%s must be greater than or equal to %s", v.%s)`,
			field.JSONName, field.JSONName, gtVal, field.Name)
		cb.Dedent()
		cb.Writeln("}")
	}

	return nil
}

type LessThanRule struct{}

func (r LessThanRule) Name() string              { return "lt" }
func (r LessThanRule) Priority() int             { return 3 }
func (r LessThanRule) RequiredImports() []string { return nil }
func (r LessThanRule) Aliases() []string         { return []string{"lte"} }

func (r LessThanRule) SupportsType(fieldType vtypes.FieldType) bool {
	return IntegerTypes.Contains(fieldType.Kind)
}

func (r LessThanRule) Generate(cb *builder.CodeBuilder, field vtypes.ValidationField, structName string) error {
	if ltVal, exists := field.Rules["lt"]; exists {
		cb.Printf(`if v.%s >= %s {`, field.Name, ltVal)
		cb.Indent()
		cb.Printf(`verr.AddFieldError("%s", "%s must be less than %s", v.%s)`,
			field.JSONName, field.JSONName, ltVal, field.Name)
		cb.Dedent()
		cb.Writeln("}")
	}

	if gtVal, exists := field.Rules["lte"]; exists {
		cb.Printf(`if v.%s > %s {`, field.Name, gtVal)
		cb.Indent()
		cb.Printf(`verr.AddFieldError("%s", "%s must be less than or equal to %s", v.%s)`,
			field.JSONName, field.JSONName, gtVal, field.Name)
		cb.Dedent()
		cb.Writeln("}")
	}

	return nil
}

type EqualFieldRule struct{}

func (r EqualFieldRule) Name() string              { return "eqfield" }
func (r EqualFieldRule) Priority() int             { return 5 }
func (r EqualFieldRule) RequiredImports() []string { return nil }
func (r EqualFieldRule) Aliases() []string         { return []string{} }

func (r EqualFieldRule) SupportsType(fieldType vtypes.FieldType) bool {
	return AllTypes.Contains(fieldType.Kind)
}

func (r EqualFieldRule) Generate(cb *builder.CodeBuilder, field vtypes.ValidationField, structName string) error {
	if targetField, exists := field.Rules["eqfield"]; exists {
		cb.Printf(`if v.%s != v.%s {`, field.Name, targetField)
		cb.Indent()
		cb.Printf(`verr.AddFieldError("%s", "%s must match %s", v.%s)`,
			field.JSONName, field.JSONName, targetField, field.Name)
		cb.Dedent()
		cb.Writeln("}")
	}
	return nil
}
