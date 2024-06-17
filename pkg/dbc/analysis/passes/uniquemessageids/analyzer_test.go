package uniquemessageids

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
			Data: `
BO_ 101 MOTOR_CMD: 1 DRIVER
 SG_ MOTOR_CMD_steer : 0|4@1- (1,-5) [-5|5] "" MOTOR
 SG_ MOTOR_CMD_drive : 4|4@1+ (1,0) [0|9] "" MOTOR
BO_ 102 MOTOR_CMD: 1 DRIVER
 SG_ MOTOR_CMD_steer : 0|4@1- (1,-5) [-5|5] "" MOTOR
 SG_ MOTOR_CMD_drive : 4|4@1+ (1,0) [0|9] "" MOTOR
			`,
		},

		{
			Name: "duplicate",
			Data: `
BO_ 103 MOTOR_CMD: 1 DRIVER
 SG_ MOTOR_CMD_steer : 0|4@1- (1,-5) [-5|5] "" MOTOR
 SG_ MOTOR_CMD_steer : 4|4@1+ (1,0) [0|9] "" MOTOR
BO_ 103 MOTOR_CMD: 1 DRIVER
 SG_ MOTOR_CMD_steer : 0|4@1- (1,-5) [-5|5] "" MOTOR
 SG_ MOTOR_CMD_steer : 4|4@1+ (1,0) [0|9] "" MOTOR
			`,
			Diagnostics: []*analysis.Diagnostic{
				{
					Pos:     scanner.Position{Line: 4, Column: 1},
					Message: "non-unique message ID",
				},
			},
		},
	})
}
