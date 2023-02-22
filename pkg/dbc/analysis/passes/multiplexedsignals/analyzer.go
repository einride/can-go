package multiplexedsignals

import (
	"github.com/blueinnovationsgroup/can-go/pkg/dbc"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "multiplexedsignals",
		Doc:  "check that multiplexed signals are valid",
		Run:  run,
	}
}

func run(pass *analysis.Pass) error {
	for _, def := range pass.File.Defs {
		message, ok := def.(*dbc.MessageDef)
		if !ok {
			continue
		}
		// locate multiplexer switch
		var multiplexerSwitch *dbc.SignalDef
		for i := range message.Signals {
			if !message.Signals[i].IsMultiplexerSwitch {
				continue
			}
			if multiplexerSwitch != nil {
				pass.Reportf(message.Signals[i].Pos, "more than one multiplexer switch")
				continue
			}
			multiplexerSwitch = &message.Signals[i]
			if multiplexerSwitch.IsSigned {
				pass.Reportf(message.Signals[i].Pos, "signed multiplexer switch")
				continue
			}
			if multiplexerSwitch.IsMultiplexed {
				pass.Reportf(message.Signals[i].Pos, "can't be multiplexer and multiplexed")
				continue
			}
		}
		for i := range message.Signals {
			signal := &message.Signals[i]
			if !signal.IsMultiplexed {
				continue
			}
			if multiplexerSwitch == nil {
				pass.Reportf(message.Signals[i].Pos, "no multiplexer switch for multiplexed signal")
				continue
			}
			multiplexerSwitchMaxValue := uint64((1 << multiplexerSwitch.Size) - 1)
			if signal.MultiplexerSwitch > multiplexerSwitchMaxValue {
				pass.Reportf(signal.Pos, "multiplexer switch exceeds max value: %v", multiplexerSwitchMaxValue)
				continue
			}
		}
	}
	return nil
}
