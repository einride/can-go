package socketcan

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"go.einride.tech/can"
	"golang.org/x/sync/errgroup"
)

type emulatorCfg struct {
	address string
	logger  *log.Logger
}

func defaultCfg() emulatorCfg {
	stdLogger := log.New(os.Stderr, "emulator: ", log.Lshortfile|log.Ltime)
	return emulatorCfg{
		address: "239.64.142.206:0",
		logger:  stdLogger,
	}
}

// EmulatorOption represents a way to configure an Emulator prior to creating it.
type EmulatorOption func(*emulatorCfg)

// WithMulticastAddress sets the address for the multicast group that the Emulator should listen on.
// A multicast address starts with 239.x.x.x, and using an address that does not conform to this
// will lead to undefined behavior.
func WithMulticastAddress(address string) EmulatorOption {
	return func(cfg *emulatorCfg) {
		cfg.address = address
	}
}

// WithLogger makes the Emulator print out status messages with the provided logger.
func WithLogger(l *log.Logger) EmulatorOption {
	return func(cfg *emulatorCfg) {
		cfg.logger = l
	}
}

// NoLogger disables logging in the Emulator.
func NoLogger(cfg *emulatorCfg) {
	cfg.logger = log.New(&writeSink{}, "", log.LstdFlags)
}

// writeSink is an io.Writer which does not write to anything.
//
// Can be thought of as a /dev/null for writers.
type writeSink struct{}

// Write returns without actually writing to anything.
func (w *writeSink) Write(buf []byte) (int, error) {
	return len(buf), nil
}

// Emulator emulates a CAN bus.
//
// Emulator emulates a CAN bus by using UDP multicast. The emulator itself
// does not own the multicast group but rather establishes a common
// address/port pair for the CAN bus to be emulated on.
// Emulator exposes a thread-safe API to callees and may therefore be
// shared among different goroutines.
type Emulator struct {
	transceiver   *udpTxRx
	logger        *log.Logger
	reqSenderChan chan chan int
	closeChan     chan struct{}
	sync.Mutex
	g  *errgroup.Group
	rg errgroup.Group
}

// NewEmulator creates an Emulator to emulate a CAN bus.
//
// If no error is returned it is safe to `socketcan.Dial` the address
// of the Emulator. The emulator will default to using multicast group `239.64.142.206`
// with a random port that's decided when calling Emulator
//
// N.B. It is not possible to simply use `net.Dial` as for UDP multicast both
// a transmitting connection and a writing connection. This is handled
// by `socketcan.Dial` under the hood.
func NewEmulator(options ...EmulatorOption) (*Emulator, error) {
	cfg := defaultCfg()
	for _, update := range options {
		update(&cfg)
	}
	c, err := udpTransceiver("udp", cfg.address)
	if err != nil {
		return nil, err
	}
	return &Emulator{
		transceiver:   c,
		logger:        cfg.logger,
		closeChan:     make(chan struct{}),
		reqSenderChan: make(chan chan int),
	}, nil
}

