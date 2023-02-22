package uniquenodenames

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
			Data: `BU_: ECU1 ECU2 ECU3`,
		},

		{
			Name: "duplicates",
			Data: `BU_: ECU1 ECU2 ECU3 ECU1`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 1, Column: 1},
					Message: "non-unique node name",
				},
			},
		},
	})
}
