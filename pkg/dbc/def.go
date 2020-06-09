package dbc

import (
	"strconv"
	"text/scanner"
)

// Def represents a single definition within a DBC file.
type Def interface {
	// Position of the definition.
	Position() scanner.Position

	// parseFrom parses the definition from a parser.
	parseFrom(*Parser)
}

// VersionDef defines the version of a DBC file.
type VersionDef struct {
	Pos     scanner.Position
	Version string
}

var _ Def = &VersionDef{}

func (d *VersionDef) parseFrom(p *Parser) {
	d.Pos = p.keyword(KeywordVersion).pos
	d.Version = p.string()
}

// Position returns the position of the definition.
func (d *VersionDef) Position() scanner.Position {
	return d.Pos
}

// NewSymbolsDef defines new symbol entries in a DBC file.
type NewSymbolsDef struct {
	Pos     scanner.Position
	Symbols []Keyword
}

var _ Def = &NewSymbolsDef{}

func (d *NewSymbolsDef) parseFrom(p *Parser) {
	p.useWhitespace(significantTab)
	defer p.useWhitespace(defaultWhitespace)
	d.Pos = p.keyword(KeywordNewSymbols).pos
	p.token(':')
	for p.peekToken().typ == '\t' {
		p.token('\t')
		d.Symbols = append(d.Symbols, Keyword(p.identifier()))
	}
}

// Position returns the position of the definition.
func (d *NewSymbolsDef) Position() scanner.Position {
	return d.Pos
}

// BitTimingDef defines the baud rate and the settings of the BTR registers of a CAN network.
//
// This definition is obsolete and not used anymore.
type BitTimingDef struct {
	Pos      scanner.Position
	BaudRate uint64
	BTR1     uint64
	BTR2     uint64
}

var _ Def = &BitTimingDef{}

func (d *BitTimingDef) parseFrom(p *Parser) {
	d.Pos = p.keyword(KeywordBitTiming).pos
	p.token(':')
	d.BaudRate = p.optionalUint()
	if p.peekToken().typ == ':' {
		d.BTR1 = p.optionalUint()
	}
	if p.peekToken().typ == ',' {
		d.BTR2 = p.optionalUint()
	}
}

// Position returns the position of the definition.
func (d *BitTimingDef) Position() scanner.Position {
	return d.Pos
}

// NodesDef defines the names of all nodes participating in the network.
//
// This definition is required in every DBC file.
//
// All node names must be unique.
type NodesDef struct {
	Pos       scanner.Position
	NodeNames []Identifier
}

var _ Def = &NodesDef{}

func (d *NodesDef) parseFrom(p *Parser) {
	p.useWhitespace(significantNewline)
	defer p.useWhitespace(defaultWhitespace)
	d.Pos = p.keyword(KeywordNodes).pos
	p.token(':')
	for p.peekToken().typ == scanner.Ident {
		d.NodeNames = append(d.NodeNames, p.identifier())
	}
	if p.peekToken().typ != scanner.EOF {
		p.token('\n')
	}
}

// Position returns the position of the definition.
func (d *NodesDef) Position() scanner.Position {
	return d.Pos
}

// ValueDescriptionDef defines a textual description for a single signal value.
//
// The value may either be a signal raw value transferred on the bus or the value of an environment variable in a
// remaining bus simulation.
type ValueDescriptionDef struct {
	Pos         scanner.Position
	Value       float64
	Description string
}

var _ Def = &ValueDescriptionDef{}

func (d *ValueDescriptionDef) parseFrom(p *Parser) {
	d.Pos = p.peekToken().pos
	d.Value = p.float()
	d.Description = p.string()
}

// Position returns the position of the definition.
func (d *ValueDescriptionDef) Position() scanner.Position {
	return d.Pos
}

// ValueTableDef defines a global value table.
//
// The value descriptions in value tables define value encodings for signal raw values.
//
// In commonly used DBC files, the global value tables aren't used, but the value descriptions are defined for each
// signal independently.
type ValueTableDef struct {
	Pos               scanner.Position
	TableName         Identifier
	ValueDescriptions []ValueDescriptionDef
}

var _ Def = &ValueTableDef{}

func (d *ValueTableDef) parseFrom(p *Parser) {
	d.Pos = p.keyword(KeywordValueTable).pos
	d.TableName = p.identifier()
	for p.peekToken().typ != ';' {
		valueDescriptionDef := ValueDescriptionDef{}
		valueDescriptionDef.parseFrom(p)
		d.ValueDescriptions = append(d.ValueDescriptions, valueDescriptionDef)
	}
	p.token(';')
}

