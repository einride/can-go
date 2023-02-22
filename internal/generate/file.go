package generate

import (
	"bytes"
	"fmt"
	"go/format"
	"go/types"
	"path"
	"regexp"
	"strings"

	"github.com/blueinnovationsgroup/can-go/pkg/descriptor"
	"github.com/shurcooL/go-goon"
)

type File struct {
	buf bytes.Buffer
	err error
}

func NewFile() *File {
	f := &File{}
	f.buf.Grow(1e5) // 100K
	return f
}

func (f *File) Write(p []byte) (int, error) {
	if f.err != nil {
		return 0, f.err
	}
	n, err := f.buf.Write(p)
	f.err = err
	return n, err
}

func (f *File) P(v ...interface{}) {
	for _, x := range v {
		_, _ = fmt.Fprint(f, x)
	}
	_, _ = fmt.Fprintln(f)
}

func (f *File) Dump(v interface{}) {
	_, _ = goon.Fdump(f, v)
}

func (f *File) Content() ([]byte, error) {
	if f.err != nil {
		return nil, fmt.Errorf("file content: %w", f.err)
	}
	formatted, err := format.Source(f.buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("file content: %s: %w", f.buf.String(), err)
	}
	return formatted, nil
}

func Database(d *descriptor.Database) ([]byte, error) {
	f := NewFile()
	Package(f, d)
	Imports(f)
	for _, m := range d.Messages {
		MessageType(f, m)
		for _, s := range m.Signals {
			if hasCustomType(s) {
				SignalCustomType(f, m, s)
			}
		}
		MarshalFrame(f, m)
		UnmarshalFrame(f, m)
	}
	if hasSendType(d) { // only code-generate nodes for schemas with send types specified
		for _, n := range d.Nodes {
			Node(f, d, n)
		}
	}
	Descriptors(f, d)
	return f.Content()
}

func Package(f *File, d *descriptor.Database) {
	packageName := strings.TrimSuffix(path.Base(d.SourceFile), path.Ext(d.SourceFile)) + "can"
	// Remove illegal characters from package name
	packageName = strings.ReplaceAll(packageName, ".", "")
	packageName = strings.ReplaceAll(packageName, "-", "")
	packageName = strings.ReplaceAll(packageName, "_", "")
	f.P("// Package ", packageName, " provides primitives for encoding and decoding ", d.Name(), " CAN messages.")
	f.P("//")
	f.P("// Source: ", d.SourceFile)
	f.P("package ", packageName)
	f.P()
}

func Imports(f *File) {
	f.P("import (")
	f.P(`"context"`)
	f.P(`"fmt"`)
	f.P(`"net"`)
	f.P(`"net/http"`)
	f.P(`"sync"`)
	f.P(`"time"`)
	f.P()
	f.P(`"github.com/blueinnovationsgroup/can-go"`)
	f.P(`"github.com/blueinnovationsgroup/can-go/pkg/socketcan"`)
	f.P(`"github.com/blueinnovationsgroup/can-go/pkg/candebug"`)
	f.P(`"github.com/blueinnovationsgroup/can-go/pkg/canrunner"`)
	f.P(`"github.com/blueinnovationsgroup/can-go/pkg/descriptor"`)
	f.P(`"github.com/blueinnovationsgroup/can-go/pkg/generated"`)
	f.P(`"github.com/blueinnovationsgroup/can-go/pkg/cantext"`)
	f.P(")")
	f.P()
	// we could use goimports for this, but it significantly slows down code generation
	f.P("// prevent unused imports")
	f.P("var (")
	f.P("_ = context.Background")
	f.P("_ = fmt.Print")
	f.P("_ = net.Dial")
	f.P("_ = http.Error")
	f.P("_ = sync.Mutex{}")
	f.P("_ = time.Now")
	f.P("_ = socketcan.Dial")
	f.P("_ = candebug.ServeMessagesHTTP")
	f.P("_ = canrunner.Run")
	f.P(")")
	f.P()
	f.P("// Generated code. DO NOT EDIT.")
}

var nonAlphaNumericRegexp = regexp.MustCompile("[^a-zA-Z0-9]+")

func slugifyString(s string) string {
	return nonAlphaNumericRegexp.ReplaceAllString(s, "")
}

