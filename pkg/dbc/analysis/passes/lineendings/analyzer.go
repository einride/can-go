package lineendings

import (
	"bytes"
	"text/scanner"

	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "lineendings",
		Doc:  `check that the file does not contain Windows line-endings (\r\n)`,
		Run:  run,
	}
}

func run(pass *analysis.Pass) error {
	if bytes.Contains(pass.File.Data, []byte{'\r', '\n'}) {
		pass.Reportf(
			scanner.Position{Filename: pass.File.Name, Line: 1, Column: 1},
			`file must not contain Windows line-endings (\r\n)`,
		)
	}
	return nil
}
