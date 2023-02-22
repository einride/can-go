package intervals

import (
	"testing"
	"text/scanner"

	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, Analyzer(), []*analysistest.Case{
		{
			Name: "attribute interval ok",
			Data: `BA_DEF_ "AttributeName" INT 0 10;`,
		},

		{
			Name: "attribute interval bad",
			Data: `BA_DEF_ "AttributeName" INT 10 0;`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 1, Column: 1},
					Message: "invalid interval: [10, 0]",
				},
			},
		},
	})
}
