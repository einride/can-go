package boolprefix

import (
	"testing"
	"text/scanner"

	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, Analyzer(), []*analysistest.Case{
		{
			Name: "prefix has",
			Data: `
BO_ 400 MOTOR_STATUS: 3 MOTOR
 SG_ HasWheelError : 0|1@1+ (1,0) [0|0] "" DRIVER,IO
`,
		},

		{
			Name: "prefix is",
			Data: `
BO_ 400 MOTOR_STATUS: 3 MOTOR
 SG_ IsOverheated : 0|1@1+ (1,0) [0|0] "" DRIVER,IO
`,
		},

		{
			Name: "missing prefix",
			Data: `
BO_ 400 MOTOR_STATUS: 3 MOTOR
 SG_ WheelError : 0|1@1+ (1,0) [0|0] "" DRIVER,IO
`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 2, Column: 2},
					Message: "bool signals (1-bit) must have prefix Is or Has",
				},
			},
		},

		{
			Name: "missing prefix with value descriptions",
			Data: `
BO_ 400 MOTOR_STATUS: 3 MOTOR
 SG_ Status : 0|1@1+ (1,0) [0|0] "" DRIVER,IO

VAL_ 400 Status 1 "ValidDataPresent" 0 "NoData" ;
`,
		},
	})
}
