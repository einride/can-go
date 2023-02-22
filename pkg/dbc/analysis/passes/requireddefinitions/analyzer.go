package requireddefinitions

import (
	"reflect"

	"github.com/blueinnovationsgroup/can-go/pkg/dbc"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "requireddefinitions",
		Doc:  "check that the file contains exactly one of all required definitions",
		Run:  run,
	}
}

func requiredDefinitions() []dbc.Def {
	return []dbc.Def{
		&dbc.BitTimingDef{},
		&dbc.NodesDef{},
	}
}

func run(pass *analysis.Pass) error {
	counts := make(map[reflect.Type]int)
	for _, def := range pass.File.Defs {
		counts[reflect.TypeOf(def)]++
	}
	for _, requiredDef := range requiredDefinitions() {
		if counts[reflect.TypeOf(requiredDef)] == 0 {
			// we have no definition to return, so return the first
			pass.Reportf(pass.File.Defs[0].Position(), "missing required definition(s)")
			break
		}
	}
	return nil
}
