package requireddefinitions

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
			Data: `
BS_:
BU_: ECU1
			`,
		},

		{
			Name: "missing bit timing",
			Data: `
BU_: ECU1
			`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 1, Column: 1},
					Message: "missing required definition(s)",
				},
			},
		},

		{
			Name: "missing nodes",
			Data: `
BS_:
			`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 1, Column: 1},
					Message: "missing required definition(s)",
				},
			},
		},
	})
}
