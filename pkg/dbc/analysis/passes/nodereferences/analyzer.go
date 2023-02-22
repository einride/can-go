package nodereferences

import (
	"github.com/blueinnovationsgroup/can-go/pkg/dbc"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "nodereferences",
		Doc:  "check that all node references refer to declared nodes",
		Run:  run,
	}
}

func run(pass *analysis.Pass) error {
	declaredNodes := map[dbc.Identifier]struct{}{
		dbc.NodePlaceholder: {}, // placeholder is implicitly declared
	}
	// collect declared nodes
	for _, def := range pass.File.Defs {
		if def, ok := def.(*dbc.NodesDef); ok {
			for _, nodeName := range def.NodeNames {
				declaredNodes[nodeName] = struct{}{}
			}
		}
	}
	// verify node references
	for _, def := range pass.File.Defs {
		switch def := def.(type) {
		case *dbc.MessageDef:
			if _, ok := declaredNodes[def.Transmitter]; !ok {
				pass.Reportf(def.Pos, "undeclared transmitter node: %v", def.Transmitter)
			}
			for i := range def.Signals {
				signal := &def.Signals[i]
				for _, receiver := range signal.Receivers {
					if _, ok := declaredNodes[receiver]; !ok {
						pass.Reportf(signal.Pos, "undeclared receiver node: %v", receiver)
					}
				}
			}
		case *dbc.EnvironmentVariableDef:
			for _, accessNode := range def.AccessNodes {
				if _, ok := declaredNodes[accessNode]; !ok {
					pass.Reportf(def.Pos, "undeclared access node: %v", accessNode)
				}
			}
		case *dbc.MessageTransmittersDef:
			for _, transmitter := range def.Transmitters {
				if _, ok := declaredNodes[transmitter]; !ok {
					pass.Reportf(def.Pos, "undeclared transmitter node: %v", transmitter)
				}
			}
		}
	}
	return nil
}