func SignalCustomType(f *File, m *descriptor.Message, s *descriptor.Signal) {
	f.P("// ", signalType(m, s), " models the ", s.Name, " signal of the ", m.Name, " message.")
	f.P("type ", signalType(m, s), " ", signalPrimitiveType(s))
	f.P()
	f.P("// Value descriptions for the ", s.Name, " signal of the ", m.Name, " message.")
	f.P("const (")
	for _, vd := range s.ValueDescriptions {
		desc := slugifyString(vd.Description)
		switch {
		case s.Length == 1 && vd.Value == 1:
			f.P(signalType(m, s), "_", desc, " ", signalType(m, s), " = true")
		case s.Length == 1 && vd.Value == 0:
			f.P(signalType(m, s), "_", desc, " ", signalType(m, s), " = false")
		default:
			f.P(signalType(m, s), "_", desc, " ", signalType(m, s), " = ", vd.Value)
		}
	}
	f.P(")")
	f.P()
	f.P("func (v ", signalType(m, s), ") String() string {")
	if s.Length == 1 {
		f.P("switch bool(v) {")
		for _, vd := range s.ValueDescriptions {
			if vd.Value == 1 {
				f.P("case true:")
			} else {
				f.P("case false:")
			}
			f.P(`return "`, vd.Description, `"`)
		}
		f.P("}")
		f.P(`return fmt.Sprintf("`, signalType(m, s), `(%t)", v)`)
	} else {
		f.P("switch v {")
		for _, vd := range s.ValueDescriptions {
			f.P("case ", vd.Value, ":")
			f.P(`return "`, vd.Description, `"`)
		}
		f.P("default:")
		f.P(`return fmt.Sprintf("`, signalType(m, s), `(%d)", v)`)
		f.P("}")
	}
	f.P("}")
}

