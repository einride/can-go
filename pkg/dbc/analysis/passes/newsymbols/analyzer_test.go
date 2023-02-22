package newsymbols

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
			Data: `NS_ :`,
		},

		{
			Name: "not ok",
			Data: `
NS_ :
	BA_DEF_DEF_REL_
	BA_DEF_SGTYPE_`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 1, Column: 1},
					Message: "new symbols should be empty",
				},
			},
		},
	})
}
