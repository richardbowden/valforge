package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/richardbowden/valforge/internal/generator"
	"github.com/richardbowden/valforge/internal/parser"
	"github.com/richardbowden/valforge/internal/typechecker"
	"github.com/richardbowden/valforge/internal/vfcontext"
	"github.com/richardbowden/valforge/internal/vtypes"
)

type ParseStage struct{}

func (s *ParseStage) Name() string { return "Parse" }

func (s *ParseStage) Execute(ctx *vfcontext.Context) error {
	p := parser.New()

	var structs []vtypes.ValidationStruct
	var packageName string
	var err error

	if ctx.Config.InputFile != "" {
		structs, packageName, err = p.ParseFile(ctx.Config.InputFile)
	} else if ctx.Config.PackagePath != "" {
		structs, packageName, err = p.ParsePackage(ctx.Config.PackagePath)
	} else if goFile := os.Getenv("GOFILE"); goFile != "" {
		ctx.Config.InputFile = goFile
		structs, packageName, err = p.ParseFile(goFile)
	} else {
		return fmt.Errorf("no input file or package specified")
	}

	if err != nil {
		return err
	}

	if len(structs) == 0 {
		return fmt.Errorf("no structs with validation tags found")
	}

	ctx.Structs = structs
	ctx.Config.PackageName = packageName

	if ctx.Config.OutputFile == "" {
		if ctx.Config.InputFile != "" {
			dir := filepath.Dir(ctx.Config.InputFile)
			base := filepath.Base(ctx.Config.InputFile)
			name := strings.TrimSuffix(base, filepath.Ext(base))
			ctx.Config.OutputFile = filepath.Join(dir, name+"_validation.gen.go")
		} else {
			ctx.Config.OutputFile = filepath.Join(ctx.Config.PackagePath, "validation.gen.go")
		}
	}

	return nil
}

type TypeCheckStage struct{}

func (s *TypeCheckStage) Name() string { return "Type Check" }

func (s *TypeCheckStage) Execute(ctx *vfcontext.Context) error {
	tc := typechecker.New(ctx.Registry)

	for _, st := range ctx.Structs {
		errors := tc.CheckStruct(st)
		ctx.Errors = append(ctx.Errors, errors...)
	}

	if ctx.Errors.HasErrors() {
		return ctx.Errors
	}

	return nil
}

type GenerateStage struct{}

func (s *GenerateStage) Name() string { return "Generate" }

func (s *GenerateStage) Execute(ctx *vfcontext.Context) error {
	gen := generator.New(ctx.Registry, ctx.Config)

	output, err := gen.Generate(ctx.Structs)
	if err != nil {
		return err
	}

	ctx.Output = output
	return nil
}

type WriteStage struct{}

func (s *WriteStage) Name() string { return "Write" }

func (s *WriteStage) Execute(ctx *vfcontext.Context) error {
	file, err := os.Create(ctx.Config.OutputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(ctx.Output)
	return err
}
