package rules

import (
	"github.com/richardbowden/valforge/internal/builder"
	"github.com/richardbowden/valforge/internal/vtypes"
)

type RequiredRule struct{}

func (r RequiredRule) Name() string              { return "required" }
func (r RequiredRule) Priority() int             { return 1 }
func (r RequiredRule) RequiredImports() []string { return nil }
func (r RequiredRule) Aliases() []string         { return []string{} }
func (r RequiredRule) SupportsType(fieldType vtypes.FieldType) bool {
	return StringTypes.Contains(fieldType.Kind) || IntegerTypes.Contains(fieldType.Kind)
}

func (r RequiredRule) Generate(cb *builder.CodeBuilder, field vtypes.ValidationField, structName string) error {
	switch field.Type.Kind {
	case vtypes.TypeString:
		cb.Printf(`if v.%s == "" {`, field.Name)
		cb.Indent()
		cb.Printf(`verr.AddFieldError("%s", "%s is required", v.%s)`,
			field.JSONName, field.JSONName, field.Name)
		cb.Dedent()
		cb.Writeln("}")

	case vtypes.TypeInt, vtypes.TypeInt8, vtypes.TypeInt16, vtypes.TypeInt32, vtypes.TypeInt64,
		vtypes.TypeUint, vtypes.TypeUint8, vtypes.TypeUint16, vtypes.TypeUint32, vtypes.TypeUint64:
		cb.Printf(`if v.%s == 0 {`, field.Name)
		cb.Indent()
		cb.Printf(`verr.AddFieldError("%s", "%s is required", v.%s)`,
			field.JSONName, field.JSONName, field.Name)
		cb.Dedent()
		cb.Writeln("}")
	}

	return nil
}