func MessageType(f *File, m *descriptor.Message) {
	f.P("// ", messageReaderInterface(m), " provides read access to a ", m.Name, " message.")
	f.P("type ", messageReaderInterface(m), " interface {")
	for _, s := range m.Signals {
		if hasPhysicalRepresentation(s) {
			f.P("// ", s.Name, " returns the physical value of the ", s.Name, " signal.")
			f.P(s.Name, "() float64")
			if len(s.ValueDescriptions) > 0 {
				f.P()
				f.P("// ", s.Name, " returns the raw (encoded) value of the ", s.Name, " signal.")
				f.P("Raw", s.Name, "() ", signalType(m, s))
			}
		} else {
			f.P("// ", s.Name, " returns the value of the ", s.Name, " signal.")
			f.P(s.Name, "()", signalType(m, s))
		}
	}
	f.P("}")
	f.P()
	f.P("// ", messageWriterInterface(m), " provides write access to a ", m.Name, " message.")
	f.P("type ", messageWriterInterface(m), " interface {")
	f.P("// CopyFrom copies all values from ", messageReaderInterface(m), ".")
	f.P("CopyFrom(", messageReaderInterface(m), ") *", messageStruct(m))
	for _, s := range m.Signals {
		if hasPhysicalRepresentation(s) {
			f.P("// Set", s.Name, " sets the physical value of the ", s.Name, " signal.")
			f.P("Set", s.Name, "(float64) *", messageStruct(m))
			if len(s.ValueDescriptions) > 0 {
				f.P()
				f.P("// SetRaw", s.Name, " sets the raw (encoded) value of the ", s.Name, " signal.")
				f.P("SetRaw", s.Name, "(", signalType(m, s), ") *", messageStruct(m))
			}
		} else {
			f.P("// Set", s.Name, " sets the value of the ", s.Name, " signal.")
			f.P("Set", s.Name, "(", signalType(m, s), ") *", messageStruct(m))
		}
	}
	f.P("}")
	f.P()
	f.P("type ", messageStruct(m), " struct {")
	for _, s := range m.Signals {
		f.P(signalField(s), " ", signalType(m, s))
	}
	f.P("}")
	f.P()
	f.P("func New", messageStruct(m), "() *", messageStruct(m), " {")
	f.P("m := &", messageStruct(m), "{}")
	f.P("m.Reset()")
	f.P("return m")
	f.P("}")
	f.P()
	f.P("func (m *", messageStruct(m), ") Reset() {")
	for _, s := range m.Signals {
		switch {
		case s.Length == 1 && s.DefaultValue == 1:
			f.P("m.", signalField(s), " = true")
		case s.Length == 1:
			f.P("m.", signalField(s), " = false")
		default:
			f.P("m.", signalField(s), " = ", s.DefaultValue)
		}
	}
	f.P("}")
	f.P()
	f.P("func (m *", messageStruct(m), ") CopyFrom(o ", messageReaderInterface(m), ") *", messageStruct(m), "{")
	for _, s := range m.Signals {
		if hasPhysicalRepresentation(s) {
			f.P("m.Set", s.Name, "(o.", s.Name, "())")
		} else {
			f.P("m.", signalField(s), " = o.", s.Name, "()")
		}
	}
	f.P("return m")
	f.P("}")
	f.P()
	f.P("// Descriptor returns the ", m.Name, " descriptor.")
	f.P("func (m *", messageStruct(m), ") Descriptor() *descriptor.Message {")
	f.P("return ", messageDescriptor(m), ".Message")
	f.P("}")
	f.P()
	f.P("// String returns a compact string representation of the message.")
	f.P("func(m *", messageStruct(m), ") String() string {")
	f.P("return cantext.MessageString(m)")
	f.P("}")
	f.P()
	for _, s := range m.Signals {
		if !hasPhysicalRepresentation(s) {
			f.P("func (m *", messageStruct(m), ") ", s.Name, "() ", signalType(m, s), " {")
			f.P("return m.", signalField(s))
			f.P("}")
			f.P()
			f.P("func (m *", messageStruct(m), ") Set", s.Name, "(v ", signalType(m, s), ") *", messageStruct(m), " {")
			if s.Length == 1 {
				f.P("m.", signalField(s), " = v")
			} else {
				f.P(
					"m.", signalField(s), " = ", signalType(m, s), "(",
					signalDescriptor(m, s), ".SaturatedCast", signalSuperType(s), "(",
					signalPrimitiveSuperType(s), "(v)))",
				)
			}
			f.P("return m")
			f.P("}")
			f.P()
			continue
		}
		f.P("func (m *", messageStruct(m), ") ", s.Name, "() float64 {")
		f.P("return ", signalDescriptor(m, s), ".ToPhysical(float64(m.", signalField(s), "))")
		f.P("}")
		f.P()
		f.P("func (m *", messageStruct(m), ") Set", s.Name, "(v float64) *", messageStruct(m), " {")
		f.P("m.", signalField(s), " = ", signalType(m, s), "(", signalDescriptor(m, s), ".FromPhysical(v))")
		f.P("return m")
		f.P("}")
		f.P()
		if len(s.ValueDescriptions) > 0 {
			f.P("func (m *", messageStruct(m), ") Raw", s.Name, "() ", signalType(m, s), " {")
			f.P("return m.", signalField(s))
			f.P("}")
			f.P()
			f.P("func (m *", messageStruct(m), ") SetRaw", s.Name, "(v ", signalType(m, s), ") *", messageStruct(m), "{")
			f.P(
				"m.", signalField(s), " = ", signalType(m, s), "(",
				signalDescriptor(m, s), ".SaturatedCast", signalSuperType(s), "(",
				signalPrimitiveSuperType(s), "(v)))",
			)
			f.P("return m")
			f.P("}")
			f.P()
		}
	}
}

