package analysis

import (
	"fmt"
	"strings"
	"text/scanner"

	"github.com/blueinnovationsgroup/can-go/pkg/dbc"
)

// An Analyzer describes an analysis function and its options.
type Analyzer struct {
	// Name of the analyzer.
	Name string

	// Doc is the documentation for the analyzer.
	Doc string

	// Run the analyzer.
	Run func(*Pass) error
}

// Title is the part before the first "\n\n" of the documentation.
func (a *Analyzer) Title() string {
	return strings.SplitN(a.Doc, "\n\n", 2)[0]
}

// Validate the analyzer metadata.
func (a *Analyzer) Validate() error {
	if a.Doc == "" {
		return fmt.Errorf("missing doc")
	}
	return nil
}

// A Diagnostic is a message associated with a source location.
type Diagnostic struct {
	Pos     scanner.Position
	Message string
}

// Pass is the interface to the run function that analyzes DBC definitions.
type Pass struct {
	Analyzer    *Analyzer
	File        *dbc.File
	Diagnostics []*Diagnostic
}

// Reportf reports a diagnostic by building a message from the provided format and args.
func (pass *Pass) Reportf(pos scanner.Position, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	pass.Diagnostics = append(pass.Diagnostics, &Diagnostic{Pos: pos, Message: msg})
}