// Position returns the position of the definition.
func (d *ValueTableDef) Position() scanner.Position {
	return d.Pos
}

// MessageDef defines a frame in the network.
//
// The definition includes the name of a frame as well as its properties and the signals transferred.
type MessageDef struct {
	// Pos is the position of the message definition.
	Pos scanner.Position

	// MessageID contains the message CAN ID.
	//
	// The CAN ID has to be unique within the DBC file.
	//
	// If the most significant bit of the message ID is set, the ID is an extended CAN ID. The extended CAN ID can be
	// determined by masking out the most significant bit with the mask 0xCFFFFFFF.
	MessageID MessageID

	// Name is the name of the message.
	//
	// The message name has to be unique within the DBC file.
	Name Identifier

	// Size specifies the size of the message in bytes.
	Size uint64

	// Transmitter specifies the name of the node transmitting the message.
	//
	// The transmitter has to be defined in the set of node names in the nodes definition.
	//
	// If the message has no transmitter, the string 'Vector__XXX' has to be given here.
	Transmitter Identifier

	// Signals specifies the signals of the message.
	Signals []SignalDef
}

var _ Def = &MessageDef{}

func (d *MessageDef) parseFrom(p *Parser) {
	d.Pos = p.keyword(KeywordMessage).pos
	d.MessageID = p.messageID()
	d.Name = p.identifier()
	p.token(':')
	d.Size = p.uint()
	d.Transmitter = p.identifier()
	for p.peekToken().typ != scanner.EOF && p.peekKeyword() == KeywordSignal {
		signalDef := SignalDef{}
		signalDef.parseFrom(p)
		d.Signals = append(d.Signals, signalDef)
	}
}

// Position returns the position of the definition.
func (d *MessageDef) Position() scanner.Position {
	return d.Pos
}

// SignalDef defines a signal within a message.
type SignalDef struct {
	// Pos is the position of the definition.
	Pos scanner.Position

	// Name of the signal.
	//
	// Has to be unique for all signals within the same message.
	Name Identifier

	// StartBit specifies the position of the signal within the data field of the frame.
	//
	// For signals with byte order Intel (little-endian) the position of the least-significant bit is given.
	//
	// For signals with byte order Motorola (big-endian) the position of the most significant bit is given.
	//
	// The bits are counted in a saw-tooth manner.
	//
	// The start bit has to be in the range of [0 ,8*message_size-1].
	StartBit uint64

	// Size specifies the size of the signal in bits.
	Size uint64

	// IsBigEndian is true if the signal's byte order is Motorola (big-endian).
	IsBigEndian bool

	// IsSigned is true if the signal is signed.
	IsSigned bool

	// IsMultiplexerSwitch is true if the signal is a multiplexer switch.
	//
	// A multiplexer indicator of 'M' defines the signal as the multiplexer switch.
	// Only one signal within a single message can be the multiplexer switch.
	IsMultiplexerSwitch bool

	// IsMultiplexed is true if the signal is multiplexed by the message's multiplexer switch.
	IsMultiplexed bool

	// MultiplexerSwitch is the multiplexer switch value of the signal.
	//
	// The multiplexed signal is transferred in the message if the switch value of the multiplexer signal is equal to
	// its multiplexer switch value.
	MultiplexerSwitch uint64

	// Offset is the signals physical value offset.
	//
	// Together with the factor, the offset defines the linear conversion rule to convert the signal's raw value into
	// the signal's physical value and vice versa.
	//
	//  physical_value = raw_value * factor + offset
	//  raw_value      = (physical_value - offset) / factor
	Offset float64

	// Factor is the signal's physical value factor.
	//
	// See: Offset.
	Factor float64

	// Minimum defines the signal's minimum physical value.
	Minimum float64

	// Maximum defines the signal's maximum physical value.
	Maximum float64

	// Unit defines the unit of the signal's physical value.
	Unit string

	// Receivers specifies the nodes receiving the signal.
	//
	// If the signal has no receiver, the string 'Vector__XXX' has to be given here.
	Receivers []Identifier
}

var _ Def = &SignalDef{}

