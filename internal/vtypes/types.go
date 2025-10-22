package vtypes

import "go/types"

// FieldType represents the type information for a struct field
type FieldType struct {
	GoType     types.Type // The actual Go type
	Kind       TypeKind   // Simplified type classification
	IsPointer  bool       // Whether it's a pointer type
	IsSlice    bool       // Whether it's a slice type
	Underlying types.Type // Underlying type for pointers/slices
}

type TypeKind int

const (
	TypeUnknown TypeKind = iota
	TypeString
	TypeInt
	TypeInt8
	TypeInt16
	TypeInt32
	TypeInt64
	TypeUint
	TypeUint8
	TypeUint16
	TypeUint32
	TypeUint64
	TypeFloat32
	TypeFloat64
	TypeBool
	TypeStruct
)

func (tk TypeKind) String() string {
	switch tk {
	case TypeString:
		return "string"
	case TypeInt:
		return "int"
	case TypeInt8:
		return "int8"
	case TypeInt16:
		return "int16"
	case TypeInt32:
		return "int32"
	case TypeInt64:
		return "int64"
	case TypeUint:
		return "uint"
	case TypeUint8:
		return "uint8"
	case TypeUint16:
		return "uint16"
	case TypeUint32:
		return "uint32"
	case TypeUint64:
		return "uint64"
	case TypeFloat32:
		return "float32"
	case TypeFloat64:
		return "float64"
	case TypeBool:
		return "bool"
	case TypeStruct:
		return "struct"
	default:
		return "unknown"
	}
}

// ValidationField represents a field with validation rules and type info
type ValidationField struct {
	Name     string
	Type     FieldType
	JSONName string
	Rules    map[string]string
}

// ValidationStruct represents a struct with validation
type ValidationStruct struct {
	Name        string
	PackageName string
	Fields      []ValidationField
}

// GenerationConfig holds configuration for code generation
type GenerationConfig struct {
	InputFile           string
	PackagePath         string
	OutputFile          string
	PackageName         string
	ValforgePackage     string // Name of the package for supporting code (default: "valforge")
	ValforgePackagePath string // Path to error package (e.g., "internal/valgen")
	ModuleName          string
	ProjectRoot         string // Project root directory
	Version             string
}

// CompilerError represents an error during compilation
type CompilerError struct {
	Type     ErrorType
	Message  string
	Field    string
	Struct   string
	Rule     string
	Position string
}

type ErrorType int

const (
	ErrorTypeIncompatible ErrorType = iota
	ErrorTypeMissing
	ErrorTypeInvalid
	ErrorTypeDuplicate
)

func (e CompilerError) Error() string {
	return e.Message
}

// CompilerErrors holds multiple compiler errors
type CompilerErrors []CompilerError

func (errs CompilerErrors) Error() string {
	if len(errs) == 0 {
		return "no errors"
	}
	if len(errs) == 1 {
		return errs[0].Error()
	}
	return errs[0].Error() + " (and more)"
}

func (errs *CompilerErrors) Add(err CompilerError) {
	*errs = append(*errs, err)
}

func (errs CompilerErrors) HasErrors() bool {
	return len(errs) > 0
}