func Descriptors(f *File, d *descriptor.Database) {
	f.P("// Nodes returns the ", d.Name(), " node descriptors.")
	f.P("func Nodes() *NodesDescriptor {")
	f.P("return nd")
	f.P("}")
	f.P()
	f.P("// NodesDescriptor contains all ", d.Name(), " node descriptors.")
	f.P("type NodesDescriptor struct{")
	for _, n := range d.Nodes {
		f.P(n.Name, " *descriptor.Node")
	}
	f.P("}")
	f.P()
	f.P("// Messages returns the ", d.Name(), " message descriptors.")
	f.P("func Messages() *MessagesDescriptor {")
	f.P("return md")
	f.P("}")
	f.P()
	f.P("// MessagesDescriptor contains all ", d.Name(), " message descriptors.")
	f.P("type MessagesDescriptor struct{")
	for _, m := range d.Messages {
		f.P(m.Name, " *", m.Name, "Descriptor")
	}
	f.P("}")
	f.P()
	f.P("// UnmarshalFrame unmarshals the provided ", d.Name(), " CAN frame.")
	f.P("func (md *MessagesDescriptor) UnmarshalFrame(f can.Frame) (generated.Message, error) {")
	f.P("switch f.ID {")
	for _, m := range d.Messages {
		f.P("case md.", m.Name, ".ID:")
		f.P("var msg ", messageStruct(m))
		f.P("if err := msg.UnmarshalFrame(f); err != nil {")
		f.P(`return nil, fmt.Errorf("unmarshal `, d.Name(), ` frame: %w", err)`)
		f.P("}")
		f.P("return &msg, nil")
	}
	f.P("default:")
	f.P(`return nil, fmt.Errorf("unmarshal `, d.Name(), ` frame: ID not in database: %d", f.ID)`)
	f.P("}")
	f.P("}")
	f.P()
	for _, m := range d.Messages {
		f.P("type ", m.Name, "Descriptor struct{")
		f.P("*descriptor.Message")
		for _, s := range m.Signals {
			f.P(s.Name, " *descriptor.Signal")
		}
		f.P("}")
		f.P()
	}
	f.P("// Database returns the ", d.Name(), " database descriptor.")
	f.P("func (md *MessagesDescriptor) Database() *descriptor.Database {")
	f.P("return d")
	f.P("}")
	f.P()
	f.P("var nd = &NodesDescriptor{")
	for ni, n := range d.Nodes {
		f.P(n.Name, ": d.Nodes[", ni, "],")
	}
	f.P("}")
	f.P()
	f.P("var md = &MessagesDescriptor{")
	for mi, m := range d.Messages {
		f.P(m.Name, ": &", m.Name, "Descriptor{")
		f.P("Message: d.Messages[", mi, "],")
		for si, s := range m.Signals {
			f.P(s.Name, ": d.Messages[", mi, "].Signals[", si, "],")
		}
		f.P("},")
	}
	f.P("}")
	f.P()
	f.P("var d = ")
	f.Dump(d)
	f.P()
}

func MarshalFrame(f *File, m *descriptor.Message) {
	f.P("// Frame returns a CAN frame representing the message.")
	f.P("func (m *", messageStruct(m), ") Frame() can.Frame {")
	f.P("md := ", messageDescriptor(m))
	f.P("f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}")
	for _, s := range m.Signals {
		if s.IsMultiplexed {
			continue
		}
		f.P(
			"md.", s.Name, ".Marshal", signalSuperType(s),
			"(&f.Data, ", signalPrimitiveSuperType(s), "(m.", signalField(s), "))",
		)
	}
	if mux, ok := m.MultiplexerSignal(); ok {
		for _, s := range m.Signals {
			if !s.IsMultiplexed {
				continue
			}
			f.P("if m.", signalField(mux), " == ", s.MultiplexerValue, " {")
			f.P(
				"md.", s.Name, ".Marshal", signalSuperType(s), "(&f.Data, ", signalPrimitiveSuperType(s),
				"(m.", signalField(s), "))",
			)
			f.P("}")
		}
	}
	f.P("return f")
	f.P("}")
	f.P()
	f.P("// MarshalFrame encodes the message as a CAN frame.")
	f.P("func (m *", messageStruct(m), ") MarshalFrame() (can.Frame, error) {")
	f.P("return m.Frame(), nil")
	f.P("}")
	f.P()
}

func UnmarshalFrame(f *File, m *descriptor.Message) {
	f.P("// UnmarshalFrame decodes the message from a CAN frame.")
	f.P("func (m *", messageStruct(m), ") UnmarshalFrame(f can.Frame) error {")
	f.P("md := ", messageDescriptor(m))
	// generate frame checks
	id := func(isExtended bool) string {
		if isExtended {
			return "extended ID"
		}
		return "standard ID"
	}
	f.P("switch {")
	f.P("case f.ID != md.ID:")
	f.P(`return fmt.Errorf(`)
	f.P(`"unmarshal `, m.Name, `: expects ID `, m.ID, ` (got %s with ID %d)", f.String(), f.ID,`)
	f.P(`)`)
	f.P("case f.Length != md.Length:")
	f.P(`return fmt.Errorf(`)
	f.P(`"unmarshal `, m.Name, `: expects length `, m.Length, ` (got %s with length %d)", f.String(), f.Length,`)
	f.P(`)`)
	f.P("case f.IsRemote:")
	f.P(`return fmt.Errorf(`)
	f.P(`"unmarshal `, m.Name, `: expects non-remote frame (got remote frame %s)", f.String(),`)
	f.P(`)`)
	f.P("case f.IsExtended != md.IsExtended:")
	f.P(`return fmt.Errorf(`)
	f.P(`"unmarshal `, m.Name, `: expects `, id(m.IsExtended), ` (got %s with `, id(!m.IsExtended), `)", f.String(),`)
	f.P(`)`)
	f.P("}")
	if len(m.Signals) == 0 {
		f.P("return nil")
		f.P("}")
		return
	}
	// generate non-multiplexed signal unmarshaling
	for _, s := range m.Signals {
		if s.IsMultiplexed {
			continue
		}
		f.P("m.", signalField(s), " = ", signalType(m, s), "(md.", s.Name, ".Unmarshal", signalSuperType(s), "(f.Data))")
	}
	// generate multiplexed signal unmarshaling
	if mux, ok := m.MultiplexerSignal(); ok {
		for _, s := range m.Signals {
			if !s.IsMultiplexed {
				continue
			}
			f.P("if m.", signalField(mux), " == ", s.MultiplexerValue, " {")
			f.P("m.", signalField(s), " = ", signalType(m, s), "(md.", s.Name, ".Unmarshal", signalSuperType(s), "(f.Data))")
			f.P("}")
		}
	}
	f.P("return nil")
	f.P("}")
	f.P()
}