func (d *SignalDef) parseFrom(p *Parser) {
	d.Pos = p.keyword(KeywordSignal).pos
	d.Name = p.identifier()
	// Parse: Multiplexing
	if p.peekToken().typ != ':' {
		tok := p.nextToken()
		if tok.typ != scanner.Ident {
			p.failf(tok.pos, "expected ident")
		}
		switch {
		case tok.txt == "M":
			d.IsMultiplexerSwitch = true
		case tok.txt[0] == 'm' && len(tok.txt) > 1:
			d.IsMultiplexed = true
			i, err := strconv.Atoi(tok.txt[1:])
			if err != nil || i < 0 {
				p.failf(tok.pos, "invalid multiplexer value")
			}
			d.MultiplexerSwitch = uint64(i)
		default:
			p.failf(tok.pos, "expected multiplexer")
		}
	}
	p.token(':')
	d.StartBit = p.uint()
	p.token('|')
	d.Size = p.uint()
	p.token('@')
	d.IsBigEndian = p.intInRange(0, 1) == 0
	d.IsSigned = p.anyOf('-', '+') == '-'
	p.token('(')
	d.Factor = p.float()
	p.token(',')
	d.Offset = p.float()
	p.token(')')
	p.token('[')
	d.Minimum = p.float()
	p.token('|')
	d.Maximum = p.float()
	p.token(']')
	d.Unit = p.string()
	// Parse: Receivers
	d.Receivers = append(d.Receivers, p.identifier())
	for p.peekToken().typ == ',' {
		p.token(',')
		d.Receivers = append(d.Receivers, p.identifier())
	}
}

// Position returns the position of the definition.
func (d *SignalDef) Position() scanner.Position {
	return d.Pos
}

// SignalValueTypeDef defines an extended type definition for a signal.
type SignalValueTypeDef struct {
	Pos             scanner.Position
	MessageID       MessageID
	SignalName      Identifier
	SignalValueType SignalValueType
}

var _ Def = &SignalValueTypeDef{}

func (d *SignalValueTypeDef) parseFrom(p *Parser) {
	d.Pos = p.keyword(KeywordSignalValueType).pos
	d.MessageID = p.messageID()
	d.SignalName = p.identifier()
	p.optionalToken(':') // SPECIAL-CASE: colon not part of spec, but encountered in the wild
	d.SignalValueType = p.signalValueType()
	p.token(';')
}

// Position returns the position of the definition.
func (d *SignalValueTypeDef) Position() scanner.Position {
	return d.Pos
}

// MessageTransmittersDef defines multiple transmitter nodes of a single message.
//
// This definition is used to describe communication data for higher layer protocols.
//
// This is not used to define CAN layer-2 communication.
type MessageTransmittersDef struct {
	Pos          scanner.Position
	MessageID    MessageID
	Transmitters []Identifier
}

var _ Def = &MessageTransmittersDef{}

func (d *MessageTransmittersDef) parseFrom(p *Parser) {
	d.Pos = p.keyword(KeywordMessageTransmitters).pos
	d.MessageID = p.messageID()
	p.token(':')
	for p.peekToken().typ != ';' {
		d.Transmitters = append(d.Transmitters, p.identifier())
		// SPECIAL-CASE: Comma not included in spec, but encountered in the wild
		p.optionalToken(',')
	}
	p.token(';')
}

// Position returns the position of the definition.
func (d *MessageTransmittersDef) Position() scanner.Position {
	return d.Pos
}

// ValueDescriptionsDef defines inline descriptions for specific raw signal values.
type ValueDescriptionsDef struct {
	Pos                     scanner.Position
	ObjectType              ObjectType
	MessageID               MessageID
	SignalName              Identifier
	EnvironmentVariableName Identifier
	ValueDescriptions       []ValueDescriptionDef
}

var _ Def = &ValueDescriptionsDef{}

func (d *ValueDescriptionsDef) parseFrom(p *Parser) {
	d.Pos = p.keyword(KeywordValueDescriptions).pos
	if p.peekToken().typ == scanner.Ident {
		d.ObjectType = ObjectTypeEnvironmentVariable
		d.EnvironmentVariableName = p.identifier()
	} else {
		d.ObjectType = ObjectTypeSignal
		d.MessageID = p.messageID()
		d.SignalName = p.identifier()
	}
	for p.peekToken().typ != ';' {
		valueDescriptionDef := ValueDescriptionDef{}
		valueDescriptionDef.parseFrom(p)
		d.ValueDescriptions = append(d.ValueDescriptions, valueDescriptionDef)
	}
	p.token(';')
}

