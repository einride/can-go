package socketcan

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.einride.tech/can"
	"golang.org/x/sync/errgroup"
)

func TestEmulate_Close(t *testing.T) {
	// Given: an emulator
	e, err := NewEmulator()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	// When: I start the emulator
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return e.Run(ctx)
	})
	// Then: I should be able to close it
	require.NoError(t, g.Wait())
}

func TestEmulate_SendToAll(t *testing.T) {
	for _, tt := range []struct {
		receivers int
	}{
		{receivers: 1},
		{receivers: 5},
		{receivers: 100},
	} {
		tt := tt
		t.Run(fmt.Sprintf("receivers:%v", tt.receivers), func(t *testing.T) {
			// Given: A listener with an Emulator
			ctx, cancel := context.WithCancel(context.Background())
			eg, eCtx := errgroup.WithContext(ctx)
			e, err := NewEmulator()
			require.NoError(t, err)
			eg.Go(func() error {
				return e.Run(eCtx)
			})
			// When: I start multiple receivers connected to the Emulator
			g := errgroup.Group{}
			for i := 0; i < tt.receivers; i++ {
				r, err := e.Receiver()
				require.NoError(t, err)
				g.Go(func() error {
					if ok := r.Receive(); !ok {
						return fmt.Errorf("failed to receive CAN frame: %w", r.Err())
					}
					if r.HasErrorFrame() {
						return fmt.Errorf("received error frame: %v", r.ErrorFrame())
					}
					return r.Err()
				})
			}
			// And then the emulator transmits a CAN frame
			txFrame := can.Frame{ID: 42, Length: 4, Data: can.Data{1, 2, 3, 4}}
			err = e.TransmitFrame(context.Background(), txFrame)
			require.NoError(t, err)
			// Then: Every receiver should receive the frame and not return an error
			require.NoError(t, g.Wait())
			cancel()
			require.NoError(t, eg.Wait())
		})
	}
}

func TestEmulate_ConnectMany(t *testing.T) {
	for _, tt := range []struct {
		noTransmitters int
		canFrames      []can.Frame
	}{
		{
			noTransmitters: 1,
			canFrames: []can.Frame{
				{ID: 42},
				{ID: 43, Length: 4, Data: can.Data{1, 2, 3, 4}},
			},
		},
		{
			noTransmitters: 10,
			canFrames: []can.Frame{
				{ID: 42},
				{ID: 43, Length: 4, Data: can.Data{1, 2, 3, 4}},
				{ID: 44, IsRemote: true},
			},
		},
		{
			noTransmitters: 50,
			canFrames: []can.Frame{
				{ID: 42},
				{ID: 43, Length: 4, Data: can.Data{1, 2, 3, 4}},
				{ID: 44, IsRemote: true},
				{ID: 45, Length: 7, Data: can.Data{1, 2, 3, 4, 5, 6, 7}},
				{ID: 46, IsExtended: false},
				{ID: 47, Length: 1, Data: can.Data{1}},
				{ID: 48, IsRemote: false},
			},
		},
	} {
		tt := tt
		name := fmt.Sprintf("transmitters:%v,frames:%v", tt.noTransmitters, len(tt.canFrames))
		t.Run(name, func(t *testing.T) {
			// Given: A listener with an Emulator
			e, err := NewEmulator(NoLogger)
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(context.Background())
			eg, eCtx := errgroup.WithContext(ctx)
			eg.Go(func() error {
				return e.Run(eCtx)
			})
			r, err := e.Receiver()
			require.NoError(t, err)
			receiver := errgroup.Group{}
			receiver.Go(func() error {
				for i := 0; i < len(tt.canFrames)*tt.noTransmitters; i++ {
					i := i
					if ok := r.Receive(); !ok {
						cancel()
						require.NoError(t, eg.Wait())
						t.Fatal("Not all CAN frames were received", i, r.Err())
					}
					require.Contains(t, tt.canFrames, r.Frame())
				}
				return nil
			})
			// When: I connect multiple transmitters and transmit CAN frame on every transmitter
			transmits, txCtx := errgroup.WithContext(ctx)
			for i := 0; i < tt.noTransmitters; i++ {
				transmits.Go(func() error {
					conn, err := DialContext(txCtx, e.Addr().Network(), e.Addr().String())
					if err != nil {
						return err
					}
					tx := NewTransmitter(conn)
					for _, frame := range tt.canFrames {
						if err := tx.TransmitFrame(txCtx, frame); err != nil {
							log.Printf("failed to transmit frame: %+v\n", err)
							return err
						}
					}
					return conn.Close()
				})
			}
			require.NoError(t, transmits.Wait())
			// Then: Every CAN frame should have been delivered to the emulator
			require.NoError(t, receiver.Wait())
			cancel()
			require.NoError(t, eg.Wait())
		})
	}
}

