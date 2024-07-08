package dbc

import (
	"os"
	"strings"
	"testing"
	"text/scanner"

	"github.com/davecgh/go-spew/spew"
	"gotest.tools/v3/assert"
)

func shouldUpdateGoldenFiles() bool {
	return os.Getenv("GOLDEN") == "true"
}

func TestParse_ExampleDBC(t *testing.T) {
	const inputFile = "../../testdata/dbc/example/example.dbc"
	const goldenFile = "../../testdata/dbc/example/example.dbc.golden"
	data, err := os.ReadFile(inputFile)
	assert.NilError(t, err)
	p := NewParser(inputFile, data)
	assert.NilError(t, p.Parse())
	if shouldUpdateGoldenFiles() {
		assert.NilError(t, os.WriteFile(goldenFile, []byte(dump(p.Defs())), 0o600))
	}
	goldenFileData, err := os.ReadFile(goldenFile)
	assert.NilError(t, err)
	assert.Equal(t, string(goldenFileData), dump(p.Defs()))
}

func TestParser_Parse(t *testing.T) {
	for _, tt := range []struct {
		name string
		text string
		defs []Def
	}{
		{
			name: "version.dbc",
			text: `VERSION "foo"`,
			defs: []Def{
				&VersionDef{
					Pos: scanner.Position{
						Filename: "version.dbc",
						Line:     1,
						Column:   1,
					},
					Version: "foo",
				},
			},
		},

		{
			name: "multiple_version.dbc",
			text: strings.Join([]string{
				`VERSION "foo"`,
				`VERSION "bar"`,
			}, "\n"),
			defs: []Def{
				&VersionDef{
					Pos: scanner.Position{
						Filename: "multiple_version.dbc",
						Line:     1,
						Column:   1,
					},
					Version: "foo",
				},
				&VersionDef{
					Pos: scanner.Position{
						Filename: "multiple_version.dbc",
						Line:     2,
						Column:   1,
						Offset:   14,
					},
					Version: "bar",
				},
			},
		},

		{
			name: "no_bus_speed.dbc",
			text: `BS_:`,
			defs: []Def{
				&BitTimingDef{
					Pos: scanner.Position{
						Filename: "no_bus_speed.dbc",
						Line:     1,
						Column:   1,
					},
				},
			},
		},

		{
			name: "bus_speed.dbc",
			text: `BS_: 250000`,
			defs: []Def{
				&BitTimingDef{
					Pos: scanner.Position{
						Filename: "bus_speed.dbc",
						Line:     1,
						Column:   1,
					},
					BaudRate: 250000,
				},
			},
		},

		{
			name: "symbols.dbc",
			text: strings.Join([]string{
				"NS_ :",
				"NS_DESC_",
				"CM_",
				"BA_DEF_",
				"BA_",
				"VAL_",
			}, "\n\t"),
			defs: []Def{
				&NewSymbolsDef{
					Pos: scanner.Position{
						Filename: "symbols.dbc",
						Line:     1,
						Column:   1,
					},
					Symbols: []Keyword{
						"NS_DESC_",
						"CM_",
						"BA_DEF_",
						"BA_",
						"VAL_",
					},
				},
			},
		},

		{
			name: "standard_message.dbc",
			text: "BO_ 804 CRUISE: 8 PCM",
			defs: []Def{
				&MessageDef{
					Pos: scanner.Position{
						Filename: "standard_message.dbc",
						Line:     1,
						Column:   1,
					},
					Name:        "CRUISE",
					MessageID:   804,
					Size:        8,
					Transmitter: "PCM",
				},
			},
		},

		{
			name: "extended_message.dbc",
			text: "BO_ 2566857412 BMS2_4: 8 Vector__XXX",
			defs: []Def{
				&MessageDef{
					Pos: scanner.Position{
						Filename: "extended_message.dbc",
						Line:     1,
						Column:   1,
					},
					Name:        "BMS2_4",
					MessageID:   2566857412,
					Size:        8,
					Transmitter: "Vector__XXX",
				},
			},
		},

		{
			name: "signal.dbc",
			text: `SG_ CellTempLowest : 32|8@0+ (1,-40) [-40|215] "C" Vector__XXX`,
			defs: []Def{
				&SignalDef{
					Pos: scanner.Position{
						Filename: "signal.dbc",
						Line:     1,
						Column:   1,
					},
					Name:        "CellTempLowest",
					StartBit:    32,
					Size:        8,
					IsBigEndian: true,
					Factor:      1,
					Offset:      -40,
					Minimum:     -40,
					Maximum:     215,
					Unit:        "C",
					Receivers:   []Identifier{"Vector__XXX"},
				},
			},
		},

		{
			name: "multiplexer_signal.dbc",
			text: `SG_ TestSignal M : 56|8@1+ (0.001,0) [0|0.255] "l/mm" XXX`,
			defs: []Def{
				&SignalDef{
					Pos: scanner.Position{
						Filename: "multiplexer_signal.dbc",
						Line:     1,
						Column:   1,
					},
					Name:                "TestSignal",
					StartBit:            56,
					Size:                8,
					Factor:              0.001,
					Offset:              0,
					Minimum:             0,
					Maximum:             0.255,
					Unit:                "l/mm",
					Receivers:           []Identifier{"XXX"},
					IsMultiplexerSwitch: true,
				},
			},
		},

		{
			name: "multiplexed_signal.dbc",
			text: `SG_ TestSignal m2 : 56|8@1+ (0.001,0) [0|0.255] "l/mm" XXX`,
			defs: []Def{
				&SignalDef{
					Pos: scanner.Position{
						Filename: "multiplexed_signal.dbc",
						Line:     1,
						Column:   1,
					},
					Name:              "TestSignal",
					StartBit:          56,
					Size:              8,
					Factor:            0.001,
					Offset:            0,
					Minimum:           0,
					Maximum:           0.255,
					Unit:              "l/mm",
					Receivers:         []Identifier{"XXX"},
					IsMultiplexed:     true,
					MultiplexerSwitch: 2,
				},
			},
		},

		{
			name: "CSS-Electronics-OBD2-v1.4.dbc",
			text: `SG_ ParameterID_Service02 m2M : 23|8@0+ (1,0) [0|255] "" Vector__XXX`,
			defs: []Def{
				&SignalDef{
					Pos: scanner.Position{
						Filename: "CSS-Electronics-OBD2-v1.4.dbc",
						Line:     1,
						Column:   1,
					},
					Name:                "ParameterID_Service02",
					IsBigEndian:         true,
					StartBit:            23,
					Size:                8,
					Factor:              1,
					Offset:              0,
					Minimum:             0,
					Maximum:             255,
					Unit:                "",
					Receivers:           []Identifier{"Vector__XXX"},
					IsMultiplexerSwitch: true,
					IsMultiplexed:       true,
					MultiplexerSwitch:   2,
				},
			},
		},

		{
			name: "CSS-Electronics-OBD2-v1.4.dbc",
			text: `SG_MUL_VAL_ 2024 S1_PID_5B_HybrBatPackRemLife ParameterID_Service01 91-91;`,
			defs: []Def{
				&SignalMultiplexValueDef{
					Pos: scanner.Position{
						Filename: "CSS-Electronics-OBD2-v1.4.dbc",
						Line:     1,
						Column:   1,
					},
					MessageID:         2024,
					Signal:            "S1_PID_5B_HybrBatPackRemLife",
					MultiplexerSwitch: "ParameterID_Service01",
					RangeStart:        91,
					RangeEnd:          91,
				},
			},
		},

		{
			name: "comment.dbc",
			text: `CM_ "comment";`,
			defs: []Def{
				&CommentDef{
					Pos: scanner.Position{
						Filename: "comment.dbc",
						Line:     1,
						Column:   1,
					},
					Comment: "comment",
				},
			},
		},

		{
			name: "node_comment.dbc",
			text: `CM_ BU_ NodeName "node comment";`,
			defs: []Def{
				&CommentDef{
					Pos: scanner.Position{
						Filename: "node_comment.dbc",
						Line:     1,
						Column:   1,
					},
					ObjectType: ObjectTypeNetworkNode,
					NodeName:   "NodeName",
					Comment:    "node comment",
				},
			},
		},

		{
			name: "message_comment.dbc",
			text: `CM_ BO_ 1234 "message comment";`,
			defs: []Def{
				&CommentDef{
					Pos: scanner.Position{
						Filename: "message_comment.dbc",
						Line:     1,
						Column:   1,
					},
					ObjectType: ObjectTypeMessage,
					MessageID:  1234,
					Comment:    "message comment",
				},
			},
		},

		{
			name: "signal_comment.dbc",
			text: `CM_ SG_ 1234 SignalName "signal comment";`,
			defs: []Def{
				&CommentDef{
					Pos: scanner.Position{
						Filename: "signal_comment.dbc",
						Line:     1,
						Column:   1,
					},
					ObjectType: ObjectTypeSignal,
					MessageID:  1234,
					SignalName: "SignalName",
					Comment:    "signal comment",
				},
			},
		},

		{
			name: "int_attribute_definition.dbc",
			text: `BA_DEF_ "AttributeName" INT 5 10;`,
			defs: []Def{
				&AttributeDef{
					Pos: scanner.Position{
						Filename: "int_attribute_definition.dbc",
						Line:     1,
						Column:   1,
					},
					Name:       "AttributeName",
					Type:       AttributeValueTypeInt,
					MinimumInt: 5,
					MaximumInt: 10,
				},
			},
		},

		{
			name: "int_attribute_definition_no_min_or_max.dbc",
			text: `BA_DEF_ "AttributeName" INT;`,
			defs: []Def{
				&AttributeDef{
					Pos: scanner.Position{
						Filename: "int_attribute_definition_no_min_or_max.dbc",
						Line:     1,
						Column:   1,
					},
					Name: "AttributeName",
					Type: AttributeValueTypeInt,
				},
			},
		},

		{
			name: "float_attribute_definition.dbc",
			text: `BA_DEF_ "AttributeName" FLOAT 0.5 1.5;`,
			defs: []Def{
				&AttributeDef{
					Pos: scanner.Position{
						Filename: "float_attribute_definition.dbc",
						Line:     1,
						Column:   1,
					},
					Name:         "AttributeName",
					Type:         AttributeValueTypeFloat,
					MinimumFloat: 0.5,
					MaximumFloat: 1.5,
				},
			},
		},

		{
			name: "float_attribute_definition_no_min_or_max.dbc",
			text: `BA_DEF_ "AttributeName" FLOAT;`,
			defs: []Def{
				&AttributeDef{
					Pos: scanner.Position{
						Filename: "float_attribute_definition_no_min_or_max.dbc",
						Line:     1,
						Column:   1,
					},
					Name: "AttributeName",
					Type: AttributeValueTypeFloat,
				},
			},
		},

		{
			name: "string_attribute.dbc",
			text: `BA_DEF_ "AttributeName" STRING;`,
			defs: []Def{
				&AttributeDef{
					Pos: scanner.Position{
						Filename: "string_attribute.dbc",
						Line:     1,
						Column:   1,
					},
					Name: "AttributeName",
					Type: AttributeValueTypeString,
				},
			},
		},

		{
			name: "enum_attribute.dbc",
			text: `BA_DEF_ "AttributeName" ENUM "value1","value2";`,
			defs: []Def{
				&AttributeDef{
					Pos: scanner.Position{
						Filename: "enum_attribute.dbc",
						Line:     1,
						Column:   1,
					},
					Name:       "AttributeName",
					Type:       AttributeValueTypeEnum,
					EnumValues: []string{"value1", "value2"},
				},
			},
		},

		{
			name: "enum_attribute_for_messages.dbc",
			text: `BA_DEF_ BO_  "VFrameFormat" ENUM  "StandardCAN","ExtendedCAN","reserved","J1939PG";`,
			defs: []Def{
				&AttributeDef{
					Pos: scanner.Position{
						Filename: "enum_attribute_for_messages.dbc",
						Line:     1,
						Column:   1,
					},
					Name:       "VFrameFormat",
					ObjectType: ObjectTypeMessage,
					Type:       AttributeValueTypeEnum,
					EnumValues: []string{"StandardCAN", "ExtendedCAN", "reserved", "J1939PG"},
				},
			},
		},

		{
			name: "attribute_default_string.dbc",
			text: strings.Join([]string{
				`BA_DEF_ "Foo" STRING;`,
				`BA_DEF_DEF_ "Foo" "string value";`,
			}, "\n"),
			defs: []Def{
				&AttributeDef{
					Pos: scanner.Position{
						Filename: "attribute_default_string.dbc",
						Line:     1,
						Column:   1,
					},
					Name: "Foo",
					Type: AttributeValueTypeString,
				},
				&AttributeDefaultValueDef{
					Pos: scanner.Position{
						Filename: "attribute_default_string.dbc",
						Line:     2,
						Column:   1,
						Offset:   22,
					},
					AttributeName:      "Foo",
					DefaultStringValue: "string value",
				},
			},
		},

		{
			name: "attribute_default_int.dbc",
			text: strings.Join([]string{
				`BA_DEF_ "Foo" INT 0 200;`,
				`BA_DEF_DEF_ "Foo" 100;`,
			}, "\n"),
			defs: []Def{
				&AttributeDef{
					Pos: scanner.Position{
						Filename: "attribute_default_int.dbc",
						Line:     1,
						Column:   1,
					},
					Name:       "Foo",
					Type:       AttributeValueTypeInt,
					MinimumInt: 0,
					MaximumInt: 200,
				},
				&AttributeDefaultValueDef{
					Pos: scanner.Position{
						Filename: "attribute_default_int.dbc",
						Line:     2,
						Column:   1,
						Offset:   25,
					},
					AttributeName:   "Foo",
					DefaultIntValue: 100,
				},
			},
		},

		{
			name: "attribute_default_float.dbc",
			text: strings.Join([]string{
				`BA_DEF_ "Foo" FLOAT 0.5 200.5;`,
				`BA_DEF_DEF_ "Foo" 100.5;`,
			}, "\n"),
			defs: []Def{
				&AttributeDef{
					Pos: scanner.Position{
						Filename: "attribute_default_float.dbc",
						Line:     1,
						Column:   1,
					},
					Name:         "Foo",
					Type:         AttributeValueTypeFloat,
					MinimumFloat: 0.5,
					MaximumFloat: 200.5,
				},
				&AttributeDefaultValueDef{
					Pos: scanner.Position{
						Filename: "attribute_default_float.dbc",
						Line:     2,
						Column:   1,
						Offset:   31,
					},
					AttributeName:     "Foo",
					DefaultFloatValue: 100.5,
				},
			},
		},

		{
			name: "attribute_value.dbc",
			text: strings.Join([]string{
				`BA_DEF_ "Foo" FLOAT;`,
				`BA_ "Foo" 100.5;`,
			}, "\n"),
			defs: []Def{
				&AttributeDef{
					Pos: scanner.Position{
						Filename: "attribute_value.dbc",
						Line:     1,
						Column:   1,
					},
					Name: "Foo",
					Type: AttributeValueTypeFloat,
				},
				&AttributeValueForObjectDef{
					Pos: scanner.Position{
						Filename: "attribute_value.dbc",
						Line:     2,
						Column:   1,
						Offset:   21,
					},
					AttributeName: "Foo",
					FloatValue:    100.5,
				},
			},
		},

		{
			name: "negative_attribute_value.dbc",
			text: strings.Join([]string{
				`BA_DEF_ "Foo" INT;`,
				`BA_ "Foo" -100;`,
			}, "\n"),
			defs: []Def{
				&AttributeDef{
					Pos: scanner.Position{
						Filename: "negative_attribute_value.dbc",
						Line:     1,
						Column:   1,
					},
					Name: "Foo",
					Type: AttributeValueTypeInt,
				},
				&AttributeValueForObjectDef{
					Pos: scanner.Position{
						Filename: "negative_attribute_value.dbc",
						Line:     2,
						Column:   1,
						Offset:   19,
					},
					AttributeName: "Foo",
					IntValue:      -100,
				},
			},
		},

		{
			name: "node_attribute_value.dbc",
			text: strings.Join([]string{
				`BA_DEF_ "Foo" INT;`,
				`BA_ "Foo" BU_ TestNode 100;`,
			}, "\n"),
			defs: []Def{
				&AttributeDef{
					Pos: scanner.Position{
						Filename: "node_attribute_value.dbc",
						Line:     1,
						Column:   1,
					},
					Name: "Foo",
					Type: AttributeValueTypeInt,
				},
				&AttributeValueForObjectDef{
					Pos: scanner.Position{
						Filename: "node_attribute_value.dbc",
						Line:     2,
						Column:   1,
						Offset:   19,
					},
					AttributeName: "Foo",
					ObjectType:    ObjectTypeNetworkNode,
					NodeName:      "TestNode",
					IntValue:      100,
				},
			},
		},

		{
			name: "message_attribute_value.dbc",
			text: strings.Join([]string{
				`BA_DEF_ "Foo" STRING;`,
				`BA_ "Foo" BO_ 1234 "string value";`,
			}, "\n"),
			defs: []Def{
				&AttributeDef{
					Pos: scanner.Position{
						Filename: "message_attribute_value.dbc",
						Line:     1,
						Column:   1,
					},
					Name: "Foo",
					Type: AttributeValueTypeString,
				},
				&AttributeValueForObjectDef{
					Pos: scanner.Position{
						Filename: "message_attribute_value.dbc",
						Line:     2,
						Column:   1,
						Offset:   22,
					},
					AttributeName: "Foo",
					ObjectType:    ObjectTypeMessage,
					MessageID:     1234,
					StringValue:   "string value",
				},
			},
		},

		{
			name: "signal_attribute_value.dbc",
			text: strings.Join([]string{
				`BA_DEF_ "Foo" STRING;`,
				`BA_ "Foo" SG_ 1234 SignalName "string value";`,
			}, "\n"),
			defs: []Def{
				&AttributeDef{
					Pos: scanner.Position{
						Filename: "signal_attribute_value.dbc",
						Line:     1,
						Column:   1,
					},
					Name: "Foo",
					Type: AttributeValueTypeString,
				},
				&AttributeValueForObjectDef{
					Pos: scanner.Position{
						Filename: "signal_attribute_value.dbc",
						Line:     2,
						Column:   1,
						Offset:   22,
					},
					AttributeName: "Foo",
					ObjectType:    ObjectTypeSignal,
					MessageID:     1234,
					SignalName:    "SignalName",
					StringValue:   "string value",
				},
			},
		},

		{
			name: "enum_attribute_value_by_index.dbc",
			text: strings.Join([]string{
				`BA_DEF_ BO_  "VFrameFormat" ENUM  "StandardCAN","ExtendedCAN","reserved","J1939PG";`,
				`BA_ "VFrameFormat" BO_ 1234 3;`,
			}, "\n"),
			defs: []Def{
				&AttributeDef{
					Pos: scanner.Position{
						Filename: "enum_attribute_value_by_index.dbc",
						Line:     1,
						Column:   1,
					},
					Name:       "VFrameFormat",
					ObjectType: ObjectTypeMessage,
					Type:       AttributeValueTypeEnum,
					EnumValues: []string{"StandardCAN", "ExtendedCAN", "reserved", "J1939PG"},
				},
				&AttributeValueForObjectDef{
					Pos: scanner.Position{
						Filename: "enum_attribute_value_by_index.dbc",
						Line:     2,
						Column:   1,
						Offset:   84,
					},
					AttributeName: "VFrameFormat",
					ObjectType:    ObjectTypeMessage,
					MessageID:     1234,
					StringValue:   "J1939PG",
				},
			},
		},

		{
			name: "value_descriptions_for_signal.dbc",
			text: `VAL_ 3 StW_AnglSens_Id 2 "MUST" 0 "PSBL" 1 "SELF";`,
			defs: []Def{
				&ValueDescriptionsDef{
					Pos: scanner.Position{
						Filename: "value_descriptions_for_signal.dbc",
						Line:     1,
						Column:   1,
					},
					ObjectType: ObjectTypeSignal,
					MessageID:  3,
					SignalName: "StW_AnglSens_Id",
					ValueDescriptions: []ValueDescriptionDef{
						{
							Pos: scanner.Position{
								Filename: "value_descriptions_for_signal.dbc",
								Line:     1,
								Column:   24,
								Offset:   23,
							},
							Value:       2,
							Description: "MUST",
						},
						{
							Pos: scanner.Position{
								Filename: "value_descriptions_for_signal.dbc",
								Line:     1,
								Column:   33,
								Offset:   32,
							},
							Value:       0,
							Description: "PSBL",
						},
						{
							Pos: scanner.Position{
								Filename: "value_descriptions_for_signal.dbc",
								Line:     1,
								Column:   42,
								Offset:   41,
							},
							Value:       1,
							Description: "SELF",
						},
					},
				},
			},
		},

		{
			name: "value_table.dbc",
			text: `VAL_TABLE_ DI_gear 7 "DI_GEAR_SNA" 4 "DI_GEAR_D";`,
			defs: []Def{
				&ValueTableDef{
					Pos: scanner.Position{
						Filename: "value_table.dbc",
						Line:     1,
						Column:   1,
					},
					TableName: "DI_gear",
					ValueDescriptions: []ValueDescriptionDef{
						{
							Pos: scanner.Position{
								Filename: "value_table.dbc",
								Line:     1,
								Column:   20,
								Offset:   19,
							},
							Value:       7,
							Description: "DI_GEAR_SNA",
						},
						{
							Pos: scanner.Position{
								Filename: "value_table.dbc",
								Line:     1,
								Column:   36,
								Offset:   35,
							},
							Value:       4,
							Description: "DI_GEAR_D",
						},
					},
				},
			},
		},

		{
			name: "node_list.dbc",
			text: `BU_: RSDS`,
			defs: []Def{
				&NodesDef{
					Pos: scanner.Position{
						Filename: "node_list.dbc",
						Line:     1,
						Column:   1,
					},
					NodeNames: []Identifier{"RSDS"},
				},
			},
		},

		{
			name: "node_list_followed_by_single_newline_followed_by_value_table.dbc",
			text: strings.Join([]string{
				`BU_: RSDS`,
				`VAL_TABLE_ TableName 3 "Value3" 2 "Value2" 1 "Value1" 0 "Value0";`,
			}, "\n"),
			defs: []Def{
				&NodesDef{
					Pos: scanner.Position{
						Filename: "node_list_followed_by_single_newline_followed_by_value_table.dbc",
						Line:     1,
						Column:   1,
					},
					NodeNames: []Identifier{"RSDS"},
				},
				&ValueTableDef{
					Pos: scanner.Position{
						Filename: "node_list_followed_by_single_newline_followed_by_value_table.dbc",
						Line:     2,
						Column:   1,
						Offset:   10,
					},
					TableName: "TableName",
					ValueDescriptions: []ValueDescriptionDef{
						{
							Pos: scanner.Position{
								Filename: "node_list_followed_by_single_newline_followed_by_value_table.dbc",
								Line:     2,
								Column:   22,
								Offset:   31,
							},
							Value:       3,
							Description: "Value3",
						},
						{
							Pos: scanner.Position{
								Filename: "node_list_followed_by_single_newline_followed_by_value_table.dbc",
								Line:     2,
								Column:   33,
								Offset:   42,
							},
							Value:       2,
							Description: "Value2",
						},
						{
							Pos: scanner.Position{
								Filename: "node_list_followed_by_single_newline_followed_by_value_table.dbc",
								Line:     2,
								Column:   44,
								Offset:   53,
							},
							Value:       1,
							Description: "Value1",
						},
						{
							Pos: scanner.Position{
								Filename: "node_list_followed_by_single_newline_followed_by_value_table.dbc",
								Line:     2,
								Column:   55,
								Offset:   64,
							},
							Value:       0,
							Description: "Value0",
						},
					},
				},
			},
		},

		{
			name: "signal_value_type.dbc",
			text: `SIG_VALTYPE_ 42 TestSignal 2;`,
			defs: []Def{
				&SignalValueTypeDef{
					Pos: scanner.Position{
						Filename: "signal_value_type.dbc",
						Line:     1,
						Column:   1,
					},
					MessageID:       42,
					SignalName:      "TestSignal",
					SignalValueType: SignalValueTypeFloat64,
				},
			},
		},

		{
			name: "signal_value_type_with_colon.dbc",
			text: `SIG_VALTYPE_ 42 TestSignal : 2;`,
			defs: []Def{
				&SignalValueTypeDef{
					Pos: scanner.Position{
						Filename: "signal_value_type_with_colon.dbc",
						Line:     1,
						Column:   1,
					},
					MessageID:       42,
					SignalName:      "TestSignal",
					SignalValueType: SignalValueTypeFloat64,
				},
			},
		},

		{
			name: "message_transmitters.dbc",
			text: `BO_TX_BU_ 42: Node1 Node2;`,
			defs: []Def{
				&MessageTransmittersDef{
					Pos: scanner.Position{
						Filename: "message_transmitters.dbc",
						Line:     1,
						Column:   1,
					},
					MessageID:    42,
					Transmitters: []Identifier{"Node1", "Node2"},
				},
			},
		},

		{
			name: "message_transmitters_comma_separated.dbc",
			text: `BO_TX_BU_ 42: Node1,Node2;`,
			defs: []Def{
				&MessageTransmittersDef{
					Pos: scanner.Position{
						Filename: "message_transmitters_comma_separated.dbc",
						Line:     1,
						Column:   1,
					},
					MessageID:    42,
					Transmitters: []Identifier{"Node1", "Node2"},
				},
			},
		},

		{
			name: "environment_variable_data.dbc",
			text: `ENVVAR_DATA_ VariableName: 42;`,
			defs: []Def{
				&EnvironmentVariableDataDef{
					Pos: scanner.Position{
						Filename: "environment_variable_data.dbc",
						Line:     1,
						Column:   1,
					},
					EnvironmentVariableName: "VariableName",
					DataSize:                42,
				},
			},
		},

		{
			name: "environment_variable_value_descriptions.dbc",
			text: `VAL_ VariableName 2 "Value2" 1 "Value1" 0 "Value0";`,
			defs: []Def{
				&ValueDescriptionsDef{
					Pos: scanner.Position{
						Filename: "environment_variable_value_descriptions.dbc",
						Line:     1,
						Column:   1,
					},
					ObjectType:              ObjectTypeEnvironmentVariable,
					EnvironmentVariableName: "VariableName",
					ValueDescriptions: []ValueDescriptionDef{
						{
							Pos: scanner.Position{
								Filename: "environment_variable_value_descriptions.dbc",
								Line:     1,
								Column:   19,
								Offset:   18,
							},
							Value:       2,
							Description: "Value2",
						},
						{
							Pos: scanner.Position{
								Filename: "environment_variable_value_descriptions.dbc",
								Line:     1,
								Column:   30,
								Offset:   29,
							},
							Value:       1,
							Description: "Value1",
						},
						{
							Pos: scanner.Position{
								Filename: "environment_variable_value_descriptions.dbc",
								Line:     1,
								Column:   41,
								Offset:   40,
							},
							Value:       0,
							Description: "Value0",
						},
					},
				},
			},
		},

		{
			name: "unknown_def.dbc",
			text: `FOO_ Bar 2 Baz;`,
			defs: []Def{
				&UnknownDef{
					Pos: scanner.Position{
						Filename: "unknown_def.dbc",
						Line:     1,
						Column:   1,
					},
					Keyword: "FOO_",
				},
			},
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.name, []byte(tt.text))
			assert.NilError(t, p.Parse())
			assert.DeepEqual(t, tt.defs, p.Defs())
		})
	}
}

func TestParser_Parse_Error(t *testing.T) {
	for _, tt := range []struct {
		name string
		text string
		err  *parseError
	}{
		{
			name: "non_string_version.dbc",
			text: "VERSION foo",
			err: &parseError{
				pos: scanner.Position{
					Filename: "non_string_version.dbc",
					Line:     1,
					Column:   9,
					Offset:   8,
				},
				reason: "expected token \"",
			},
		},
		{
			name: "invalid_utf8.dbc",
			text: `VERSION "foo` + string([]byte{0xc3, 0x28}) + `"`,
			err: &parseError{
				pos: scanner.Position{
					Filename: "invalid_utf8.dbc",
					Line:     1,
					Column:   13,
					Offset:   12,
				},
				reason: "invalid UTF-8 encoding",
			},
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.name, []byte(tt.text))
			assert.Error(t, p.Parse(), tt.err.Error())
		})
	}
}

func dump(data interface{}) string {
	spewConfig := spew.ConfigState{
		Indent:                  " ",
		DisablePointerAddresses: true,
		DisableCapacities:       true,
		SortKeys:                true,
	}
	return spewConfig.Sdump(data)
}