// Position returns the position of the definition.
func (d *ValueDescriptionsDef) Position() scanner.Position {
	return d.Pos
}

// EnvironmentVariableDef defines an environment variable.
//
// DBC files that describe the CAN communication and don't define any additional data for system or remaining bus
// simulations don't include environment variables.
type EnvironmentVariableDef struct {
	Pos          scanner.Position
	Name         Identifier
	Type         EnvironmentVariableType
	Minimum      float64
	Maximum      float64
	Unit         string
	InitialValue float64
	ID           uint64
	AccessType   AccessType
	AccessNodes  []Identifier
}

var _ Def = &EnvironmentVariableDef{}

func (d *EnvironmentVariableDef) parseFrom(p *Parser) {
	d.Pos = p.keyword(KeywordEnvironmentVariable).pos
	d.Name = p.identifier()
	p.token(':')
	d.Type = p.environmentVariableType()
	p.token('[')
	d.Minimum = p.float()
	p.token('|')
	d.Maximum = p.float()
	p.token(']')
	d.Unit = p.string()
	d.InitialValue = p.float()
	d.ID = p.uint()
	d.AccessType = p.accessType()
	d.AccessNodes = append(d.AccessNodes, p.identifier())
	for p.peekToken().typ == ',' {
		p.token(',')
		d.AccessNodes = append(d.AccessNodes, p.identifier())
	}
	p.token(';')
}

// Position returns the position of the definition.
func (d *EnvironmentVariableDef) Position() scanner.Position {
	return d.Pos
}

// EnvironmentVariableDataDef defines an environment variable as being of type "data".
//
// Environment variables of this type can store an arbitrary binary data of the given length.
// The length is given in bytes.
type EnvironmentVariableDataDef struct {
	Pos scanner.Position
	// EnvironmentVariableName is the name of the environment variable.
	EnvironmentVariableName Identifier
	// DataSize is the size of the environment variable data in bytes.
	DataSize uint64
}

var _ Def = &EnvironmentVariableDataDef{}

func (d *EnvironmentVariableDataDef) parseFrom(p *Parser) {
	d.Pos = p.keyword(KeywordEnvironmentVariableData).pos
	d.EnvironmentVariableName = p.identifier()
	p.token(':')
	d.DataSize = p.uint()
	p.token(';')
}

// Position returns the position of the definition.
func (d *EnvironmentVariableDataDef) Position() scanner.Position {
	return d.Pos
}

// CommentDef defines a comment.
type CommentDef struct {
	Pos                     scanner.Position
	ObjectType              ObjectType
	NodeName                Identifier
	MessageID               MessageID
	SignalName              Identifier
	EnvironmentVariableName Identifier
	Comment                 string
}

var _ Def = &CommentDef{}

func (d *CommentDef) parseFrom(p *Parser) {
	d.Pos = p.keyword(KeywordComment).pos
	d.ObjectType = p.optionalObjectType()
	switch d.ObjectType {
	case ObjectTypeNetworkNode:
		d.NodeName = p.identifier()
	case ObjectTypeMessage:
		d.MessageID = p.messageID()
	case ObjectTypeSignal:
		d.MessageID = p.messageID()
		d.SignalName = p.identifier()
	case ObjectTypeEnvironmentVariable:
		d.EnvironmentVariableName = p.identifier()
	}
	d.Comment = p.string()
	p.token(';')
}

// Position returns the position of the definition.
func (d *CommentDef) Position() scanner.Position {
	return d.Pos
}

// AttributeDef defines a user-defined attribute.
//
// User-defined attributes are a means to extend the object properties of the DBC file.
//
// These additional attributes have to be defined using an attribute definition with an attribute default value.
//
// For each object having a value defined for the attribute, an attribute value entry has to be defined.
//
// If no attribute value entry is defined for an object, the value of the object's attribute is the attribute's default.
type AttributeDef struct {
	Pos          scanner.Position
	ObjectType   ObjectType
	Name         Identifier
	Type         AttributeValueType
	MinimumInt   int64
	MaximumInt   int64
	MinimumFloat float64
	MaximumFloat float64
	EnumValues   []string
}

var _ Def = &AttributeDef{}

