package builder

import (
	"fmt"
	"strings"
)

// CodeBuilder helps build code with proper indentation
type CodeBuilder struct {
	content     strings.Builder
	indentLevel int
}

// NewCodeBuilder creates a new code builder
func NewCodeBuilder() *CodeBuilder {
	return &CodeBuilder{
		indentLevel: 0,
	}
}

// Write writes a string without indentation or newline
func (cb *CodeBuilder) Write(s string) {
	cb.content.WriteString(s)
}

// Writeln writes a line with proper indentation
func (cb *CodeBuilder) Writeln(s string) {
	for i := 0; i < cb.indentLevel; i++ {
		cb.content.WriteString("    ")
	}
	cb.content.WriteString(s)
	cb.content.WriteString("\n")
}

// Printf writes a formatted line with indentation
func (cb *CodeBuilder) Printf(format string, args ...any) {
	cb.Writeln(fmt.Sprintf(format, args...))
}

// Indent increases indentation level
func (cb *CodeBuilder) Indent() {
	cb.indentLevel++
}

// Dedent decreases indentation level
func (cb *CodeBuilder) Dedent() {
	if cb.indentLevel > 0 {
		cb.indentLevel--
	}
}

// Newline adds an empty line
func (cb *CodeBuilder) Newline() {
	cb.content.WriteString("\n")
}

// String returns the built string
func (cb *CodeBuilder) String() string {
	return cb.content.String()
}

