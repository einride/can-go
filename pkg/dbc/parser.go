package dbc

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"text/scanner"
	"unicode/utf8"
)

const defaultScannerMode = scanner.ScanIdents | scanner.ScanFloats

const (
	defaultWhitespace  = scanner.GoWhitespace
	significantNewline = defaultWhitespace & ^uint64(1<<'\n')
	significantTab     = defaultWhitespace & ^uint64(1<<'\t')
)

type token struct {
	typ rune
	pos scanner.Position
	txt string
}

type Parser struct {
	sc           scanner.Scanner
	curr         token
	lookahead    token
	hasLookahead bool
	data         []byte
	defs         []Def
}

func NewParser(filename string, data []byte) *Parser {
	p := &Parser{data: data}
	p.sc.Init(bytes.NewReader(data))
	p.sc.Mode = defaultScannerMode
	p.sc.Whitespace = defaultWhitespace
	p.sc.Filename = filename
	p.sc.Error = func(sc *scanner.Scanner, msg string) {
		p.failf(sc.Pos(), msg)
	}
	return p
}

func (p *Parser) Defs() []Def {
	return p.defs
}

func (p *Parser) File() *File {
	return &File{
		Name: p.sc.Filename,
		Data: p.data,
		Defs: p.defs,
	}
}

func (p *Parser) Parse() (err Error) {
	defer func() {
		if r := recover(); r != nil {
			// recover from parse errors only
			if errParse, ok := r.(*parseError); ok {
				err = errParse
			} else {
				panic(r)
			}
		}
	}()
	for p.peekToken().typ != scanner.EOF {
		var def Def
		switch p.peekKeyword() {
		case KeywordVersion:
			def = &VersionDef{}
		case KeywordBitTiming:
			def = &BitTimingDef{}
		case KeywordNewSymbols:
			def = &NewSymbolsDef{}
		case KeywordNodes:
			def = &NodesDef{}
		case KeywordMessage:
			def = &MessageDef{}
		case KeywordSignal:
			def = &SignalDef{}
		case KeywordEnvironmentVariable:
			def = &EnvironmentVariableDef{}
		case KeywordComment:
			def = &CommentDef{}
		case KeywordAttribute:
			def = &AttributeDef{}
		case KeywordAttributeDefault:
			def = &AttributeDefaultValueDef{}
		case KeywordAttributeValue:
			def = &AttributeValueForObjectDef{}
		case KeywordValueDescriptions:
			def = &ValueDescriptionsDef{}
		case KeywordValueTable:
			def = &ValueTableDef{}
		case KeywordSignalValueType:
			def = &SignalValueTypeDef{}
		case KeywordMessageTransmitters:
			def = &MessageTransmittersDef{}
		case KeywordEnvironmentVariableData:
			def = &EnvironmentVariableDataDef{}
		default:
			def = &UnknownDef{}
		}
		def.parseFrom(p)
		p.defs = append(p.defs, def)
	}
	return nil
}

func (p *Parser) failf(pos scanner.Position, format string, a ...interface{}) {
	panic(&parseError{pos: pos, reason: fmt.Sprintf(format, a...)})
}

//
// Whitespace
//

func (p *Parser) useWhitespace(whitespace uint64) {
	p.sc.Whitespace = whitespace
}

//
// Characters
//

func (p *Parser) nextRune() rune {
	if p.hasLookahead {
		if utf8.RuneCountInString(p.lookahead.txt) > 1 {
			p.failf(p.lookahead.pos, "cannot get next rune when lookahead contains a token")
		}
		p.hasLookahead = false
		r, _ := utf8.DecodeRuneInString(p.lookahead.txt)
		return r
	}
	return p.sc.Next()
}

func (p *Parser) peekRune() rune {
	if p.hasLookahead {
		if utf8.RuneCountInString(p.lookahead.txt) > 1 {
			p.failf(p.lookahead.pos, "cannot peek next rune when lookahead contains a token")
		}
		r, _ := utf8.DecodeRuneInString(p.lookahead.txt)
		return r
	}
	return p.sc.Peek()
}

