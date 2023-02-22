package siunits

import (
	"github.com/blueinnovationsgroup/can-go/pkg/dbc"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "unitsuffixes",
		Doc:  "check that signals with SI units have the correct symbols",
		Run:  run,
	}
}

const (
	metersPerSecond   = "m/s"
	kilometersPerHour = "km/h"
	meters            = "m"
	degrees           = "°"
	radians           = "rad"
)

// symbolMap returns a map from non-standard unit symbols to SI unit symbols.
func symbolMap() map[string]string {
	return map[string]string{
		"kph":        kilometersPerHour,
		"mps":        metersPerSecond,
		"meters/sec": metersPerSecond,
		"meters":     meters,
		"deg":        degrees,
		"degrees":    degrees,
		"radians":    radians,
	}
}

func run(pass *analysis.Pass) error {
	symbols := symbolMap()
	for _, def := range pass.File.Defs {
		message, ok := def.(*dbc.MessageDef)
		if !ok {
			continue
		}
		for _, signal := range message.Signals {
			if symbol, ok := symbols[signal.Unit]; ok {
				pass.Reportf(signal.Pos, "signal with unit %s should have SI unit %s", signal.Unit, symbol)
			}
		}
	}
	return nil
}
