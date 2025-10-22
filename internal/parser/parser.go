package parser

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"strings"

	"github.com/richardbowden/valforge/internal/vtypes"
)

type Parser struct {
	fset *token.FileSet
	info *types.Info
}

func New() *Parser {
	return &Parser{
		fset: token.NewFileSet(),
		info: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
			Defs:  make(map[*ast.Ident]types.Object),
			Uses:  make(map[*ast.Ident]types.Object),
		},
	}
}

func (p *Parser) ParseFile(filePath string) ([]vtypes.ValidationStruct, string, error) {
	file, err := parser.ParseFile(p.fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, "", err
	}

	config := types.Config{
		Importer: importer.Default(),
		Error: func(err error) {
			// Ignore import errors for generated files
		},
	}
	pkg, err := config.Check("", p.fset, []*ast.File{file}, p.info)
	if err != nil {
		// Try without type checking if import resolution fails
		visitor := &structVisitor{
			info:        nil, // No type info available
			structs:     []vtypes.ValidationStruct{},
			packageName: file.Name.Name,
		}
		ast.Walk(visitor, file)
		return visitor.structs, file.Name.Name, nil
	}

	visitor := &structVisitor{
		info:        p.info,
		structs:     []vtypes.ValidationStruct{},
		packageName: pkg.Name(),
	}

	ast.Walk(visitor, file)
	return visitor.structs, visitor.packageName, nil
}

func (p *Parser) ParsePackage(packagePath string) ([]vtypes.ValidationStruct, string, error) {
	pkgs, err := parser.ParseDir(p.fset, packagePath, nil, parser.ParseComments)
	if err != nil {
		return nil, "", err
	}

	var allFiles []*ast.File
	var packageName string

	// Collect all files from the package
	for _, pkg := range pkgs {
		packageName = pkg.Name
		for _, file := range pkg.Files {
			allFiles = append(allFiles, file)
		}
	}

	if len(allFiles) == 0 {
		return nil, "", fmt.Errorf("no Go files found in package %s", packagePath)
	}

	// Type check the entire package
	config := types.Config{
		Importer: importer.Default(),
		Error: func(err error) {
			// Ignore import errors for generated files
		},
	}
	_, err = config.Check(packageName, p.fset, allFiles, p.info)

	var allStructs []vtypes.ValidationStruct
	useTypeInfo := err == nil

	// Visit all files to collect structs
	for _, file := range allFiles {
		visitor := &structVisitor{
			info:        nil,
			structs:     []vtypes.ValidationStruct{},
			packageName: packageName,
		}
		if useTypeInfo {
			visitor.info = p.info
		}
		ast.Walk(visitor, file)
		allStructs = append(allStructs, visitor.structs...)
	}

	return allStructs, packageName, nil
}

type structVisitor struct {
	info        *types.Info
	structs     []vtypes.ValidationStruct
	packageName string
}

func (v *structVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.TypeSpec:
		if structType, ok := n.Type.(*ast.StructType); ok {
			if s := v.parseStruct(n.Name.Name, structType, n); s != nil {
				v.structs = append(v.structs, *s)
			}
		}
	}
	return v
}

func (v *structVisitor) parseStruct(name string, structType *ast.StructType, typeSpec *ast.TypeSpec) *vtypes.ValidationStruct {
	var fields []vtypes.ValidationField
	hasValidation := false

	for _, field := range structType.Fields.List {
		if field.Tag == nil {
			continue
		}

		tagValue := strings.Trim(field.Tag.Value, "`")
		tags := parseStructTags(tagValue)

		if validateTag, exists := tags["validate"]; exists {
			hasValidation = true
			for _, fieldName := range field.Names {
				fieldType := v.extractFieldType(field.Type)
				fields = append(fields, vtypes.ValidationField{
					Name:     fieldName.Name,
					Type:     fieldType,
					JSONName: getJSONName(tags, fieldName.Name),
					Rules:    parseValidationRules(validateTag),
				})
			}
		}
	}

	if !hasValidation {
		return nil
	}

	return &vtypes.ValidationStruct{
		Name:        name,
		PackageName: v.packageName,
		Fields:      fields,
	}
}

