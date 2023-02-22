package definitiontypeorder

import (
	"math"
	"reflect"

	"github.com/blueinnovationsgroup/can-go/pkg/dbc"
	"github.com/blueinnovationsgroup/can-go/pkg/dbc/analysis"
)

func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "definitiontypeorder",
		Doc:  "check that definitions are in the correct order",
		Run:  run,
	}
}

func orderOf(def dbc.Def) uint64 {
	for i, orderDef := range []dbc.Def{
		&dbc.VersionDef{},
		&dbc.NewSymbolsDef{},
		&dbc.BitTimingDef{},
		&dbc.NodesDef{},
		&dbc.ValueTableDef{},
		&dbc.MessageDef{},
		&dbc.MessageTransmittersDef{},
		&dbc.EnvironmentVariableDef{},
		&dbc.EnvironmentVariableDataDef{},
		&dbc.CommentDef{},
		&dbc.AttributeDef{},
		&dbc.AttributeDefaultValueDef{},
		&dbc.AttributeValueForObjectDef{},
		&dbc.ValueDescriptionsDef{},
	} {
		if reflect.TypeOf(def) == reflect.TypeOf(orderDef) {
			return uint64(i)
		}
	}
	return math.MaxUint64
}

func run(pass *analysis.Pass) error {
	minOrder := uint64(math.MaxUint64)
	for i := range pass.File.Defs {
		// diagnostics make more sense when going backwards
		def := pass.File.Defs[len(pass.File.Defs)-i-1]
		currOrder := orderOf(def)
		if currOrder > minOrder {
			pass.Reportf(def.Position(), "definition out of order")
		} else {
			minOrder = currOrder
		}
	}
	return nil
}