func Node(f *File, d *descriptor.Database, n *descriptor.Node) {
	rxMessages := collectRxMessages(d, n)
	txMessages := collectTxMessages(d, n)
	f.P("type ", nodeInterface(n), " interface {")
	f.P("sync.Locker")
	f.P("Tx() ", txGroupInterface(n))
	f.P("Rx() ", rxGroupInterface(n))
	f.P("Run(ctx context.Context) error")
	f.P("}")
	f.P()
	f.P("type ", rxGroupInterface(n), " interface {")
	f.P("http.Handler // for debugging")
	for _, m := range rxMessages {
		f.P(m.Name, "() ", rxMessageInterface(n, m))
	}
	f.P("}")
	f.P()
	f.P("type ", txGroupInterface(n), " interface {")
	f.P("http.Handler // for debugging")
	for _, m := range txMessages {
		f.P(m.Name, "() ", txMessageInterface(n, m))
	}
	f.P("}")
	f.P()
	for _, m := range rxMessages {
		f.P("type ", rxMessageInterface(n, m), " interface {")
		f.P(messageReaderInterface(m))
		f.P("ReceiveTime() time.Time")
		f.P("SetAfterReceiveHook(h func(context.Context) error)")
		f.P("}")
		f.P()
	}
	for _, m := range txMessages {
		f.P("type ", txMessageInterface(n, m), " interface {")
		f.P(messageReaderInterface(m))
		f.P(messageWriterInterface(m))
		f.P("TransmitTime() time.Time")
		f.P("Transmit(ctx context.Context) error")
		f.P("SetBeforeTransmitHook(h func(context.Context) error)")
		if m.SendType == descriptor.SendTypeCyclic {
			f.P("// SetCyclicTransmissionEnabled enables/disables cyclic transmission.")
			f.P("SetCyclicTransmissionEnabled(bool)")
			f.P("// IsCyclicTransmissionEnabled returns whether cyclic transmission is enabled/disabled.")
			f.P("IsCyclicTransmissionEnabled() bool")
		}
		f.P("}")
		f.P()
	}
	f.P("type ", nodeStruct(n), " struct {")
	f.P("sync.Mutex // protects all node state")
	f.P("network string")
	f.P("address string")
	f.P("rx ", rxGroupStruct(n))
	f.P("tx ", txGroupStruct(n))
	f.P("}")
	f.P()
	f.P("var _ ", nodeInterface(n), " = &", nodeStruct(n), "{}")
	f.P("var _ canrunner.Node = &", nodeStruct(n), "{}")
	f.P()
	f.P("func New", nodeInterface(n), "(network, address string) ", nodeInterface(n), " {")
	f.P("n := &", nodeStruct(n), "{network: network, address: address}")
	f.P("n.rx.parentMutex = &n.Mutex")
	f.P("n.tx.parentMutex = &n.Mutex")
	for _, m := range rxMessages {
		f.P("n.rx.", messageField(m), ".init()")
		f.P("n.rx.", messageField(m), ".Reset()")
	}
	for _, m := range txMessages {
		f.P("n.tx.", messageField(m), ".init()")
		f.P("n.tx.", messageField(m), ".Reset()")
	}
	f.P("return n")
	f.P("}")
	f.P()
	f.P("func (n *", nodeStruct(n), ") Run(ctx context.Context) error {")
	f.P("return canrunner.Run(ctx, n)")
	f.P("}")
	f.P()
	f.P("func (n *", nodeStruct(n), ") Rx() ", rxGroupInterface(n), " {")
	f.P("return &n.rx")
	f.P("}")
	f.P()
	f.P("func (n *", nodeStruct(n), ") Tx() ", txGroupInterface(n), " {")
	f.P("return &n.tx")
	f.P("}")
	f.P()
	f.P("type ", rxGroupStruct(n), " struct {")
	f.P("parentMutex *sync.Mutex")
	for _, m := range rxMessages {
		f.P(messageField(m), " ", rxMessageStruct(n, m))
	}
	f.P("}")
	f.P()
	f.P("var _ ", rxGroupInterface(n), " = &", rxGroupStruct(n), "{}")
	f.P()
	f.P("func (rx *", rxGroupStruct(n), ") ServeHTTP(w http.ResponseWriter, r *http.Request) {")
	f.P("rx.parentMutex.Lock()")
	f.P("defer rx.parentMutex.Unlock()")
	f.P("candebug.ServeMessagesHTTP(w, r, []generated.Message{")
	for _, m := range rxMessages {
		f.P("&rx.", messageField(m), ",")
	}
	f.P("})")
	f.P("}")
	f.P()
	for _, m := range rxMessages {
		f.P("func (rx *", rxGroupStruct(n), ") ", m.Name, "() ", rxMessageInterface(n, m), " {")
		f.P("return &rx.", messageField(m))
		f.P("}")
		f.P()
	}
	f.P()
	f.P("type ", txGroupStruct(n), " struct {")
	f.P("parentMutex *sync.Mutex")
	for _, m := range txMessages {
		f.P(messageField(m), " ", txMessageStruct(n, m))
	}
	f.P("}")
	f.P()
	f.P("var _ ", txGroupInterface(n), " = &", txGroupStruct(n), "{}")
	f.P()
	f.P("func (tx *", txGroupStruct(n), ") ServeHTTP(w http.ResponseWriter, r *http.Request) {")
	f.P("tx.parentMutex.Lock()")
	f.P("defer tx.parentMutex.Unlock()")
	f.P("candebug.ServeMessagesHTTP(w, r, []generated.Message{")
	for _, m := range txMessages {
		f.P("&tx.", messageField(m), ",")
	}
	f.P("})")
	f.P("}")
	f.P()
	for _, m := range txMessages {
		f.P("func (tx *", txGroupStruct(n), ") ", m.Name, "() ", txMessageInterface(n, m), " {")
		f.P("return &tx.", messageField(m))
		f.P("}")
		f.P()
	}
	f.P()
	f.P("func (n *", nodeStruct(n), ") Descriptor() *descriptor.Node {")
	f.P("return ", nodeDescriptor(n))
	f.P("}")
	f.P()
	f.P("func (n *", nodeStruct(n), ") Connect() (net.Conn, error) {")
	f.P("return socketcan.Dial(n.network, n.address)")
	f.P("}")
	f.P()
	f.P("func (n *", nodeStruct(n), ") ReceivedMessage(id uint32) (canrunner.ReceivedMessage, bool) {")
	f.P("switch id {")
	for _, m := range rxMessages {
		f.P("case ", m.ID, ":")
		f.P("return &n.rx.", messageField(m), ", true")
	}
	f.P("default:")
	f.P("return nil, false")
	f.P("}")
	f.P("}")
	f.P()
	f.P("func (n *", nodeStruct(n), ") TransmittedMessages() []canrunner.TransmittedMessage {")
	f.P("return []canrunner.TransmittedMessage{")
	for _, m := range txMessages {
		f.P("&n.tx.", messageField(m), ",")
	}
	f.P("}")
	f.P("}")
	f.P()
	for _, m := range rxMessages {
		f.P("type ", rxMessageStruct(n, m), " struct {")
		f.P(messageStruct(m))
		f.P("receiveTime time.Time")
		f.P("afterReceiveHook func(context.Context) error")
		f.P("}")
		f.P()
		f.P("func (m *", rxMessageStruct(n, m), ") init() {")
		f.P("m.afterReceiveHook = func(context.Context) error { return nil }")
		f.P("}")
		f.P()
		f.P("func (m *", rxMessageStruct(n, m), ") SetAfterReceiveHook(h func(context.Context) error) {")
		f.P("m.afterReceiveHook = h")
		f.P("}")
		f.P()
		f.P("func (m *", rxMessageStruct(n, m), ") AfterReceiveHook() func(context.Context) error {")
		f.P("return m.afterReceiveHook")
		f.P("}")
		f.P()
		f.P("func (m *", rxMessageStruct(n, m), ") ReceiveTime() time.Time {")
		f.P("return m.receiveTime")
		f.P("}")
		f.P()
		f.P("func (m *", rxMessageStruct(n, m), ") SetReceiveTime(t time.Time) {")
		f.P("m.receiveTime = t")
		f.P("}")
		f.P()
		f.P("var _ canrunner.ReceivedMessage = &", rxMessageStruct(n, m), "{}")
		f.P()
	}
	for _, m := range txMessages {
		f.P("type ", txMessageStruct(n, m), " struct {")
		f.P(messageStruct(m))
		f.P("transmitTime time.Time")
		f.P("beforeTransmitHook func(context.Context) error")
		f.P("isCyclicEnabled bool")
		f.P("wakeUpChan chan struct{}")
		f.P("transmitEventChan chan struct{}")
		f.P("}")
		f.P()
		f.P("var _ ", txMessageInterface(n, m), " = &", txMessageStruct(n, m), "{}")
		f.P("var _ canrunner.TransmittedMessage = &", txMessageStruct(n, m), "{}")
		f.P()
		f.P("func (m *", txMessageStruct(n, m), ") init() {")
		f.P("m.beforeTransmitHook = func(context.Context) error { return nil }")
		f.P("m.wakeUpChan = make(chan struct{}, 1)")
		f.P("m.transmitEventChan = make(chan struct{})")
		f.P("}")
		f.P()
		f.P("func (m *", txMessageStruct(n, m), ") SetBeforeTransmitHook(h func(context.Context) error) {")
		f.P("m.beforeTransmitHook = h")
		f.P("}")
		f.P()
		f.P("func (m *", txMessageStruct(n, m), ") BeforeTransmitHook() func(context.Context) error {")
		f.P("return m.beforeTransmitHook")
		f.P("}")
		f.P()
		f.P("func (m *", txMessageStruct(n, m), ") TransmitTime() time.Time {")
		f.P("return m.transmitTime")
		f.P("}")
		f.P()
		f.P("func (m *", txMessageStruct(n, m), ") SetTransmitTime(t time.Time) {")
		f.P("m.transmitTime = t")
		f.P("}")
		f.P()
		f.P("func (m *", txMessageStruct(n, m), ") IsCyclicTransmissionEnabled() bool {")
		f.P("return m.isCyclicEnabled")
		f.P("}")
		f.P()
		f.P("func (m *", txMessageStruct(n, m), ") SetCyclicTransmissionEnabled(b bool) {")
		f.P("m.isCyclicEnabled = b")
		f.P("select {")
		f.P("case m.wakeUpChan <-struct{}{}:")
		f.P("default:")
		f.P("}")
		f.P("}")
		f.P()
		f.P("func (m *", txMessageStruct(n, m), ") WakeUpChan() <-chan struct{} {")
		f.P("return m.wakeUpChan")
		f.P("}")
		f.P()
		f.P("func (m *", txMessageStruct(n, m), ") Transmit(ctx context.Context) error {")
		f.P("select {")
		f.P("case m.transmitEventChan <- struct{}{}:")
		f.P("return nil")
		f.P("case <-ctx.Done():")
		f.P(`return fmt.Errorf("event-triggered transmit of `, m.Name, `: %w", ctx.Err())`)
		f.P("}")
		f.P("}")
		f.P()
		f.P("func (m *", txMessageStruct(n, m), ") TransmitEventChan() <-chan struct{} {")
		f.P("return m.transmitEventChan")
		f.P("}")
		f.P()
		f.P("var _ canrunner.TransmittedMessage = &", txMessageStruct(n, m), "{}")
		f.P()
	}
}