// Run an Emulator.
//
// This starts the listener and waits until the context is canceled
// before tidying up.
func (e *Emulator) Run(ctx context.Context) error {
	e.Lock()
	e.g, ctx = errgroup.WithContext(ctx)
	ctxDone := ctx.Done()

	// Listen for incoming frames.
	// Keep track of unique senders, and notify newSenderChan.
	newSenderChan := make(chan struct{})
	e.g.Go(func() error {
		e.logger.Printf("waiting for SocketCAN connection requests on udp://%s\n", e.Addr().String())
		registeredSenders := make(map[string]bool)
		for {
			buffer := make([]byte, 8096)
			_, _, src, err := e.transceiver.rx.ReadFrom(buffer)
			if err != nil {
				if isClosedError(err) {
					return nil
				}
				return fmt.Errorf("read from udp: %w", err)
			}
			if !registeredSenders[src.String()] {
				e.logger.Printf("received first frame from %s", src.String())
				registeredSenders[src.String()] = true
				select {
				case <-ctxDone:
					return nil
				case newSenderChan <- struct{}{}:
				}
			}
		}
	})

	// Close multicast listener when ctx is canceled
	e.g.Go(func() error {
		<-ctxDone
		e.logger.Println("closing SocketCAN listener...")
		e.Lock()
		defer e.Unlock()
		return e.transceiver.Close()
	})

	// Stop all started receivers when ctx is canceled
	e.g.Go(func() error {
		<-ctxDone
		e.logger.Println("stopping receivers...")
		close(e.closeChan)
		e.Lock()
		defer e.Unlock()
		return e.rg.Wait()
	})

	// Keep track of the number of unique senders of the received frames, when the number of senders
	// are requested on the reqSenderChan, send them on the provided channel.
	e.g.Go(func() error {
		nSenders := 0
		for {
			select {
			case <-ctxDone:
				return nil
			case <-newSenderChan:
				nSenders++
			case req := <-e.reqSenderChan:
				req <- nSenders
			}
		}
	})
	e.Unlock()
	e.logger.Println("started emulator, waiting for cancel signal")
	return e.g.Wait()
}

// Addr returns the address of the Emulator's multicast group.
func (e *Emulator) Addr() net.Addr {
	return e.transceiver.RemoteAddr()
}

// Receiver returns a Receiver connected to the Emulator.
//
// The emulator owns the underlying network connection an
// will close it when the emulator is closed.
func (e *Emulator) Receiver() (*Receiver, error) {
	conn, err := udpTransceiver(e.Addr().Network(), e.Addr().String())
	if err != nil {
		return nil, err
	}
	e.Lock()
	e.rg.Go(func() error {
		<-e.closeChan
		return conn.Close()
	})
	e.Unlock()
	return NewReceiver(conn), nil
}

// TransmitFrame sends a CAN frame to the Emulator's multicast group.
func (e *Emulator) TransmitFrame(ctx context.Context, f can.Frame) error {
	conn, err := udpTransceiver(e.Addr().Network(), e.Addr().String())
	if err != nil {
		return fmt.Errorf("transmit CAN frame: %w", err)
	}
	errChan := make(chan error)
	go func() {
		if err := NewTransmitter(conn).TransmitFrame(ctx, f); err != nil {
			errChan <- fmt.Errorf("transmit CAN frame: %w", err)
		}
		close(errChan)
	}()
	select {
	case <-ctx.Done():
		_ = conn.Close()
		return ctx.Err()
	case err := <-errChan:
		_ = conn.Close()
		if err != nil {
			return fmt.Errorf("emulator: %w", err)
		}
		return nil
	}
}

// TransmitMessage sends a CAN message to every emulator connection.
func (e *Emulator) TransmitMessage(ctx context.Context, m can.Message) error {
	f, err := m.MarshalFrame()
	if err != nil {
		return fmt.Errorf("transmit CAN message: %w", err)
	}
	if err := e.TransmitFrame(ctx, f); err != nil {
		return fmt.Errorf("transmit CAN message: %w", err)
	}
	return nil
}

// WaitForSenders waits until either, n unique senders have been sending messages to the
// multicast group, or the timeout is reached.
func (e *Emulator) WaitForSenders(n int, timeout time.Duration) error {
	reqChan := make(chan int)
	timeoutChannel := time.After(timeout)
	for {
		select {
		case <-timeoutChannel:
			return fmt.Errorf("emulator timeout waiting for senders")
		case e.reqSenderChan <- reqChan:
			conns := <-reqChan
			if conns < n {
				// We don't want to keep the emulator
				// busy with our requests all the time.
				time.Sleep(time.Millisecond)
				continue
			}
			return nil
		}
	}
}

func isClosedError(e error) bool {
	if e == nil {
		return false
	}
	return strings.Contains(e.Error(), "closed")
}
