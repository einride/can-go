package messagenames

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
			Data: `BO_ 100 DriverHeartbeat: 1 DRIVER`,
		},

		{
			Name: "not ok",
			Data: `BO_ 100 DRIVER_HEARTBEAT: 1 DRIVER`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 1, Column: 1},
					Message: "message names must be CamelCase",
				},
			},
		},
	})
}
