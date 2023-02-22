package nodereferences

import (
	"testing"
	"text/scanner"

	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, Analyzer(), []*analysistest.Case{
		{
			Name: "valid",
			Data: `
BU_: ECU1 ECU2
BO_ 42 TestMessage: 8 ECU2
 SG_ CellTempLowest : 32|8@0+ (1,-40) [-40|215] "C" ECU1
			`,
		},

		{
			Name: "undeclared transmitter",
			Data: `
BU_: ECU1 ECU2
BO_ 42 TestMessage: 8 ECU3
 SG_ CellTempLowest : 32|8@0+ (1,-40) [-40|215] "C" ECU1
			`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 2, Column: 1},
					Message: "undeclared transmitter node: ECU3",
				},
			},
		},

		{
			Name: "undeclared receiver",
			Data: `
BU_: ECU1 ECU2
BO_ 42 TestMessage: 8 ECU2
 SG_ CellTempLowest : 32|8@0+ (1,-40) [-40|215] "C" ECU2,ECU3
			`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 3, Column: 2},
					Message: "undeclared receiver node: ECU3",
				},
			},
		},
	})
}
