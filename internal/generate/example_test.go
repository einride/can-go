package generate

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/blueinnovationsgroup/can-go"
	"github.com/blueinnovationsgroup/can-go/pkg/generated"
	"github.com/blueinnovationsgroup/can-go/pkg/socketcan"
	examplecan "github.com/blueinnovationsgroup/can-go/testdata/gen/go/example"
	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"
)

func TestExampleDatabase_MarshalUnmarshal(t *testing.T) {
	for _, tt := range []struct {
		name string
		m    can.Message
		f    can.Frame
	}{
		{
			name: "IODebug",
			m: examplecan.NewIODebug().
				SetTestUnsigned(5).
				SetTestEnum(examplecan.IODebug_TestEnum_Two).
				SetTestSigned(-42).
				SetTestFloat(61.5).
				SetTestBoolEnum(examplecan.IODebug_TestBoolEnum_One).
				SetRawTestScaledEnum(examplecan.IODebug_TestScaledEnum_Four),
			f: can.Frame{
				ID:     500,
				Length: 6,
				Data:   can.Data{5, 2, 214, 123, 1, 2},
			},
		},

		{
			name: "MotorStatus1",
			m: examplecan.NewMotorStatus().
				SetSpeedKph(0.423).
				SetWheelError(true),
			f: can.Frame{
				ID:     400,
				Length: 3,
				Data:   can.Data{0x1, 0xa7, 0x1},
			},
		},

		{
			name: "MotorStatus2",
			m: examplecan.NewMotorStatus().
				SetSpeedKph(12),
			f: can.Frame{
				ID:     400,
				Length: 3,
				Data:   can.Data{0x00, 0xe0, 0x2e},
			},
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f, err := tt.m.MarshalFrame()
			assert.NilError(t, err)
			assert.Equal(t, tt.f, f)
			// allocate new message of same type as tt.m
			msg := reflect.New(reflect.ValueOf(tt.m).Elem().Type()).Interface().(generated.Message)
			assert.NilError(t, msg.UnmarshalFrame(f))
			assert.Assert(t, reflect.DeepEqual(tt.m, msg))
		})
	}
}

func TestExampleDatabase_UnmarshalFrame_Error(t *testing.T) {
	for _, tt := range []struct {
		name string
		f    can.Frame
		m    generated.Message
		err  string
	}{
		{
			name: "wrong ID",
			f:    can.Frame{ID: 11, Length: 8},
			m:    examplecan.NewSensorSonars(),
			err:  "unmarshal SensorSonars: expects ID 200 (got 00B#0000000000000000 with ID 11)",
		},
		{
			name: "wrong length",
			f:    can.Frame{ID: 200, Length: 4},
			m:    examplecan.NewSensorSonars(),
			err:  "unmarshal SensorSonars: expects length 8 (got 0C8#00000000 with length 4)",
		},
		{
			name: "remote frame",
			f:    can.Frame{ID: 200, Length: 8, IsRemote: true},
			m:    examplecan.NewSensorSonars(),
			err:  "unmarshal SensorSonars: expects non-remote frame (got remote frame 0C8#R8)",
		},
		{
			name: "extended ID",
			f:    can.Frame{ID: 200, Length: 8, IsExtended: true},
			m:    examplecan.NewSensorSonars(),
			err:  "unmarshal SensorSonars: expects standard ID (got 000000C8#0000000000000000 with extended ID)",
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.err, tt.m.UnmarshalFrame(tt.f).Error())
		})
	}
}

func TestExampleDatabase_TestEnum_String(t *testing.T) {
	assert.Equal(t, "One", examplecan.IODebug_TestEnum_One.String())
	assert.Equal(t, "Two", examplecan.IODebug_TestEnum_Two.String())
	assert.Equal(t, "IODebug_TestEnum(3)", examplecan.IODebug_TestEnum(3).String())
}

func TestExampleDatabase_Message_String(t *testing.T) {
	const expected = "{WheelError: true, SpeedKph: 42km/h}"
	msg := examplecan.NewMotorStatus().
		SetSpeedKph(42).
		SetWheelError(true)
	assert.Equal(t, expected, msg.String())
	assert.Equal(t, expected, fmt.Sprintf("%v", msg))
}

