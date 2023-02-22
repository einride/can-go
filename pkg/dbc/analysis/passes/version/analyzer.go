package version

import (
	"github.com/blueinnovationsgroup/can-go/pkg/dbc"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "version",
		Doc:  "check that the version definition is empty",
		Run:  run,
	}
}

func run(pass *analysis.Pass) error {
	for _, def := range pass.File.Defs {
		versionDef, ok := def.(*dbc.VersionDef)
		if !ok {
			continue // not a version definition
		}
		if len(versionDef.Version) > 0 {
			pass.Reportf(versionDef.Pos, "version should be empty")
		}
	}
	return nil
}