func (v *structVisitor) extractFieldType(expr ast.Expr) vtypes.FieldType {
	ft := vtypes.FieldType{}

	// Get type info from type checker if available
	if v.info != nil {
		if typeInfo, ok := v.info.Types[expr]; ok {
			ft.GoType = typeInfo.Type
			ft.Kind = classifyType(typeInfo.Type)
		}
	}

	// Handle pointers and slices - fallback to AST analysis if no type info
	switch t := expr.(type) {
	case *ast.StarExpr:
		ft.IsPointer = true
		if v.info != nil {
			if innerType, ok := v.info.Types[t.X]; ok {
				ft.Underlying = innerType.Type
				ft.Kind = classifyType(innerType.Type)
			}
		} else {
			// Fallback: analyze AST structure
			ft.Kind = inferTypeFromAST(t.X)
		}
	case *ast.ArrayType:
		ft.IsSlice = true
		if v.info != nil {
			if innerType, ok := v.info.Types[t.Elt]; ok {
				ft.Underlying = innerType.Type
				ft.Kind = classifyType(innerType.Type)
			}
		} else {
			// Fallback: analyze AST structure
			ft.Kind = inferTypeFromAST(t.Elt)
		}
	default:
		if ft.Kind == vtypes.TypeUnknown {
			// Fallback: infer from AST if type checking failed
			ft.Kind = inferTypeFromAST(expr)
		}
	}

	return ft
}

// inferTypeFromAST attempts to infer the type from AST structure when type checking fails
func inferTypeFromAST(expr ast.Expr) vtypes.TypeKind {
	switch t := expr.(type) {
	case *ast.Ident:
		switch t.Name {
		case "string":
			return vtypes.TypeString
		case "int":
			return vtypes.TypeInt
		case "int8":
			return vtypes.TypeInt8
		case "int16":
			return vtypes.TypeInt16
		case "int32":
			return vtypes.TypeInt32
		case "int64":
			return vtypes.TypeInt64
		case "uint":
			return vtypes.TypeUint
		case "uint8":
			return vtypes.TypeUint8
		case "uint16":
			return vtypes.TypeUint16
		case "uint32":
			return vtypes.TypeUint32
		case "uint64":
			return vtypes.TypeUint64
		case "float32":
			return vtypes.TypeFloat32
		case "float64":
			return vtypes.TypeFloat64
		case "bool":
			return vtypes.TypeBool
		default:
			return vtypes.TypeUnknown
		}
	case *ast.StarExpr:
		return inferTypeFromAST(t.X)
	case *ast.ArrayType:
		return inferTypeFromAST(t.Elt)
	default:
		return vtypes.TypeUnknown
	}
}

func classifyType(t types.Type) vtypes.TypeKind {
	if t == nil {
		return vtypes.TypeUnknown
	}

	basic, ok := t.Underlying().(*types.Basic)
	if !ok {
		return vtypes.TypeUnknown
	}

	switch basic.Kind() {
	case types.String:
		return vtypes.TypeString
	case types.Int:
		return vtypes.TypeInt
	case types.Int8:
		return vtypes.TypeInt8
	case types.Int16:
		return vtypes.TypeInt16
	case types.Int32:
		return vtypes.TypeInt32
	case types.Int64:
		return vtypes.TypeInt64
	case types.Uint:
		return vtypes.TypeUint
	case types.Uint8:
		return vtypes.TypeUint8
	case types.Uint16:
		return vtypes.TypeUint16
	case types.Uint32:
		return vtypes.TypeUint32
	case types.Uint64:
		return vtypes.TypeUint64
	case types.Float32:
		return vtypes.TypeFloat32
	case types.Float64:
		return vtypes.TypeFloat64
	case types.Bool:
		return vtypes.TypeBool
	default:
		return vtypes.TypeUnknown
	}
}

func parseStructTags(tag string) map[string]string {
	tags := make(map[string]string)
	parts := strings.Fields(tag)

	for _, part := range parts {
		colonIndex := strings.Index(part, ":")
		if colonIndex == -1 {
			continue
		}

		key := part[:colonIndex]
		value := part[colonIndex+1:]

		if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		}

		tags[key] = value
	}

	return tags
}

func getJSONName(tags map[string]string, fieldName string) string {
	if jsonTag, exists := tags["json"]; exists {
		parts := strings.Split(jsonTag, ",")
		if len(parts) > 0 && parts[0] != "" && parts[0] != "-" {
			return parts[0]
		}
	}
	return toSnakeCase(fieldName)
}

func parseValidationRules(validateTag string) map[string]string {
	rules := make(map[string]string)
	parts := strings.Split(validateTag, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "=") {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 {
				rules[kv[0]] = kv[1]
			}
		} else {
			rules[part] = ""
		}
	}

	return rules
}

func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
