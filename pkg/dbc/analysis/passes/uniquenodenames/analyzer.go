package uniquenodenames

import (
	"go.einride.tech/can/pkg/dbc"
	"go.einride.tech/can/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "uniquenodenames",
		Doc:  "check that all declared node names are unique",
		Run:  run,
	}
}

func run(pass *analysis.Pass) error {
	nodeNames := make(map[dbc.Identifier]struct{})
	for _, def := range pass.File.Defs {
		if def, ok := def.(*dbc.NodesDef); ok {
			for _, nodeName := range def.NodeNames {
				if _, ok := nodeNames[nodeName]; ok {
					pass.Reportf(def.Pos, "non-unique node name")
				}
				nodeNames[nodeName] = struct{}{}
			}
		}
	}
	return nil
}