func txGroupInterface(n *descriptor.Node) string {
	return n.Name + "_Tx"
}

func txGroupStruct(n *descriptor.Node) string {
	return "xxx_" + n.Name + "_Tx"
}

func rxGroupInterface(n *descriptor.Node) string {
	return n.Name + "_Rx"
}

func rxGroupStruct(n *descriptor.Node) string {
	return "xxx_" + n.Name + "_Rx"
}

func rxMessageInterface(n *descriptor.Node, m *descriptor.Message) string {
	return n.Name + "_Rx_" + m.Name
}

func rxMessageStruct(n *descriptor.Node, m *descriptor.Message) string {
	return "xxx_" + n.Name + "_Rx_" + m.Name
}

func txMessageInterface(n *descriptor.Node, m *descriptor.Message) string {
	return n.Name + "_Tx_" + m.Name
}

func txMessageStruct(n *descriptor.Node, m *descriptor.Message) string {
	return "xxx_" + n.Name + "_Tx_" + m.Name
}

func collectTxMessages(d *descriptor.Database, n *descriptor.Node) []*descriptor.Message {
	tx := make([]*descriptor.Message, 0, len(d.Messages))
	for _, m := range d.Messages {
		if m.SenderNode == n.Name && m.SendType != descriptor.SendTypeNone {
			tx = append(tx, m)
		}
	}
	return tx
}

