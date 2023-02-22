package canrunner

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/blueinnovationsgroup/can-go"
	"github.com/blueinnovationsgroup/can-go/internal/clock"
	"github.com/blueinnovationsgroup/can-go/pkg/descriptor"
	"github.com/blueinnovationsgroup/can-go/pkg/generated"
	"github.com/blueinnovationsgroup/can-go/pkg/socketcan"
	"golang.org/x/sync/errgroup"
)

// defaultSendTimeout is the send timeout used for messages without a cycle time.
const defaultSendTimeout = time.Second

// Node is an interface for a CAN node to be run by the runner.
type Node interface {
	sync.Locker
	Connect() (net.Conn, error)
	Descriptor() *descriptor.Node
	TransmittedMessages() []TransmittedMessage
	ReceivedMessage(id uint32) (ReceivedMessage, bool)
}

// TransmittedMessage is an interface for a message to be transmitted by the runner.
type TransmittedMessage interface {
	generated.Message
	// SetTransmitTime sets the time the message was last transmitted.
	SetTransmitTime(time.Time)
	// IsCyclicTransmissionEnabled returns true when cyclic transmission is enabled.
	IsCyclicTransmissionEnabled() bool
	// WakeUpChan returns a channel for waking up and checking if cyclic transmission is enabled.
	WakeUpChan() <-chan struct{}
	// TransmitEventChan returns channel for event-based transmission of the message.
	TransmitEventChan() <-chan struct{}
	// BeforeTransmitHook returns a function to be called before the message is transmitted.
	//
	// If the hook returns an error, the transmitter halt.
	BeforeTransmitHook() func(context.Context) error
}

// ReceivedMessage is an interface for a message to be received by the runner.
type ReceivedMessage interface {
	generated.Message
	// SetReceiveTime sets the time the message was last received.
	SetReceiveTime(time.Time)
	// AfterReceiveHook returns a function to be called after the message has been received.
	//
	// If the hook returns an error, the receiver will halt.
	AfterReceiveHook() func(context.Context) error
}

// FrameTransmitter is an interface for the the CAN frame transmitter used by the runner.
type FrameTransmitter interface {
	TransmitFrame(context.Context, can.Frame) error
}

// FrameReceiver is an interface for the CAN frame receiver used by the runner.
type FrameReceiver interface {
	Receive() bool
	Frame() can.Frame
	Err() error
}

func Run(ctx context.Context, n Node) error {
	conn, err := n.Connect()
	if err != nil {
		return fmt.Errorf("run %s node: %w", n.Descriptor().Name, err)
	}
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		<-ctx.Done()
		return conn.Close()
	})
	g.Go(func() error {
		rx := socketcan.NewReceiver(conn)
		return RunMessageReceiver(ctx, rx, n, clock.System())
	})
	for _, m := range n.TransmittedMessages() {
		m := m
		g.Go(func() error {
			tx := socketcan.NewTransmitter(conn)
			return RunMessageTransmitter(ctx, tx, n, m, clock.System())
		})
	}
	if err := g.Wait(); err != nil {
		if strings.Contains(err.Error(), "closed") {
			return nil
		}
		return fmt.Errorf("run %s node: %w", n.Descriptor().Name, err)
	}
	return nil
}

func RunMessageReceiver(ctx context.Context, rx FrameReceiver, n Node, c clock.Clock) error {
	for rx.Receive() {
		f := rx.Frame()
		m, ok := n.ReceivedMessage(f.ID)
		if !ok {
			continue
		}
		n.Lock()
		hook := m.AfterReceiveHook()
		m.SetReceiveTime(c.Now())
		err := m.UnmarshalFrame(f)
		n.Unlock()
		if err != nil {
			return fmt.Errorf("receiver: %w", err)
		}
		if err := hook(ctx); err != nil {
			return fmt.Errorf("receiver: %w", err)
		}
	}
	if err := rx.Err(); err != nil {
		return fmt.Errorf("receiver: %w", err)
	}
	return nil
}

func RunMessageTransmitter(
	ctx context.Context,
	tx FrameTransmitter,
	l sync.Locker,
	m TransmittedMessage,
	c clock.Clock,
) error {
	sendTimeout := m.Descriptor().CycleTime
	if sendTimeout == 0 {
		sendTimeout = defaultSendTimeout
	}
	var cyclicTransmissionTicker *time.Ticker
	var cyclicTransmissionTickChan <-chan time.Time
	enableCyclicTransmission := func() {
		isCyclic := m.Descriptor().SendType == descriptor.SendTypeCyclic
		hasCycleTime := m.Descriptor().CycleTime > 0
		if !isCyclic || !hasCycleTime || cyclicTransmissionTicker != nil {
			return
		}
		cyclicTransmissionTicker = time.NewTicker(m.Descriptor().CycleTime)
		cyclicTransmissionTickChan = cyclicTransmissionTicker.C
	}
	disableCyclicTransmission := func() {
		if cyclicTransmissionTicker == nil {
			return
		}
		cyclicTransmissionTicker.Stop()
		cyclicTransmissionTicker = nil
	}
	setCyclicTransmission := func() {
		l.Lock()
		isCyclicTransmissionEnabled := m.IsCyclicTransmissionEnabled()
		l.Unlock()
		if isCyclicTransmissionEnabled {
			enableCyclicTransmission()
		} else {
			disableCyclicTransmission()
		}
	}
	transmit := func() error {
		l.Lock()
		hook := m.BeforeTransmitHook()
		m.SetTransmitTime(c.Now())
		l.Unlock()
		if err := hook(ctx); err != nil {
			return fmt.Errorf("%s transmitter: %w", m.Descriptor().Name, err)
		}
		l.Lock()
		f := m.Frame()
		l.Unlock()
		ctx, cancel := context.WithTimeout(ctx, sendTimeout)
		err := tx.TransmitFrame(ctx, f)
		cancel()
		if err != nil {
			return fmt.Errorf("%s transmitter: %w", m.Descriptor().Name, err)
		}
		return nil
	}
	ctxDone := ctx.Done()
	transmitEventChan := m.TransmitEventChan()
	setCyclicTransmission()
	wakeUpChan := m.WakeUpChan()
	for {
		select {
		case <-ctxDone:
			return nil
		case <-wakeUpChan:
			setCyclicTransmission()
		case <-transmitEventChan:
			if err := transmit(); err != nil {
				return err
			}
		case <-cyclicTransmissionTickChan:
			if err := transmit(); err != nil {
				return err
			}
		}
	}
}
