package rules

import (
	"github.com/richardbowden/valforge/internal/builder"
	"github.com/richardbowden/valforge/internal/vtypes"
)

type MinLenRule struct{}

func (r MinLenRule) Name() string              { return "minlen" }
func (r MinLenRule) Priority() int             { return 2 }
func (r MinLenRule) RequiredImports() []string { return nil }
func (r MinLenRule) Aliases() []string         { return []string{} }

func (r MinLenRule) SupportsType(fieldType vtypes.FieldType) bool {
	return StringTypes.Contains(fieldType.Kind)
}

func (r MinLenRule) Generate(cb *builder.CodeBuilder, field vtypes.ValidationField, structName string) error {
	if minVal, exists := field.Rules["minlen"]; exists {
		cb.Printf(`if len(v.%s) < %s {`, field.Name, minVal)
		cb.Indent()
		cb.Printf(`verr.AddFieldError("%s", "%s must be at least %s characters", v.%s)`,
			field.JSONName, field.JSONName, minVal, field.Name)
		cb.Dedent()
		cb.Writeln("}")
	}
	return nil
}

type MaxLenRule struct{}

func (r MaxLenRule) Name() string              { return "maxlen" }
func (r MaxLenRule) Priority() int             { return 2 }
func (r MaxLenRule) RequiredImports() []string { return nil }
func (r MaxLenRule) Aliases() []string         { return []string{} }

func (r MaxLenRule) SupportsType(fieldType vtypes.FieldType) bool {
	return StringTypes.Contains(fieldType.Kind)
}

func (r MaxLenRule) Generate(cb *builder.CodeBuilder, field vtypes.ValidationField, structName string) error {
	if maxVal, exists := field.Rules["maxlen"]; exists {
		cb.Printf(`if len(v.%s) > %s {`, field.Name, maxVal)
		cb.Indent()
		cb.Printf(`verr.AddFieldError("%s", "%s must be at most %s characters", v.%s)`,
			field.JSONName, field.JSONName, maxVal, field.Name)
		cb.Dedent()
		cb.Writeln("}")
	}
	return nil
}

type LenRule struct{}

func (r LenRule) Name() string              { return "len" }
func (r LenRule) Priority() int             { return 2 }
func (r LenRule) RequiredImports() []string { return nil }
func (r LenRule) Aliases() []string         { return []string{} }

func (r LenRule) SupportsType(fieldType vtypes.FieldType) bool {
	return StringTypes.Contains(fieldType.Kind)
}

func (r LenRule) Generate(cb *builder.CodeBuilder, field vtypes.ValidationField, structName string) error {
	if lenVal, exists := field.Rules["len"]; exists {
		cb.Printf(`if len(v.%s) != %s {`, field.Name, lenVal)
		cb.Indent()
		cb.Printf(`verr.AddFieldError("%s", "%s must be exactly %s characters", v.%s)`,
			field.JSONName, field.JSONName, lenVal, field.Name)
		cb.Dedent()
		cb.Writeln("}")
	}
	return nil
}