func TestExampleDatabase_OutOfBoundsValue(t *testing.T) {
	const expected = examplecan.IODebug_TestEnum(63)
	actual := examplecan.NewIODebug().SetTestEnum(255).TestEnum()
	assert.Equal(t, expected, actual)
}

func TestExampleDatabase_MultiplexedSignals(t *testing.T) {
	// Given a message with multiplexed signals
	msg := examplecan.NewSensorSonars().
		SetErrCount(1).
		SetMux(1).
		SetLeft(20).
		SetMiddle(30).
		SetRight(40).
		SetRear(50).
		SetNoFiltLeft(60).
		SetNoFiltMiddle(70).
		SetNoFiltRight(80).
		SetNoFiltRear(90)
	for _, tt := range []struct {
		expectedMux          uint8
		expectedErrCount     uint16
		expectedLeft         float64
		expectedMiddle       float64
		expectedRight        float64
		expectedRear         float64
		expectedNoFiltLeft   float64
		expectedNoFiltMiddle float64
		expectedNoFiltRight  float64
		expectedNoFiltRear   float64
	}{
		{
			expectedMux:          0,
			expectedErrCount:     1,
			expectedLeft:         20,
			expectedMiddle:       30,
			expectedRight:        40,
			expectedRear:         50,
			expectedNoFiltLeft:   0,
			expectedNoFiltMiddle: 0,
			expectedNoFiltRight:  0,
			expectedNoFiltRear:   0,
		},
		{
			expectedMux:          1,
			expectedErrCount:     1,
			expectedLeft:         0,
			expectedMiddle:       0,
			expectedRight:        0,
			expectedRear:         0,
			expectedNoFiltLeft:   60,
			expectedNoFiltMiddle: 70,
			expectedNoFiltRight:  80,
			expectedNoFiltRear:   90,
		},
	} {
		tt := tt
		t.Run(fmt.Sprintf("mux=%v", tt.expectedMux), func(t *testing.T) {
			unmarshal1 := examplecan.NewSensorSonars()
			// When the multiplexer signal is 0 and we marshal the message
			// to a CAN frame
			msg.SetMux(tt.expectedMux)
			f1, err := msg.MarshalFrame()
			assert.NilError(t, err)
			// When we unmarshal the CAN frame back to a message
			assert.NilError(t, unmarshal1.UnmarshalFrame(f1))
			// Then only the multiplexed signals with multiplexer value 0
			// should be unmarshaled
			assert.Equal(t, tt.expectedMux, unmarshal1.Mux(), "Mux")
			assert.Equal(t, tt.expectedErrCount, unmarshal1.ErrCount(), "ErrCount")
			assert.Equal(t, tt.expectedLeft, unmarshal1.Left(), "Left")
			assert.Equal(t, tt.expectedMiddle, unmarshal1.Middle(), "Middle")
			assert.Equal(t, tt.expectedRight, unmarshal1.Right(), "Right")
			assert.Equal(t, tt.expectedRear, unmarshal1.Rear(), "Rear")
			assert.Equal(t, tt.expectedNoFiltLeft, unmarshal1.NoFiltLeft(), "NoFiltLeft")
			assert.Equal(t, tt.expectedNoFiltMiddle, unmarshal1.NoFiltMiddle(), "NoFiltMiddle")
			assert.Equal(t, tt.expectedNoFiltRight, unmarshal1.NoFiltRight(), "NoFiltRight")
			assert.Equal(t, tt.expectedNoFiltRear, unmarshal1.NoFiltRear(), "NoFiltRear")
		})
	}
}

func TestExampleDatabase_CopyFrom(t *testing.T) {
	// Given: an original message
	from := examplecan.NewIODebug().
		SetRawTestScaledEnum(examplecan.IODebug_TestScaledEnum_Four).
		SetTestBoolEnum(true).
		SetTestFloat(0.1).
		SetTestSigned(-10).
		SetTestUnsigned(10)
	// When: another message copies from the original message
	to := examplecan.NewIODebug().CopyFrom(from)
	// Then:
	// all fields should be equal...
	assert.Equal(t, from.String(), to.String())
	assert.Equal(t, from.TestScaledEnum(), to.TestScaledEnum())
	assert.Equal(t, from.TestBoolEnum(), to.TestBoolEnum())
	assert.Equal(t, from.TestFloat(), to.TestFloat())
	assert.Equal(t, from.TestSigned(), to.TestSigned())
	assert.Equal(t, from.TestUnsigned(), to.TestUnsigned())
	// ...and changes to the original should not affect the new message
	from.SetTestUnsigned(100)
	assert.Equal(t, uint8(10), to.TestUnsigned())
}

