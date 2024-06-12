# :electric_plug: go.einride.tech/can

[![PkgGoDev](https://pkg.go.dev/badge/go.einride.tech/can)](https://pkg.go.dev/go.einride.tech/can)
[![GoReportCard](https://goreportcard.com/badge/go.einride.tech/can)](https://goreportcard.com/report/go.einride.tech/can)
[![Codecov](https://codecov.io/gh/einride/can-go/branch/master/graph/badge.svg)](https://codecov.io/gh/einride/can-go)

CAN toolkit for Go programmers.

can-go makes use of the Linux SocketCAN abstraction for CAN communication. (See
the [SocketCAN](https://www.kernel.org/doc/Documentation/networking/can.txt)
documentation for more details).

## Installation

```
go get -u go.einride.tech/can
```

## Examples

### Setting up a CAN interface

```go

import "go.einride.tech/can/pkg/candevice"

func main() {
	// Error handling omitted to keep example simple
	d, _ := candevice.New("can0")
	_ := d.SetBitrate(250000)
	_ := d.SetUp()
	defer d.SetDown()
}
```

### Decoding CAN messages

Decoding CAN messages from byte arrays can be done using `can.Payload`  

```go
func main() {
	// DBC file
	var dbcFile = []byte(`
	VERSION ""
    NS_ :
    BS_:
    BU_: DBG DRIVER IO MOTOR SENSOR
    
    BO_ 400 MOTOR_STATUS: 3 MOTOR
	  SG_ MOTOR_STATUS_wheel_error : 0|1@1+ (1,0) [0|0] "" DRIVER,IO
	  SG_ MOTOR_STATUS_speed_kph : 8|16@1+ (0.001,0) [0|0] "kph" DRIVER,IO

	BO_ 200 SENSOR_SONARS: 8 SENSOR
	  SG_ SENSOR_SONARS_mux M : 0|4@1+ (1,0) [0|0] "" DRIVER,IO
	  SG_ SENSOR_SONARS_err_count : 4|12@1+ (1,0) [0|0] "" DRIVER,IO
	  SG_ SENSOR_SONARS_left m0 : 16|12@1+ (0.1,0) [0|0] "" DRIVER,IO
	  SG_ SENSOR_SONARS_middle m0 : 28|12@1+ (0.1,0) [0|0] "" DRIVER,IO
	  SG_ SENSOR_SONARS_right m0 : 40|12@1+ (0.1,0) [0|0] "" DRIVER,IO
	  SG_ SENSOR_SONARS_rear m0 : 52|12@1+ (0.1,0) [0|0] "" DRIVER,IO
	  SG_ SENSOR_SONARS_no_filt_left m1 : 16|12@1+ (0.1,0) [0|0] "" DBG
	  SG_ SENSOR_SONARS_no_filt_middle m1 : 28|12@1+ (0.1,0) [0|0] "" DBG
	  SG_ SENSOR_SONARS_no_filt_right m1 : 40|12@1+ (0.1,0) [0|0] "" DBG
	  SG_ SENSOR_SONARS_no_filt_rear m1 : 52|12@1+ (0.1,0) [0|0] "" DBG
    `)

	// Create payload from hex string
	byteStringHex := "004faf"
	p, _ := can.PayloadFromHex(byteStringHex)

	// Load example dbc file
	c, _ := generate.Compile("test.dbc", dbcFile)
	db := *c.Database

	// Decode message frame ID 400
	message, _ := db.Message(uint32(400))
	decodedSignals := message.Decode(&p)
	for _, signal := range decodedSignals {
		fmt.Printf("Signal: %s, Value: %f, Description: %s\n", signal.Signal.Name, signal.Value, signal.Description)
	}
}
```

```
Signal: MOTOR_STATUS_wheel_error, Value: 0.000000, Description: 
Signal: MOTOR_STATUS_speed_kph, Value: 44.879000, Description: 
```  

#### Multiplexed Signals  

```go
func main() {
	// DBC file
	var dbcFile = []byte(`
	VERSION ""
    NS_ :
    BS_:
    BU_: DBG DRIVER IO MOTOR SENSOR
    
    BO_ 400 MOTOR_STATUS: 3 MOTOR
	  SG_ MOTOR_STATUS_wheel_error : 0|1@1+ (1,0) [0|0] "" DRIVER,IO
	  SG_ MOTOR_STATUS_speed_kph : 8|16@1+ (0.001,0) [0|0] "kph" DRIVER,IO

	BO_ 200 SENSOR_SONARS: 8 SENSOR
	  SG_ SENSOR_SONARS_mux M : 0|4@1+ (1,0) [0|0] "" DRIVER,IO
	  SG_ SENSOR_SONARS_err_count : 4|12@1+ (1,0) [0|0] "" DRIVER,IO
	  SG_ SENSOR_SONARS_left m0 : 16|12@1+ (0.1,0) [0|0] "" DRIVER,IO
	  SG_ SENSOR_SONARS_middle m0 : 28|12@1+ (0.1,0) [0|0] "" DRIVER,IO
	  SG_ SENSOR_SONARS_right m0 : 40|12@1+ (0.1,0) [0|0] "" DRIVER,IO
	  SG_ SENSOR_SONARS_rear m0 : 52|12@1+ (0.1,0) [0|0] "" DRIVER,IO
	  SG_ SENSOR_SONARS_no_filt_left m1 : 16|12@1+ (0.1,0) [0|0] "" DBG
	  SG_ SENSOR_SONARS_no_filt_middle m1 : 28|12@1+ (0.1,0) [0|0] "" DBG
	  SG_ SENSOR_SONARS_no_filt_right m1 : 40|12@1+ (0.1,0) [0|0] "" DBG
	  SG_ SENSOR_SONARS_no_filt_rear m1 : 52|12@1+ (0.1,0) [0|0] "" DBG
    `)

	// Create payload from hex string
	byteStringHex := "01af79f4aa3b459f"
	p, _ := can.PayloadFromHex(byteStringHex)

	// Load example dbc file
	c, _ := generate.Compile("test.dbc", dbcFile)
	db := *c.Database

	// Decode message frame ID 200
	message, _ := db.Message(uint32(200))
	decodedSignals := message.Decode(&p)
	for _, signal := range decodedSignals {
		fmt.Printf("Signal: %s, Value: %f, Description: %s\n", signal.Signal.Name, signal.Value, signal.Description)
	}
}
```  

```
Signal: SENSOR_SONARS_mux, Value: 1.000000, Description: 
Signal: SENSOR_SONARS_err_count, Value: 2800.000000, Description: 
Signal: SENSOR_SONARS_no_filt_left, Value: 114.500000, Description: 
Signal: SENSOR_SONARS_no_filt_middle, Value: 273.500000, Description: 
Signal: SENSOR_SONARS_no_filt_right, Value: 133.900000, Description: 
Signal: SENSOR_SONARS_no_filt_rear, Value: 254.800000, Description: 
```

### Receiving CAN frames

Receiving CAN frames from a socketcan interface.

```go
import "go.einride.tech/can/pkg/socketcan"

func main() {
	// Error handling omitted to keep example simple
	conn, _ := socketcan.DialContext(context.Background(), "can", "can0")

	recv := socketcan.NewReceiver(conn)
	for recv.Receive() {
		frame := recv.Frame()
		fmt.Println(frame.String())
	}
}
```

### Sending CAN frames/messages

Sending CAN frames to a socketcan interface.

```go
import "go.einride.tech/can/pkg/socketcan"

func main() {
	// Error handling omitted to keep example simple

	conn, _ := socketcan.DialContext(context.Background(), "can", "can0")

	frame := can.Frame{}
	tx := socketcan.NewTransmitter(conn)
	_ = tx.TransmitFrame(context.Background(), frame)
}
```

### Generating Go code from a DBC file

It is possible to generate Go code from a `.dbc` file.

```
$ go run go.einride.tech/can/cmd/cantool generate <dbc file root folder> <output folder>
```

In order to generate Go code that makes sense, we currently perform some
validations when parsing the DBC file so there may need to be some changes on
the DBC file to make it work

After generating Go code we can marshal a message to a frame:

```go
// import etruckcan "github.com/myproject/myrepo/gen"

auxMsg := etruckcan.NewAuxiliary().SetHeadLights(etruckcan.Auxiliary_HeadLights_LowBeam)
frame := auxMsg.Frame()
```

Or unmarshal a frame to a message:

```go
// import etruckcan "github.com/myproject/myrepo/gen"

// Error handling omitted for simplicity
_ := recv.Receive()
frame := recv.Frame()

var auxMsg *etruckcan.Auxiliary
_ = auxMsg.UnmarshalFrame(frame)

```

## Running integration tests

Building the tests:

```shell
$ make build-integration-tests
```

Built tests are placed in build/tests.

The candevice test requires access to physical HW, so run it using sudo.
Example:

```shell
$ sudo ./build/tests/candevice.test
> PASS
```
