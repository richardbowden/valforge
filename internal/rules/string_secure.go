package rules

import (
	"github.com/richardbowden/valforge/internal/builder"
	"github.com/richardbowden/valforge/internal/vtypes"
)

type EqualFieldSecureRule struct{}

func (r EqualFieldSecureRule) Name() string              { return "eqfieldsecure" }
func (r EqualFieldSecureRule) Priority() int             { return 5 }
func (r EqualFieldSecureRule) RequiredImports() []string { return []string{"crypto/subtle"} }
func (r EqualFieldSecureRule) Aliases() []string         { return []string{} }

func (r EqualFieldSecureRule) SupportsType(fieldType vtypes.FieldType) bool {
	return StringTypes.Contains(fieldType.Kind)
}

func (r EqualFieldSecureRule) Generate(cb *builder.CodeBuilder, field vtypes.ValidationField, structName string) error {
	if targetField, exists := field.Rules["eqfieldsecure"]; exists {
		cb.Printf(`if subtle.ConstantTimeCompare([]byte(v.%s), []byte(v.%s)) == 0 {`, field.Name, targetField)
		cb.Indent()
		cb.Printf(`verr.AddFieldError("%s", "%s must match %s", v.%s)`, field.JSONName, field.JSONName, targetField, field.Name)
		cb.Dedent()
		cb.Writeln("}")
	}
	return nil
}
