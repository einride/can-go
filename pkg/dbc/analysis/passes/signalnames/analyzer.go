package signalnames

import (
	"github.com/blueinnovationsgroup/can-go/internal/identifiers"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "signalnames",
		Doc:  "check that signal names are valid CamelCase identifiers",
		Run:  run,
	}
}

func run(pass *analysis.Pass) error {
	for _, d := range pass.File.Defs {
		messageDef, ok := d.(*dbc.MessageDef)
		if !ok {
			continue
		}
		for _, signalDef := range messageDef.Signals {
			if !identifiers.IsCamelCase(string(signalDef.Name)) {
				pass.Reportf(signalDef.Pos, "signal names must be CamelCase")
			}
		}
	}
	return nil
}
