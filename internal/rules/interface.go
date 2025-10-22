package rules

import (
	"github.com/richardbowden/valforge/internal/builder"
	"github.com/richardbowden/valforge/internal/vtypes"
)

type Rule interface {
	Name() string
	Aliases() []string
	Priority() int
	SupportsType(fieldType vtypes.FieldType) bool
	RequiredImports() []string
	Generate(cb *builder.CodeBuilder, field vtypes.ValidationField, structName string) error
}

type TypeSet []vtypes.TypeKind

func (ts TypeSet) Contains(kind vtypes.TypeKind) bool {
	for _, t := range ts {
		if t == kind {
			return true
		}
	}
	return false
}

var (
	StringTypes  = TypeSet{vtypes.TypeString}
	IntegerTypes = TypeSet{
		vtypes.TypeInt, vtypes.TypeInt8, vtypes.TypeInt16, vtypes.TypeInt32, vtypes.TypeInt64,
		vtypes.TypeUint, vtypes.TypeUint8, vtypes.TypeUint16, vtypes.TypeUint32, vtypes.TypeUint64,
	}
	AllTypes = TypeSet{
		vtypes.TypeString, vtypes.TypeInt, vtypes.TypeInt8, vtypes.TypeInt16, vtypes.TypeInt32, vtypes.TypeInt64,
		vtypes.TypeUint, vtypes.TypeUint8, vtypes.TypeUint16, vtypes.TypeUint32, vtypes.TypeUint64,
		vtypes.TypeFloat32, vtypes.TypeFloat64, vtypes.TypeBool,
	}
)
