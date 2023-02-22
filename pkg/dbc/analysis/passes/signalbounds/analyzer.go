package signalbounds

import (
	"github.com/blueinnovationsgroup/can-go/pkg/dbc"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "signalbounds",
		Doc:  "check that signal start and end bits are within bounds of the message size",
		Run:  run,
	}
}

func run(pass *analysis.Pass) error {
	for _, def := range pass.File.Defs {
		message, ok := def.(*dbc.MessageDef)
		if !ok || dbc.IsIndependentSignalsMessage(message) {
			continue
		}
		for i := range message.Signals {
			signal := &message.Signals[i]
			if signal.StartBit >= 8*message.Size {
				pass.Reportf(signal.Pos, "start bit out of bounds")
			}
			// TODO: Check end bit
		}
	}
	return nil
}
