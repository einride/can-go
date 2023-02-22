package definitiontypeorder

import (
	"testing"
	"text/scanner"

	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, Analyzer(), []*analysistest.Case{
		{
			Name: "correct order",
			Data: `
VERSION "foo"
NS_ :
BS_:
BU_:
			`,
		},

		{
			Name: "incorrect order",
			Data: `
VERSION "foo"
NS_ :
BU_:
BS_:
			`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 3, Column: 1},
					Message: "definition out of order",
				},
			},
		},

		{
			Name: "unknown defs last",
			Data: `
VERSION "foo"
NS_ :
BS_:
FOO "bar"
BU_:
			`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 4, Column: 1},
					Message: "definition out of order",
				},
			},
		},
	})
}
