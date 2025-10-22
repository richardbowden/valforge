package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/richardbowden/valforge/internal/pipeline"
	"github.com/richardbowden/valforge/internal/rules"
	"github.com/richardbowden/valforge/internal/vfcontext"
	"github.com/richardbowden/valforge/internal/vtypes"
)

var showVersion bool

type Stage interface {
	Name() string
	Execute(ctx *vfcontext.Context) error
}

type Pipeline struct {
	stages []Stage
}

func New() *Pipeline {
	return &Pipeline{}
}

func (p *Pipeline) AddStage(stage Stage) {
	p.stages = append(p.stages, stage)
}

func (p *Pipeline) Execute(ctx *vfcontext.Context) error {
	for _, stage := range p.stages {
		fmt.Printf("→ %s\n", stage.Name())

		if err := stage.Execute(ctx); err != nil {
			return fmt.Errorf("stage %s failed: %w", stage.Name(), err)
		}
	}
	return nil
}

func setupRegistry() *rules.Registry {
	registry := rules.NewRegistry()
	registry.Register(&rules.RequiredRule{})
	registry.Register(&rules.GreaterThanRule{})
	registry.Register(&rules.LessThanRule{})
	registry.Register(&rules.EqualFieldRule{})

	registry.Register(&rules.MinLenRule{})
	registry.Register(&rules.MaxLenRule{})
	registry.Register(&rules.LenRule{})
	registry.Register(&rules.EqualFieldSecureRule{})

	registry.Register(&rules.EmailRule{})
	return registry
}

func main() {
	config := parseFlags()

	if showVersion {
		fmt.Printf("valforge: %s\n", GetVersion())
		os.Exit(0)
	}

	registry := setupRegistry()
	pipe := New()

	pipe.AddStage(&pipeline.ParseStage{})
	pipe.AddStage(&pipeline.TypeCheckStage{})
	pipe.AddStage(&pipeline.ValforgePackageStage{})
	pipe.AddStage(&pipeline.GenerateStage{})
	pipe.AddStage(&pipeline.WriteStage{})

	ctx := &vfcontext.Context{
		Config:   config,
		Registry: registry,
	}

	if err := pipe.Execute(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✓ Generated validation for %d structs in %s\n",
		len(ctx.Structs), config.OutputFile)
}

func parseFlags() vtypes.GenerationConfig {
	var config vtypes.GenerationConfig

	flag.StringVar(&config.InputFile, "file", "", "Go file to scan for validation tags")
	flag.StringVar(&config.PackagePath, "package", "", "Package directory to scan for validation tags")
	flag.StringVar(&config.OutputFile, "output", "", "Output file")
	flag.StringVar(&config.ValforgePackage, "valforge-package", "valgen", "Name of the valfore supporting code package")
	flag.StringVar(&config.ValforgePackagePath, "valforge-path", "", "Path to error package (default: internal/valgen)")
	flag.BoolVar(&showVersion, "version", false, "shows version then exits")
	flag.Parse()

	config.Version = GetVersion()
	return config
}
