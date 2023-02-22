package noreservedsignals

import (
	"strings"

	"github.com/blueinnovationsgroup/can-go/pkg/dbc"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "noreservedsignals",
		Doc:  `checks that no signals have the prefix "Reserved"`,
		Run:  run,
	}
}

func run(pass *analysis.Pass) error {
	for _, d := range pass.File.Defs {
		messageDef, ok := d.(*dbc.MessageDef)
		if !ok {
			continue
		}
		for _, signalDef := range messageDef.Signals {
			if strings.HasPrefix(string(signalDef.Name), "Reserved") {
				pass.Reportf(signalDef.Pos, "remove reserved signals")
			}
		}
	}
	return nil
}
