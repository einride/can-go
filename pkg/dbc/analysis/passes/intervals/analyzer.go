package intervals

import (
	"github.com/blueinnovationsgroup/can-go/pkg/dbc"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "intervals",
		Doc:  "check that all intervals are valid (min <= max)",
		Run:  run,
	}
}

func run(pass *analysis.Pass) error {
	for _, def := range pass.File.Defs {
		switch def := def.(type) {
		case *dbc.EnvironmentVariableDef:
			if def.Minimum > def.Maximum {
				pass.Reportf(def.Pos, "invalid interval: [%f, %f]", def.Minimum, def.Maximum)
			}
		case *dbc.MessageDef:
			for i := range def.Signals {
				signal := &def.Signals[i]
				if signal.Minimum > signal.Maximum {
					pass.Reportf(def.Pos, "invalid interval: [%f, %f]", signal.Minimum, signal.Maximum)
				}
			}
		case *dbc.AttributeDef:
			if def.MinimumInt > def.MaximumInt || def.MinimumFloat > def.MaximumFloat {
				pass.Reportf(def.Pos, "invalid interval: [%d, %d]", def.MinimumInt, def.MaximumInt)
			}
			if def.MinimumFloat > def.MaximumFloat {
				pass.Reportf(def.Pos, "invalid interval: [%f, %f]", def.MinimumFloat, def.MaximumFloat)
			}
		}
	}
	return nil
}
