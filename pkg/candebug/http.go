package candebug

import (
	"bytes"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/blueinnovationsgroup/can-go/pkg/cantext"
	"github.com/blueinnovationsgroup/can-go/pkg/descriptor"
	"github.com/blueinnovationsgroup/can-go/pkg/generated"
)

func ServeMessagesHTTP(w http.ResponseWriter, r *http.Request, msgs []generated.Message) {
	base := path.Base(r.URL.Path)
	// if path ends with a message name, serve only that message
	for _, m := range msgs {
		if m.Descriptor().Name == base {
			serveMessagesHTTP(w, r, []generated.Message{m})
			return
		}
	}
	serveMessagesHTTP(w, r, msgs)
}

func serveMessagesHTTP(w http.ResponseWriter, _ *http.Request, msgs []generated.Message) {
	var buf []byte
	for i, m := range msgs {
		buf = appendMessage(buf, m)
		if i != len(msgs)-1 {
			buf = append(buf, "\n\n\n"...)
		}
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf)
}

func appendMessage(buf []byte, m generated.Message) []byte {
	name := m.Descriptor().Name
	sep := append(bytes.Repeat([]byte{'='}, len(name)), '\n')
	buf = append(buf, name...)
	buf = append(buf, '\n')
	buf = append(buf, sep...)
	buf = cantext.AppendID(buf, m.Descriptor())
	buf = append(buf, '\n')
	buf = cantext.AppendSender(buf, m.Descriptor())
	buf = append(buf, '\n')
	buf = cantext.AppendSendType(buf, m.Descriptor())
	buf = append(buf, '\n')
	if m.Descriptor().SendType == descriptor.SendTypeCyclic {
		if enabler, ok := m.(interface{ IsCyclicTransmissionEnabled() bool }); ok {
			buf = append(buf, "Enabled: "...)
			buf = strconv.AppendBool(buf, enabler.IsCyclicTransmissionEnabled())
			buf = append(buf, '\n')
		}
		buf = cantext.AppendCycleTime(buf, m.Descriptor())
		buf = append(buf, '\n')
	}
	if m.Descriptor().DelayTime != 0 {
		buf = cantext.AppendDelayTime(buf, m.Descriptor())
		buf = append(buf, '\n')
	}
	buf = append(buf, sep...)
	if timer, ok := m.(interface{ ReceiveTime() time.Time }); ok {
		buf = append(buf, "Received: "...)
		buf = appendTime(buf, timer.ReceiveTime())
		buf = append(buf, '\n')
		buf = append(buf, sep...)
	}
	if timer, ok := m.(interface{ TransmitTime() time.Time }); ok {
		buf = append(buf, "Transmitted: "...)
		buf = appendTime(buf, timer.TransmitTime())
		buf = append(buf, '\n')
		buf = append(buf, sep...)
	}
	f := m.Frame()
	for i, s := range m.Descriptor().Signals {
		buf = cantext.AppendSignal(buf, s, f.Data)
		if i < len(m.Descriptor().Signals)-1 {
			buf = append(buf, '\n')
		}
	}
	return buf
}

func appendTime(buf []byte, t time.Time) []byte {
	if t.IsZero() {
		buf = append(buf, "never"...)
		return buf
	}
	buf = append(buf, time.Since(t).String()...)
	buf = append(buf, " ago ("...)
	buf = t.AppendFormat(buf, "15:04:05.000000000")
	buf = append(buf, ")"...)
	return buf
}
