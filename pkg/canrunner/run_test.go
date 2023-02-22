package canrunner_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/blueinnovationsgroup/can-go"
	"github.com/blueinnovationsgroup/can-go/internal/mocks/gen/mockcanrunner"
	"github.com/blueinnovationsgroup/can-go/internal/mocks/gen/mockclock"
	"github.com/blueinnovationsgroup/can-go/pkg/canrunner"
	"github.com/blueinnovationsgroup/can-go/pkg/descriptor"
	"github.com/golang/mock/gomock"
	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"
)

func TestRunMessageReceiver_NoMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	rx := mockcanrunner.NewMockFrameReceiver(ctrl)
	node := mockcanrunner.NewMockNode(ctrl)
	clock := mockclock.NewMockClock(ctrl)
	ctx := context.Background()
	// when the first receive fails
	rx.EXPECT().Receive().Return(false)
	rx.EXPECT().Err().Return(os.ErrClosed)
	// then an error is returned
	assert.Assert(t, errors.Is(canrunner.RunMessageReceiver(ctx, rx, node, clock), os.ErrClosed))
}

func TestRunMessageReceiver_ReceiveMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	rx := mockcanrunner.NewMockFrameReceiver(ctrl)
	node := mockcanrunner.NewMockNode(ctrl)
	clock := mockclock.NewMockClock(ctrl)
	msg := mockcanrunner.NewMockReceivedMessage(ctrl)
	ctx := context.Background()
	// when the first receive succeeds
	frame := can.Frame{ID: 42}
	rx.EXPECT().Receive().Return(true)
	rx.EXPECT().Frame().Return(frame)
	// then the receiver should do a message lookup
	node.EXPECT().ReceivedMessage(frame.ID).Return(msg, true)
	// and the node should be locked
	node.EXPECT().Lock()
	// and the message should be queried for a hook with the same context
	afterReceiveHook := func(c context.Context) error {
		assert.DeepEqual(t, ctx, c)
		return nil
	}
	msg.EXPECT().AfterReceiveHook().Return(afterReceiveHook)
	// and the receive time should be set
	now := time.Unix(0, 1)
	clock.EXPECT().Now().Return(now)
	msg.EXPECT().SetReceiveTime(now)
	// and the message should be called to unmarshal the frame
	msg.EXPECT().UnmarshalFrame(frame)
	// and the node should be unlocked
	node.EXPECT().Unlock()
	// when the next receive fails
	rx.EXPECT().Receive().Return(false)
	rx.EXPECT().Err().Return(nil)
	// then the receiver should return
	assert.NilError(t, canrunner.RunMessageReceiver(ctx, rx, node, clock))
}

func TestRunMessageTransmitter_TransmitEventMessage(t *testing.T) {
	t.Skip() // TODO: fix deadlock flakynes.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tx := mockcanrunner.NewMockFrameTransmitter(ctrl)
	node := mockcanrunner.NewMockNode(ctrl)
	msg := mockcanrunner.NewMockTransmittedMessage(ctrl)
	clock := mockclock.NewMockClock(ctrl)
	desc := &descriptor.Message{
		Name:     "TestMessage",
		SendType: descriptor.SendTypeEvent,
	}
	transmitEventChan := make(chan struct{})
	wakeUpChan := make(chan struct{})
	msg.EXPECT().Descriptor().AnyTimes().Return(desc)
	msg.EXPECT().TransmitEventChan().Return(transmitEventChan)
	msg.EXPECT().WakeUpChan().Return(wakeUpChan)
	// given a running transmitter
	ctx, cancel := context.WithCancel(context.Background())
	var g errgroup.Group
	g.Go(func() error {
		return canrunner.RunMessageTransmitter(ctx, tx, node, msg, clock)
	})
	// then message should be queried for if it has cyclic transmission enabled
	node.EXPECT().Lock()
	msg.EXPECT().IsCyclicTransmissionEnabled()
	node.EXPECT().Unlock()
	// then the node should be locked
	node.EXPECT().Lock()
	// and the time should be queried
	now := time.Unix(0, 1)
	clock.EXPECT().Now().Return(now)
	// and the transmit hook should be queried with the same context
	hook := func(c context.Context) error {
		assert.Equal(t, ctx, c)
		return nil
	}
	msg.EXPECT().BeforeTransmitHook().Return(hook)
	// and the message should be marshaled to a CAN frame
	frame := can.Frame{ID: 42}
	// and the transmit time should be set
	msg.EXPECT().SetTransmitTime(now)
	// and the node should be unlocked
	node.EXPECT().Unlock()
	node.EXPECT().Lock()
	msg.EXPECT().Frame().Return(frame)
	node.EXPECT().Unlock()
	// and the CAN frame should be transmitted
	tx.EXPECT().TransmitFrame(gomock.Any(), frame)
	// when the transmitter receives a transmit event
	transmitEventChan <- struct{}{}
	cancel()
	assert.NilError(t, g.Wait())
}