func collectRxMessages(d *descriptor.Database, n *descriptor.Node) []*descriptor.Message {
	rx := make([]*descriptor.Message, 0, len(d.Messages))
Loop:
	for _, m := range d.Messages {
		for _, s := range m.Signals {
			for _, node := range s.ReceiverNodes {
				if node != n.Name {
					continue
				}
				rx = append(rx, m)
				continue Loop
			}
		}
	}
	return rx
}

func hasPhysicalRepresentation(s *descriptor.Signal) bool {
	hasScale := s.Scale != 0 && s.Scale != 1
	hasOffset := s.Offset != 0
	hasRange := s.Min != 0 || s.Max != 0
	var hasConstrainedRange bool
	if s.IsSigned {
		hasConstrainedRange = s.Min > float64(s.MinSigned()) || s.Max < float64(s.MaxSigned())
	} else {
		hasConstrainedRange = s.Min > 0 || s.Max < float64(s.MaxUnsigned())
	}
	return hasScale || hasOffset || hasRange && hasConstrainedRange
}

func hasCustomType(s *descriptor.Signal) bool {
	return len(s.ValueDescriptions) > 0
}

func hasSendType(d *descriptor.Database) bool {
	for _, m := range d.Messages {
		if m.SendType != descriptor.SendTypeNone {
			return true
		}
	}
	return false
}

