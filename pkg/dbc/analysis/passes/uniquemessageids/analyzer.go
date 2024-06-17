package uniquemessageids

import (
	"go.einride.tech/can/pkg/dbc"
	"go.einride.tech/can/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "uniquemessageids",
		Doc:  "check that all message IDs are unique",
		Run:  run,
	}
}

func run(pass *analysis.Pass) error {
	messageIDs := make(map[dbc.MessageID]struct{})
	for _, def := range pass.File.Defs {
		message, ok := def.(*dbc.MessageDef)
		if !ok || dbc.IsIndependentSignalsMessage(message) {
			continue
		}
		if _, ok := messageIDs[message.MessageID]; ok {
			pass.Reportf(message.Pos, "non-unique message ID")
		} else {
			messageIDs[message.MessageID] = struct{}{}
		}
	}
	return nil
}