func (p *Parser) discardLine() {
	p.useWhitespace(significantNewline)
	defer p.useWhitespace(defaultWhitespace)
	// skip all non-newline tokens
	for p.nextToken().typ != '\n' && p.nextToken().typ != scanner.EOF {
		_ = p.curr // fool the linter about the empty loop
	}
}

//
// Tokens
//

func (p *Parser) nextToken() token {
	if p.hasLookahead {
		p.hasLookahead = false
		p.curr = p.lookahead
		return p.lookahead
	}
	p.curr = token{typ: p.sc.Scan(), pos: p.sc.Position, txt: p.sc.TokenText()}
	return p.curr
}

func (p *Parser) peekToken() token {
	if p.hasLookahead {
		return p.lookahead
	}
	p.hasLookahead = true
	p.lookahead = token{typ: p.sc.Scan(), pos: p.sc.Position, txt: p.sc.TokenText()}
	return p.lookahead
}

//
// Data types
//

// string parses a string that may contain newlines.
func (p *Parser) string() string {
	tok := p.nextToken()
	if tok.typ != '"' {
		p.failf(tok.pos, `expected token "`)
	}
	var b strings.Builder
ReadLoop:
	for {
		switch r := p.nextRune(); r {
		case scanner.EOF:
			p.failf(tok.pos, "unterminated string")
		case '"':
			break ReadLoop
		case '\n':
			if _, err := b.WriteRune(' '); err != nil {
				p.failf(tok.pos, err.Error())
			}
		case '\\':
			if p.peekRune() == '"' {
				_ = p.nextRune() // include escaped quotes in string
				if _, err := b.WriteString(`\"`); err != nil {
					p.failf(tok.pos, err.Error())
				}
				continue
			}
			fallthrough
		default:
			if _, err := b.WriteRune(r); err != nil {
				p.failf(tok.pos, err.Error())
			}
		}
	}
	return b.String()
}

func (p *Parser) identifier() Identifier {
	tok := p.nextToken()
	if tok.typ != scanner.Ident {
		p.failf(tok.pos, "expected ident")
	}
	id := Identifier(tok.txt)
	if err := id.Validate(); err != nil {
		p.failf(tok.pos, err.Error())
	}
	return id
}

func (p *Parser) stringIdentifier() Identifier {
	tok := p.peekToken()
	id := Identifier(p.string())
	if err := id.Validate(); err != nil {
		p.failf(tok.pos, err.Error())
	}
	return id
}

func (p *Parser) keyword(kw Keyword) token {
	if p.peekKeyword() != kw {
		p.failf(p.peekToken().pos, "expected keyword: %v", kw)
	}
	return p.nextToken()
}

func (p *Parser) peekKeyword() Keyword {
	tok := p.peekToken()
	if tok.typ != scanner.Ident {
		p.failf(p.peekToken().pos, "expected ident")
	}
	return Keyword(tok.txt)
}

func (p *Parser) token(typ rune) {
	if tok := p.nextToken(); tok.typ != typ {
		p.failf(
			p.peekToken().pos,
			"expected token: %v, found: %v (%v)",
			scanner.TokenString(typ),
			scanner.TokenString(tok.typ),
			tok.txt,
		)
	}
}

func (p *Parser) optionalToken(typ rune) {
	if p.peekToken().typ == typ {
		p.token(typ)
	}
}

func (p *Parser) enumValue(values []string) string {
	tok := p.peekToken()
	if tok.typ == scanner.Int {
		// SPECIAL-CASE: Enum values by index encountered in the wild
		i := p.uint()
		if i >= uint64(len(values)) {
			p.failf(tok.pos, "enum index out of bounds")
		}
		return values[i]
	}
	return p.string()
}

func (p *Parser) float() float64 {
	var isNegative bool
	if p.peekToken().typ == '-' {
		p.token('-')
		isNegative = true
	}
	tok := p.nextToken()
	if tok.typ != scanner.Int && tok.typ != scanner.Float {
		p.failf(p.peekToken().pos, "expected int or float")
	}
	f, err := strconv.ParseFloat(tok.txt, 64)
	if err != nil {
		p.failf(tok.pos, "invalid float")
	}
	if isNegative {
		f *= -1
	}
	return f
}

