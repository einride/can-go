package uniquesignalnames

import (
	"github.com/blueinnovationsgroup/can-go/pkg/dbc"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "uniquesignalnames",
		Doc:  "check that all signal names are unique",
		Run:  run,
	}
}

func run(pass *analysis.Pass) error {
	for _, def := range pass.File.Defs {
		message, ok := def.(*dbc.MessageDef)
		if !ok || dbc.IsIndependentSignalsMessage(message) {
			continue
		}
		signalNames := make(map[dbc.Identifier]struct{})
		for i := range message.Signals {
			signal := &message.Signals[i]
			if _, ok := signalNames[signal.Name]; ok {
				pass.Reportf(signal.Pos, "non-unique signal name")
			} else {
				signalNames[signal.Name] = struct{}{}
			}
		}
	}
	return nil
}
