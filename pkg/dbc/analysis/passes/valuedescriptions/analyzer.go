package valuedescriptions

import (
	"go.einride.tech/can/internal/identifiers"
	"go.einride.tech/can/pkg/dbc"
	"go.einride.tech/can/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "valuedescriptions",
		Doc:  "check that value descriptions are valid CamelCase",
		Run:  run,
	}
}

func run(pass *analysis.Pass) error {
	for _, def := range pass.File.Defs {
		var valueDescriptions []dbc.ValueDescriptionDef
		switch def := def.(type) {
		case *dbc.ValueTableDef:
			valueDescriptions = def.ValueDescriptions
		case *dbc.ValueDescriptionsDef:
			valueDescriptions = def.ValueDescriptions
		default:
			continue
		}
		for _, vd := range valueDescriptions {
			vd := vd
			if !identifiers.IsCamelCase(vd.Description) {
				pass.Reportf(vd.Pos, "value description must be CamelCase")
			}
		}
	}
	return nil
}