func (p *Parser) int() int64 {
	var isNegative bool
	if p.peekToken().typ == '-' {
		p.token('-')
		isNegative = true
	}
	tok := p.nextToken()
	if tok.typ != scanner.Int && tok.typ != scanner.Float {
		p.failf(tok.pos, "expected int or float")
	}
	f, err := strconv.ParseFloat(tok.txt, 64)
	if err != nil {
		p.failf(tok.pos, "invalid int")
	}
	i := int64(f)
	if f > math.MaxInt64 {
		i = math.MaxInt64
	} else if f < math.MinInt64 {
		i = math.MinInt64
	}
	if isNegative {
		i *= -1
	}
	return i
}

func (p *Parser) uint() uint64 {
	tok := p.nextToken()
	if tok.typ != scanner.Int {
		p.failf(tok.pos, "expected int")
	}
	i, err := strconv.ParseUint(tok.txt, 10, 64)
	if err != nil {
		p.failf(tok.pos, "invalid uint")
	}
	return i
}

func (p *Parser) intInRange(min, max int) int {
	var isNegative bool
	if p.peekToken().typ == '-' {
		p.token('-')
		isNegative = true
	}
	tok := p.nextToken()
	i, err := strconv.Atoi(tok.txt)
	if err != nil {
		p.failf(tok.pos, "invalid int")
	}
	if isNegative {
		i *= -1
	}
	if i < min || i > max {
		p.failf(tok.pos, "invalid value")
	}
	return i
}

func (p *Parser) optionalUint() uint64 {
	if p.peekToken().typ != scanner.Int {
		return 0
	}
	tok := p.nextToken()
	i, err := strconv.ParseUint(tok.txt, 10, 64)
	if err != nil {
		p.failf(tok.pos, "invalid uint")
	}
	return i
}

func (p *Parser) anyOf(tokenTypes ...rune) rune {
	tok := p.nextToken()
	for _, tokenType := range tokenTypes {
		if tok.typ == tokenType {
			return tok.typ
		}
	}
	p.failf(tok.pos, "unexpected token")
	return 0
}

func (p *Parser) optionalObjectType() ObjectType {
	tok := p.peekToken()
	if tok.typ != scanner.Ident {
		return ObjectTypeUnspecified
	}
	objectType := ObjectType(p.identifier())
	if err := objectType.Validate(); err != nil {
		p.failf(tok.pos, err.Error())
	}
	return objectType
}

func (p *Parser) messageID() MessageID {
	tok := p.peekToken()
	messageID := MessageID(p.uint())
	if err := messageID.Validate(); err != nil {
		p.failf(tok.pos, err.Error())
	}
	return messageID
}

func (p *Parser) signalValueType() SignalValueType {
	tok := p.peekToken()
	signalValueType := SignalValueType(p.uint())
	if err := signalValueType.Validate(); err != nil {
		p.failf(tok.pos, err.Error())
	}
	return signalValueType
}

func (p *Parser) environmentVariableType() EnvironmentVariableType {
	tok := p.peekToken()
	environmentVariableType := EnvironmentVariableType(p.uint())
	if err := environmentVariableType.Validate(); err != nil {
		p.failf(tok.pos, err.Error())
	}
	return environmentVariableType
}

func (p *Parser) attributeValueType() AttributeValueType {
	tok := p.peekToken()
	attributeValueType := AttributeValueType(p.identifier())
	if err := attributeValueType.Validate(); err != nil {
		p.failf(tok.pos, err.Error())
	}
	return attributeValueType
}

func (p *Parser) accessType() AccessType {
	tok := p.peekToken()
	accessType := AccessType(p.identifier())
	if err := accessType.Validate(); err != nil {
		p.failf(tok.pos, "invalid access type")
	}
	return accessType
}