func TestEmulate_SendReceive(t *testing.T) {
	for _, tt := range []struct {
		transmitters int
		receivers    int
	}{
		{
			transmitters: 1,
			receivers:    2,
		},
		{
			transmitters: 10,
			receivers:    50,
		},
		{
			transmitters: 50,
			receivers:    50,
		},
	} {
		tt := tt
		name := fmt.Sprintf("transmitters: %v,receivers: %v", tt.transmitters, tt.receivers)
		t.Run(name, func(t *testing.T) {
			// Given: A listener and an emulator
			e, err := NewEmulator()
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(context.Background())
			eg, eCtx := errgroup.WithContext(ctx)
			eg.Go(func() error {
				return e.Run(eCtx)
			})
			canFrames := make(map[uint32]can.Frame)
			canFrames[42] = can.Frame{ID: 42}
			canFrames[43] = can.Frame{ID: 43, IsRemote: true}
			canFrames[44] = can.Frame{ID: 44, IsExtended: true}
			// When: I start a number of receivers
			rx := errgroup.Group{}
			for i := 0; i < tt.receivers; i++ {
				r, err := e.Receiver()
				require.NoError(t, err)
				rx.Go(func() error {
					for i := 0; i < tt.transmitters*len(canFrames); i++ {
						if ok := r.Receive(); !ok {
							return fmt.Errorf("receive frames: %w", r.Err())
						}
						if r.HasErrorFrame() {
							return fmt.Errorf("received error frame: %v", r.ErrorFrame())
						}
						if _, ok := canFrames[r.Frame().ID]; !ok {
							return fmt.Errorf("received unexpected frame: %v", r.Frame())
						}
					}
					return nil
				})
			}
			// And then start a number of transmitters that will transmit a number of CAN frames
			tx, txCtx := errgroup.WithContext(ctx)
			for i := 0; i < tt.transmitters; i++ {
				conn, err := DialContext(txCtx, e.Addr().Network(), e.Addr().String())
				require.NoError(t, err)
				tx.Go(func() (err error) {
					t := NewTransmitter(conn)
					for _, f := range canFrames {
						if err := t.TransmitFrame(txCtx, f); err != nil {
							return fmt.Errorf("transmit frame: %w", err)
						}
					}
					if err := conn.Close(); err != nil {
						return fmt.Errorf("close transmitter: %w", err)
					}
					return nil
				})
			}
			// Then: The transmissions should not fail
			require.NoError(t, tx.Wait())
			// And every receiver should receive every CAN frame
			require.NoError(t, rx.Wait())
			cancel()
			require.NoError(t, eg.Wait())
		})
	}
}

func TestEmulator_Isolation(t *testing.T) {
	// Given 5 separate emulators
	const nEmulators = 5
	var emulators []*Emulator
	ctx, cancel := context.WithCancel(context.Background())
	eg, eCtx := errgroup.WithContext(ctx)
	for i := 0; i < nEmulators; i++ {
		e, err := NewEmulator()
		require.NoError(t, err)
		emulators = append(emulators, e)
		eg.Go(func() error {
			return e.Run(eCtx)
		})
	}
	// When starting one transmitter/receiver pair per emulator sending 10 frames
	const nFrames = 10
	rx := errgroup.Group{}
	tx := errgroup.Group{}
	for i := 0; i < nEmulators; i++ {
		i := i
		r, err := emulators[i].Receiver()
		require.NoError(t, err)
		rx.Go(func() error {
			for j := 0; j < nFrames; j++ {
				if ok := r.Receive(); !ok {
					return fmt.Errorf("receive frame: %w", r.Err())
				}
				if r.HasErrorFrame() {
					return fmt.Errorf("received error frame: %v", r.ErrorFrame())
				}
				if r.Frame().ID != uint32(i) {
					return fmt.Errorf("receiver(%v) received unexpected frame: %v", i, r.Frame())
				}
			}
			return nil
		})
		for j := 0; j < nFrames; j++ {
			frame := can.Frame{ID: uint32(i)}
			tx.Go(func() error {
				return emulators[i].TransmitFrame(context.Background(), frame)
			})
		}
	}
	// Then all transmitted frames should be received by correct receiver
	require.NoError(t, rx.Wait())
	require.NoError(t, tx.Wait())
	cancel()
	require.NoError(t, eg.Wait())
}

func TestEmulator_WaitForSenders(t *testing.T) {
	// Given a started emulator
	ctx, cancel := context.WithCancel(context.Background())
	eg, eCtx := errgroup.WithContext(ctx)
	e, err := NewEmulator()
	require.NoError(t, err)
	eg.Go(func() error {
		return e.Run(eCtx)
	})
	// When one transmitter is transmitting a frame
	txg := errgroup.Group{}
	txg.Go(func() error {
		return e.TransmitFrame(context.Background(), can.Frame{ID: 1234})
	})
	// Then WaitForSenders should return without an error
	err = e.WaitForSenders(1, time.Second)
	require.NoError(t, err)
	require.NoError(t, txg.Wait())
	cancel()
	require.NoError(t, eg.Wait())
}

func TestEmulator_WaitForSenders_Multiple(t *testing.T) {
	// Given a started emulator
	ctx, cancel := context.WithCancel(context.Background())
	eg, eCtx := errgroup.WithContext(ctx)
	e, err := NewEmulator()
	require.NoError(t, err)
	eg.Go(func() error {
		return e.Run(eCtx)
	})
	// When one transmitter is transmitting a frame
	txg := errgroup.Group{}
	txg.Go(func() error {
		return e.TransmitFrame(context.Background(), can.Frame{ID: 1234})
	})
	txg.Go(func() error {
		return e.TransmitFrame(context.Background(), can.Frame{ID: 4321})
	})
	// Then WaitForSenders should return without an error
	err = e.WaitForSenders(2, time.Second)
	require.NoError(t, err)
	require.NoError(t, txg.Wait())
	cancel()
	require.NoError(t, eg.Wait())
}

func TestEmulator_WaitForSenders_Timeout(t *testing.T) {
	// Given a started emulator
	ctx, cancel := context.WithCancel(context.Background())
	eg, eCtx := errgroup.WithContext(ctx)
	e, err := NewEmulator()
	require.NoError(t, err)
	eg.Go(func() error {
		return e.Run(eCtx)
	})
	// When no transmitters have connected and transmitted frames
	// Then WaitForSenders should timeout
	err = e.WaitForSenders(1, 100*time.Millisecond)
	require.Error(t, err)
	cancel()
	require.NoError(t, eg.Wait())
}