func signalType(m *descriptor.Message, s *descriptor.Signal) string {
	if hasCustomType(s) {
		return m.Name + "_" + s.Name
	}
	return signalPrimitiveType(s).String()
}

func signalPrimitiveType(s *descriptor.Signal) types.Type {
	var t types.BasicKind
	switch {
	case s.Length == 1:
		t = types.Bool
	case s.Length <= 8 && s.IsSigned:
		t = types.Int8
	case s.Length <= 8:
		t = types.Uint8
	case s.Length <= 16 && s.IsSigned:
		t = types.Int16
	case s.Length <= 16:
		t = types.Uint16
	case s.Length <= 32 && s.IsSigned:
		t = types.Int32
	case s.Length <= 32:
		t = types.Uint32
	case s.Length <= 64 && s.IsSigned:
		t = types.Int64
	default:
		t = types.Uint64
	}
	return types.Typ[t]
}

func signalPrimitiveSuperType(s *descriptor.Signal) types.Type {
	var t types.BasicKind
	switch {
	case s.Length == 1:
		t = types.Bool
	case s.IsSigned:
		t = types.Int64
	default:
		t = types.Uint64
	}
	return types.Typ[t]
}

func signalSuperType(s *descriptor.Signal) string {
	switch {
	case s.Length == 1:
		return "Bool"
	case s.IsSigned:
		return "Signed"
	default:
		return "Unsigned"
	}
}

func nodeInterface(n *descriptor.Node) string {
	return n.Name
}

func nodeStruct(n *descriptor.Node) string {
	return "xxx_" + n.Name
}

func messageStruct(m *descriptor.Message) string {
	return m.Name
}

func messageReaderInterface(m *descriptor.Message) string {
	return m.Name + "Reader"
}

func messageWriterInterface(m *descriptor.Message) string {
	return m.Name + "Writer"
}

func messageField(m *descriptor.Message) string {
	return "xxx_" + m.Name
}

func signalField(s *descriptor.Signal) string {
	return "xxx_" + s.Name
}

func nodeDescriptor(n *descriptor.Node) string {
	return "Nodes()." + n.Name
}

func messageDescriptor(m *descriptor.Message) string {
	return "Messages()." + m.Name
}

func signalDescriptor(m *descriptor.Message, s *descriptor.Signal) string {
	return messageDescriptor(m) + "." + s.Name
}
