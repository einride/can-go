package singletondefinitions

import (
	"reflect"

	"github.com/blueinnovationsgroup/can-go/pkg/dbc"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "singletondefinitions",
		Doc:  "check that the file contains at most one of all singleton definitions",
		Run:  run,
	}
}

func singletonDefinitions() []dbc.Def {
	return []dbc.Def{
		&dbc.VersionDef{},
		&dbc.NewSymbolsDef{},
		&dbc.BitTimingDef{},
		&dbc.NodesDef{},
	}
}

func run(pass *analysis.Pass) error {
	defsByType := make(map[reflect.Type][]dbc.Def)
	for _, def := range pass.File.Defs {
		t := reflect.TypeOf(def)
		defsByType[t] = append(defsByType[t], def)
	}
	for _, singletonDef := range singletonDefinitions() {
		singletonDefs := defsByType[reflect.TypeOf(singletonDef)]
		if len(singletonDefs) > 1 {
			for i := 1; i < len(singletonDefs); i++ {
				pass.Reportf(singletonDefs[i].Position(), "more than one definition not allowed")
			}
		}
	}
	return nil
}
