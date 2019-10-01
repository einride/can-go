package valuedescriptions

import (
	"testing"
	"text/scanner"

	"go.einride.tech/can/pkg/dbc/analysis"
	"go.einride.tech/can/pkg/dbc/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, Analyzer(), []*analysistest.Case{
		{
			Name: "ok",
			Data: `VAL_ 100 Command 2 "Reboot" 1 "Sync" 0 "Noop";`,
		},

		{
			Name: "underscore",
			Data: `VAL_ 100 Command 2 "Reboot_Command" 1 "Sync" 0 "Noop";`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 1, Column: 18},
					Message: "value description must be CamelCase",
				},
			},
		},
	})
}
