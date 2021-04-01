# :electric_plug: CAN Go

[![PkgGoDev][pkg-badge]][pkg]
[![GoReportCard][report-badge]][report]
[![Codecov][codecov-badge]][codecov]

[pkg-badge]: https://pkg.go.dev/badge/go.einride.tech/can
[pkg]: https://pkg.go.dev/go.einride.tech/can
[report-badge]: https://goreportcard.com/badge/go.einride.tech/can
[report]: https://goreportcard.com/report/go.einride.tech/can
[codecov-badge]: https://codecov.io/gh/einride/can-go/branch/master/graph/badge.svg
[codecov]: https://codecov.io/gh/einride/can-go

CAN toolkit for Go programmers.

can-go makes use of the Linux SocketCAN abstraction for CAN communication.
(See the [SocketCAN][socketcan] documentation for more details).

[socketcan]: https://www.kernel.org/doc/Documentation/networking/can.txt

## Examples

### Decoding CAN messages

Decoding CAN messages from byte arrays can be done using `can.Payload`

```go
func main() {
    // Create payload from hex string
    byteStringHex := "8000000420061880000005200600"
    p, _ := can.PayloadFromHex(byteStringHex)

    // Load example dbc file
    dbcFile := "./testdata/dbc/example/example_payload.dbc"
    input, _ := ioutil.ReadFile(dbcFile)
    c, _ := generate.Compile(dbcFile, input)
    db := *c.Database

    // Decode message frame ID 1530
    message, _ := db.Message(uint32(1530))
    decodedSignals := message.Decode(&p)
    for _, signal := range decodedSignals {
        fmt.Printf("Signal: %s, Value: %f, Description: %s\n", signal.Signal.Name, signal.Value, signal.Description)
    }
}
```

```
Signal: TargetSpeedRearLeft, Value: 0.000000, Description: 
Signal: DisconnectStateRearLeftTarget, Value: 0.000000, Description: 
Signal: CurrentRearLeft, Value: 4.000000, Description: 
Signal: LockCountRearLeft, Value: 1560.000000, Description: 
Signal: DisconnectStateRearLeft, Value: 2.000000, Description: Unlocked
Signal: TargetSpeedRearRight, Value: 0.000000, Description: 
Signal: DisconnectStateRearRightTarget, Value: 0.000000, Description: 
Signal: CurrentRearRight, Value: 5.000000, Description: 
Signal: LockCountRearRight, Value: 1536.000000, Description: 
Signal: DisconnectStateRearRight, Value: 2.000000, Description: Unlocked
```


### Receiving CAN frames

Receiving CAN frames from a socketcan interface.

```go
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
validations when parsing the DBC file so there may need to be some changes
on the DBC file to make it work

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
