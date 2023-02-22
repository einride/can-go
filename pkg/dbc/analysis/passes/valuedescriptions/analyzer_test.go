package valuedescriptions

import (
	"testing"
	"text/scanner"

	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, Analyzer(), []*analysistest.Case{
		{
			Name: "ok",
			Data: `VAL_ 100 Command 2 "Reboot" 1 "Sync" 0 "Noop";`,
		},
		{
			Name: "ok",
			Data: `VAL_ 100 Command 2 "11Reboot" 1 "123" 0 "Noop";`,
		},
		{
			Name: "underscore",
			Data: `VAL_ 100 Command 2 "Reboot_Command" 1 "Sync" 0 "Noop";`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 1, Column: 21},
					Message: "value description must be CamelCase (numbers ignored)",
				},
			},
		},
		{
			Name: "several digits value",
			Data: `VAL_ 100 Command 234 "Reboot_Command" 1 "Sync" 0 "Noop";`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 1, Column: 23},
					Message: "value description must be CamelCase (numbers ignored)",
				},
			},
		},
	})
}
