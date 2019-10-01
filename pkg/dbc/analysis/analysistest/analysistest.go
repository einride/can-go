package analysistest

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.einride.tech/can/pkg/dbc"
	"go.einride.tech/can/pkg/dbc/analysis"
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
		require.NoError(t, p.Parse())
		pass := &analysis.Pass{
			Analyzer: a,
			File:     p.File(),
		}
		require.NoError(t, a.Run(pass))
		// allow omitting byte offsets and file names
		for _, d := range c.Diagnostics {
			d.Pos.Offset = 0
			d.Pos.Filename = c.Name
		}
		for _, d := range pass.Diagnostics {
			d.Pos.Offset = 0
		}
		require.Equal(t, c.Diagnostics, pass.Diagnostics)
	}
}
