package newsymbols

import (
	"github.com/blueinnovationsgroup/can-go/pkg/dbc"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "newsymbols",
		Doc:  "check that the new symbols definition is empty",
		Run:  run,
	}
}

func run(pass *analysis.Pass) error {
	for _, def := range pass.File.Defs {
		newSymbolsDef, ok := def.(*dbc.NewSymbolsDef)
		if !ok {
			continue // not a new symbols definition
		}
		if len(newSymbolsDef.Symbols) > 0 {
			pass.Reportf(newSymbolsDef.Pos, "new symbols should be empty")
		}
	}
	return nil
}