func TestExample_Nodes(t *testing.T) {
	const testTimeout = 2 * time.Second
	requireVCAN0(t)
	// given a DRIVER node and a MOTOR node
	motor := examplecan.NewMOTOR("can", "vcan0")
	driver := examplecan.NewDRIVER("can", "vcan0")
	// when starting them
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return motor.Run(ctx)
	})
	g.Go(func() error {
		return driver.Run(ctx)
	})
	// and the MOTOR node is configured to send a speed report
	const expectedSpeedKph = 42
	motor.Lock()
	motor.Tx().MotorStatus().SetSpeedKph(expectedSpeedKph)
	motor.Tx().MotorStatus().SetCyclicTransmissionEnabled(true)
	motor.Unlock()
	// and the DRIVER node is configured to send a steering command
	const expectedSteer = -4
	driver.Lock()
	driver.Tx().MotorCommand().SetSteer(expectedSteer)
	driver.Tx().MotorCommand().SetCyclicTransmissionEnabled(true)
	driver.Unlock()
	// and the MOTOR node is listening for the steering command
	expectedSteerReceivedChan := make(chan struct{})
	motor.Lock()
	motor.Rx().MotorCommand().SetAfterReceiveHook(func(context.Context) error {
		motor.Lock()
		if motor.Rx().MotorCommand().Steer() == expectedSteer {
			close(expectedSteerReceivedChan)
			motor.Rx().MotorCommand().SetAfterReceiveHook(func(context.Context) error { return nil })
		}
		motor.Unlock()
		return nil
	})
	motor.Unlock()
	// and the DRIVER node is listening for the speed report
	expectedSpeedReceivedChan := make(chan struct{})
	driver.Lock()
	driver.Rx().MotorStatus().SetAfterReceiveHook(func(context.Context) error {
		driver.Lock()
		if driver.Rx().MotorStatus().SpeedKph() == expectedSpeedKph {
			close(expectedSpeedReceivedChan)
			driver.Rx().MotorStatus().SetAfterReceiveHook(func(context.Context) error { return nil })
		}
		driver.Unlock()
		return nil
	})
	driver.Unlock()
	// then the steer command transmitted by DRIVER should be received by MOTOR
	select {
	case <-expectedSteerReceivedChan:
	case <-ctx.Done():
		t.Fatalf("expected steer not received: %v", expectedSteer)
	}
	// and the speed report transmitted by MOTOR should be received by DRIVER
	select {
	case <-expectedSpeedReceivedChan:
	case <-ctx.Done():
		t.Fatalf("expected speed not received: %v", expectedSpeedKph)
	}
	cancel()
	assert.NilError(t, g.Wait())
}

func TestExample_Node_NoEmptyMessages(t *testing.T) {
	const testTimeout = 2 * time.Second
	requireVCAN0(t)
	// given a DRIVER node and a MOTOR node
	motor := examplecan.NewMOTOR("can", "vcan0")
	// when starting them
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	handler := func(ctx context.Context) error {
		motor.Lock()
		motor.Tx().MotorStatus().SetSpeedKph(100).SetWheelError(true)
		motor.Unlock()
		return nil
	}
	motor.Tx().MotorStatus().SetBeforeTransmitHook(handler)
	motor.Tx().MotorStatus().SetCyclicTransmissionEnabled(true)
	c, err := socketcan.Dial("can", "vcan0")
	r := socketcan.NewReceiver(c)
	assert.NilError(t, err)
	g := errgroup.Group{}
	g.Go(func() error {
		return motor.Run(ctx)
	})
	assert.Assert(t, r.Receive())
	assert.Equal(t, examplecan.NewMotorStatus().SetSpeedKph(100).SetWheelError(true).Frame(), r.Frame())
	cancel()
	assert.NilError(t, g.Wait())
}

func requireVCAN0(t *testing.T) {
	t.Helper()
	if _, err := net.InterfaceByName("vcan0"); err != nil {
		t.Skip("interface vcan0 does not exist")
	}
}
