package vfcontext

import (
	"github.com/richardbowden/valforge/internal/builder"
	"github.com/richardbowden/valforge/internal/vtypes"
)

type PackageOptions uint8

const (
	EmailPO = 1 << iota
)

type Context struct {
	Config   vtypes.GenerationConfig
	Registry interface {
		GetRequiredImports(fields []vtypes.ValidationField) []string
		GetForTypeCheck(name string) (interface{ SupportsType(vtypes.FieldType) bool }, bool)
		GetAllForGeneration() map[string]interface {
			Generate(*builder.CodeBuilder, vtypes.ValidationField, string) error
			Priority() int
		}
	}
	Structs        []vtypes.ValidationStruct
	Output         string
	Errors         vtypes.CompilerErrors
	PackageOptions PackageOptions
}
