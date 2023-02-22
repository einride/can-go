package boolprefix

import (
	"strings"

	"github.com/blueinnovationsgroup/can-go/pkg/dbc"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "boolprefix",
		Doc:  "check that bools (1-bit signals) have a correct prefix",
		Run:  run,
	}
}

func allowedPrefixes() []string {
	return []string{
		"Is",
		"Has",
	}
}

func run(pass *analysis.Pass) error {
	for _, d := range pass.File.Defs {
		messageDef, ok := d.(*dbc.MessageDef)
		if !ok {
			continue
		}
	SignalLoop:
		for _, signalDef := range messageDef.Signals {
			if signalDef.Size != 1 {
				continue // skip all non-bool signals
			}
			for _, allowedPrefix := range allowedPrefixes() {
				if strings.HasPrefix(string(signalDef.Name), allowedPrefix) {
					continue SignalLoop // has allowed prefix
				}
			}
			// edge-case: allow non-prefixed 1-bit signals with value descriptions
			for _, d := range pass.File.Defs {
				valueDescriptionsDef, ok := d.(*dbc.ValueDescriptionsDef)
				if !ok {
					continue // not value descriptions
				}
				if valueDescriptionsDef.MessageID == messageDef.MessageID &&
					valueDescriptionsDef.SignalName == signalDef.Name {
					continue SignalLoop // has value descriptions
				}
			}
			pass.Reportf(
				signalDef.Pos,
				"bool signals (1-bit) must have prefix %s",
				strings.Join(allowedPrefixes(), " or "),
			)
		}
	}
	return nil
}
