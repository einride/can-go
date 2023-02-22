package unitsuffixes

import (
	"strings"

	"github.com/blueinnovationsgroup/can-go/pkg/dbc"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "unitsuffixes",
		Doc:  "check that signals with units have correct name suffixes",
		Run:  run,
	}
}

func unitSuffixes() map[string]string {
	return map[string]string{
		"°":    "Degrees",
		"rad":  "Radians",
		"%":    "Percent",
		"km/h": "Kph",
		"m/s":  "Mps",
	}
}

func run(pass *analysis.Pass) error {
	suffixes := unitSuffixes()
	for _, def := range pass.File.Defs {
		message, ok := def.(*dbc.MessageDef)
		if !ok {
			continue
		}
		for _, signal := range message.Signals {
			if suffix, ok := suffixes[signal.Unit]; ok {
				if !strings.HasSuffix(string(signal.Name), suffix) {
					pass.Reportf(signal.Pos, "signal with unit %s must have suffix %s", signal.Unit, suffix)
				}
			}
		}
	}
	return nil
}
