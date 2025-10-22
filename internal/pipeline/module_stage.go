package pipeline

import (
	"path/filepath"

	"github.com/richardbowden/valforge/internal/modulegen"
	"github.com/richardbowden/valforge/internal/project"
	"github.com/richardbowden/valforge/internal/vfcontext"
)

type ValforgePackageStage struct{}

func (s *ValforgePackageStage) Name() string { return "Error Package" }

func (s *ValforgePackageStage) Execute(ctx *vfcontext.Context) error {
	// Find project root if not set
	if ctx.Config.ProjectRoot == "" {
		if ctx.Config.InputFile != "" {
			root, moduleName, err := project.FindProjectRoot(filepath.Dir(ctx.Config.InputFile))
			if err != nil {
				ctx.Config.ProjectRoot = "."
			} else {
				ctx.Config.ProjectRoot = root
				if ctx.Config.ModuleName == "" {
					ctx.Config.ModuleName = moduleName
				}
			}
		} else {
			ctx.Config.ProjectRoot = "."
		}
	}

	// Set default error package name
	if ctx.Config.ValforgePackage == "" {
		ctx.Config.ValforgePackage = "valgen"
	}

	// Generate error package
	errorGen := modulegen.NewGenerator(ctx.Config)
	return errorGen.EnsurePackages(ctx)
}