func (d *AttributeDef) parseFrom(p *Parser) {
	d.Pos = p.keyword(KeywordAttribute).pos
	d.ObjectType = p.optionalObjectType()
	d.Name = p.stringIdentifier()
	d.Type = p.attributeValueType()
	switch d.Type {
	case AttributeValueTypeInt, AttributeValueTypeHex:
		if p.peekToken().typ != ';' {
			d.MinimumInt = p.int()
			d.MaximumInt = p.int()
		}
	case AttributeValueTypeFloat:
		if p.peekToken().typ != ';' {
			// SPECIAL CASE: Support attributes without min/max
			d.MinimumFloat = p.float()
			d.MaximumFloat = p.float()
		}
	case AttributeValueTypeEnum:
		d.EnumValues = append(d.EnumValues, p.string())
		for p.peekToken().typ == ',' {
			p.token(',')
			d.EnumValues = append(d.EnumValues, p.string())
		}
	}
	p.token(';')
}

// Position returns the position of the definition.
func (d *AttributeDef) Position() scanner.Position {
	return d.Pos
}

// AttributeDefaultValueDef defines the default value for an attribute.
type AttributeDefaultValueDef struct {
	Pos                scanner.Position
	AttributeName      Identifier
	DefaultIntValue    int64
	DefaultFloatValue  float64
	DefaultStringValue string
}

var _ Def = &AttributeDefaultValueDef{}

func (d *AttributeDefaultValueDef) parseFrom(p *Parser) {
	d.Pos = p.keyword(KeywordAttributeDefault).pos
	d.AttributeName = Identifier(p.string())
	// look up attribute type
	for _, prevDef := range p.defs {
		if attributeDef, ok := prevDef.(*AttributeDef); ok && attributeDef.Name == d.AttributeName {
			switch attributeDef.Type {
			case AttributeValueTypeInt, AttributeValueTypeHex:
				d.DefaultIntValue = p.int()
			case AttributeValueTypeFloat:
				d.DefaultFloatValue = p.float()
			case AttributeValueTypeString:
				d.DefaultStringValue = p.string()
			case AttributeValueTypeEnum:
				d.DefaultStringValue = p.enumValue(attributeDef.EnumValues)
			}
			break
		}
	}
	p.token(';')
}

// Position returns the position of the definition.
func (d *AttributeDefaultValueDef) Position() scanner.Position {
	return d.Pos
}

// AttributeValueForObjectDef defines a value for an attribute and an object.
type AttributeValueForObjectDef struct {
	Pos                     scanner.Position
	AttributeName           Identifier
	ObjectType              ObjectType
	MessageID               MessageID
	SignalName              Identifier
	NodeName                Identifier
	EnvironmentVariableName Identifier
	IntValue                int64
	FloatValue              float64
	StringValue             string
}

var _ Def = &AttributeValueForObjectDef{}

func (d *AttributeValueForObjectDef) parseFrom(p *Parser) {
	d.Pos = p.keyword(KeywordAttributeValue).pos
	d.AttributeName = Identifier(p.string())
	d.ObjectType = p.optionalObjectType()
	switch d.ObjectType {
	case ObjectTypeMessage:
		d.MessageID = p.messageID()
	case ObjectTypeSignal:
		d.MessageID = p.messageID()
		d.SignalName = p.identifier()
	case ObjectTypeNetworkNode:
		d.NodeName = p.identifier()
	case ObjectTypeEnvironmentVariable:
		d.EnvironmentVariableName = p.identifier()
	}
	// look up attribute type
	for _, prevDef := range p.defs {
		if attributeDef, ok := prevDef.(*AttributeDef); ok && attributeDef.Name == d.AttributeName {
			switch attributeDef.Type {
			case AttributeValueTypeInt, AttributeValueTypeHex:
				d.IntValue = p.int()
			case AttributeValueTypeFloat:
				d.FloatValue = p.float()
			case AttributeValueTypeString:
				d.StringValue = p.string()
			case AttributeValueTypeEnum:
				d.StringValue = p.enumValue(attributeDef.EnumValues)
			}
			break
		}
	}
	p.token(';')
}

// Position returns the position of the definition.
func (d *AttributeValueForObjectDef) Position() scanner.Position {
	return d.Pos
}

// UnknownDef represents an unknown or unsupported DBC definition.
type UnknownDef struct {
	Pos     scanner.Position
	Keyword Keyword
}

var _ Def = &UnknownDef{}

func (d *UnknownDef) parseFrom(p *Parser) {
	tok := p.peekToken()
	d.Pos = tok.pos
	d.Keyword = Keyword(tok.txt)
	p.discardLine()
}

// Position returns the position of the definition.
func (d *UnknownDef) Position() scanner.Position {
	return d.Pos
}
