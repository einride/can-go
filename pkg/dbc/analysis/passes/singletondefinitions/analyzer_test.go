package singletondefinitions

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
VERSION "foo"
NS_:
BS_:
BU_: ECU1
			`,
		},

		{
			Name: "multiple versions",
			Data: `
VERSION "foo"
VERSION "foo"
NS_:
BS_:
BU_: ECU1
			`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 2, Column: 1},
					Message: "more than one definition not allowed",
				},
			},
		},

		{
			Name: "multiple new symbols",
			Data: `
VERSION "foo"
NS_:
NS_:
BS_:
BU_: ECU1
			`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 3, Column: 1},
					Message: "more than one definition not allowed",
				},
			},
		},

		{
			Name: "multiple bit timing",
			Data: `
VERSION "foo"
NS_:
BS_:
BS_:
BU_: ECU1
			`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 4, Column: 1},
					Message: "more than one definition not allowed",
				},
			},
		},

		{
			Name: "multiple nodes",
			Data: `
VERSION "foo"
NS_:
BS_:
BU_: ECU1
BU_: ECU2
			`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 5, Column: 1},
					Message: "more than one definition not allowed",
				},
			},
		},
	})
}
