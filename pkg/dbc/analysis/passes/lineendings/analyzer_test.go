package lineendings

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
			Data: `NS_ :`,
		},

		{
			Name: "not ok",
			Data: "NS_ :\r\n",
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 1, Column: 1},
					Message: `file must not contain Windows line-endings (\r\n)`,
				},
			},
		},
	})
}
