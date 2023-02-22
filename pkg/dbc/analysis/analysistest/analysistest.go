package analysistest

import (
	"strings"
	"testing"

	"github.com/blueinnovationsgroup/can-go/pkg/dbc"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
	"gotest.tools/v3/assert"
)

type Case struct {
	Name        string
	Data        string
	Diagnostics []*analysis.Diagnostic
}

func Run(t *testing.T, a *analysis.Analyzer, cs []*Case) {
	t.Helper()
	for _, c := range cs {
		p := dbc.NewParser(c.Name, []byte(strings.TrimLeft(c.Data, "\n")))
		assert.NilError(t, p.Parse())
		pass := &analysis.Pass{
			Analyzer: a,
			File:     p.File(),
		}
		assert.NilError(t, a.Run(pass))
		// allow omitting byte offsets and file names
		for _, d := range c.Diagnostics {
			d.Pos.Offset = 0
			d.Pos.Filename = c.Name
		}
		for _, d := range pass.Diagnostics {
			d.Pos.Offset = 0
		}
		assert.DeepEqual(t, c.Diagnostics, pass.Diagnostics)
	}
}
