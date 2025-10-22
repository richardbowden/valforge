package rules

import (
	"github.com/richardbowden/valforge/internal/builder"
	"github.com/richardbowden/valforge/internal/vtypes"
)

type EmailRule struct{}

func (r EmailRule) Name() string              { return "email" }
func (r EmailRule) Priority() int             { return 2 }
func (r EmailRule) RequiredImports() []string { return nil }
func (r EmailRule) Aliases() []string         { return []string{} }

func (r EmailRule) SupportsType(fieldType vtypes.FieldType) bool {
	return StringTypes.Contains(fieldType.Kind)
}

func (r EmailRule) Generate(cb *builder.CodeBuilder, field vtypes.ValidationField, structName string) error {

	cb.Printf("err := valgen.ValidateEmail(v.%s)", field.Name)
	cb.Printf("if err != nil {")
	cb.Indent()
	cb.Printf(`verr.AddFieldError("%s", err.Error(), v.%s)`, field.JSONName, field.Name)
	cb.Dedent()
	cb.Writeln("}")

	return nil
}
