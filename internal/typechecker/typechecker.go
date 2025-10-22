package typechecker

import (
	"fmt"
	"strconv"

	"github.com/richardbowden/valforge/internal/vtypes"
)

type TypeChecker struct {
	registry interface {
		GetForTypeCheck(name string) (interface{ SupportsType(vtypes.FieldType) bool }, bool)
	}
}

func New(registry interface {
	GetForTypeCheck(name string) (interface{ SupportsType(vtypes.FieldType) bool }, bool)
}) *TypeChecker {
	return &TypeChecker{
		registry: registry,
	}
}

func (tc *TypeChecker) CheckStruct(s vtypes.ValidationStruct) vtypes.CompilerErrors {
	var errors vtypes.CompilerErrors

	// Build field map for cross-field validations
	fieldMap := make(map[string]vtypes.ValidationField)
	for _, field := range s.Fields {
		fieldMap[field.Name] = field
	}

	for _, field := range s.Fields {
		fieldErrors := tc.checkField(field, s.Name, fieldMap)
		errors = append(errors, fieldErrors...)
	}

	return errors
}

func (tc *TypeChecker) checkField(field vtypes.ValidationField, structName string, fieldMap map[string]vtypes.ValidationField) vtypes.CompilerErrors {
	var errors vtypes.CompilerErrors

	for ruleName, ruleValue := range field.Rules {
		rule, exists := tc.registry.GetForTypeCheck(ruleName)
		if !exists {
			errors.Add(vtypes.CompilerError{
				Type:    vtypes.ErrorTypeMissing,
				Message: fmt.Sprintf("unknown validation rule '%s'", ruleName),
				Field:   field.Name,
				Struct:  structName,
				Rule:    ruleName,
			})
			continue
		}

		// Check if rule is compatible with field type
		if !rule.SupportsType(field.Type) {
			errors.Add(vtypes.CompilerError{
				Type:    vtypes.ErrorTypeIncompatible,
				Message: fmt.Sprintf("rule '%s' is not compatible with type '%s'", ruleName, field.Type.Kind),
				Field:   field.Name,
				Struct:  structName,
				Rule:    ruleName,
			})
			continue
		}

		// Validate rule parameters
		if err := tc.validateRuleParams(ruleName, field, ruleValue, structName, fieldMap); err != nil {
			errors.Add(*err)
		}
	}

	return errors
}

func (tc *TypeChecker) validateRuleParams(ruleName string, field vtypes.ValidationField, ruleValue, structName string, fieldMap map[string]vtypes.ValidationField) *vtypes.CompilerError {
	switch ruleName {
	case "gt", "lt", "lte", "gte":
		if ruleValue == "" {
			return &vtypes.CompilerError{
				Type:    vtypes.ErrorTypeInvalid,
				Message: fmt.Sprintf("rule '%s' requires a numeric value", ruleName),
				Field:   field.Name,
				Struct:  structName,
				Rule:    ruleName,
			}
		}

		if !tc.isIntegerType(field.Type.Kind) {
			return &vtypes.CompilerError{
				Type:    vtypes.ErrorTypeIncompatible,
				Message: fmt.Sprintf("rule '%s' can only be used with integer types", ruleName),
				Field:   field.Name,
				Struct:  structName,
				Rule:    ruleName,
			}
		}

		if _, err := strconv.ParseInt(ruleValue, 10, 64); err != nil {
			return &vtypes.CompilerError{
				Type:    vtypes.ErrorTypeInvalid,
				Message: fmt.Sprintf("rule '%s' value must be a valid integer", ruleName),
				Field:   field.Name,
				Struct:  structName,
				Rule:    ruleName,
			}
		}
	case "minlen", "maxlen", "len":
		if ruleValue == "" {
			return &vtypes.CompilerError{
				Type:    vtypes.ErrorTypeInvalid,
				Message: fmt.Sprintf("rule '%s' requires a numeric value", ruleName),
				Field:   field.Name,
				Struct:  structName,
				Rule:    ruleName,
			}
		}

		if val, err := strconv.Atoi(ruleValue); err != nil || val < 0 {
			return &vtypes.CompilerError{
				Type:    vtypes.ErrorTypeInvalid,
				Message: fmt.Sprintf("rule '%s' value must be a non-negative integer", ruleName),
				Field:   field.Name,
				Struct:  structName,
				Rule:    ruleName,
			}
		}
	case "eqfield":
		if ruleValue == "" {
			return &vtypes.CompilerError{
				Type:    vtypes.ErrorTypeInvalid,
				Message: "eqfield rule requires a field name",
				Field:   field.Name,
				Struct:  structName,
				Rule:    ruleName,
			}
		}

		targetField, exists := fieldMap[ruleValue]
		if !exists {
			return &vtypes.CompilerError{
				Type:    vtypes.ErrorTypeMissing,
				Message: fmt.Sprintf("eqfield references unknown field '%s'", ruleValue),
				Field:   field.Name,
				Struct:  structName,
				Rule:    ruleName,
			}
		}

		if field.Type.Kind != targetField.Type.Kind {
			return &vtypes.CompilerError{
				Type:    vtypes.ErrorTypeIncompatible,
				Message: fmt.Sprintf("eqfield field types must match: '%s' vs '%s'", field.Type.Kind, targetField.Type.Kind),
				Field:   field.Name,
				Struct:  structName,
				Rule:    ruleName,
			}
		}
	case "eqfieldsecure":
		if ruleValue == "" {
			return &vtypes.CompilerError{
				Type:    vtypes.ErrorTypeInvalid,
				Message: "eqfieldsecure rule requires a field name",
				Field:   field.Name,
				Struct:  structName,
				Rule:    ruleName,
			}
		}

		targetField, exists := fieldMap[ruleValue]
		if !exists {
			return &vtypes.CompilerError{
				Type:    vtypes.ErrorTypeMissing,
				Message: fmt.Sprintf("eqfieldsecure references unknown field '%s'", ruleValue),
				Field:   field.Name,
				Struct:  structName,
				Rule:    ruleName,
			}
		}

		if field.Type.Kind != vtypes.TypeString || targetField.Type.Kind != vtypes.TypeString {
			return &vtypes.CompilerError{
				Type:    vtypes.ErrorTypeIncompatible,
				Message: "eqfieldsecure can only be used with string fields",
				Field:   field.Name,
				Struct:  structName,
				Rule:    ruleName,
			}
		}
	}
	return nil
}

func (tc *TypeChecker) isIntegerType(kind vtypes.TypeKind) bool {
	switch kind {
	case vtypes.TypeInt, vtypes.TypeInt8, vtypes.TypeInt16, vtypes.TypeInt32, vtypes.TypeInt64,
		vtypes.TypeUint, vtypes.TypeUint8, vtypes.TypeUint16, vtypes.TypeUint32, vtypes.TypeUint64:
		return true
	default:
		return false
	}
}
