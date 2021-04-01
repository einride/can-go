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
    // DBC file
    var dbcFile = []byte(`
    VERSION ""
    NS_ :
    BS_:
    BU_: DBG DRIVER IO MOTOR SENSOR
    
    BO_ 1530 DisconnectState: 14 MOTOR
    SG_ LockCountRearRight : 91|20@0+ (1,0) [0|1048575] ""  IO
    SG_ DisconnectStateRearRight : 95|4@0+ (1,0) [0|5] ""  IO
    SG_ CurrentRearRight : 79|16@0+ (1,0) [0|65535] ""  IO
    SG_ DisconnectStateRearRightTarget : 64|1@0+ (1,0) [0|1] ""  IO
    SG_ TargetSpeedRearRight : 63|15@0+ (0.125,-2048) [-2048|2047.875] "rad/s"  IO
    SG_ LockCountRearLeft : 35|20@0+ (1,0) [0|1048575] ""  IO
    SG_ DisconnectStateRearLeft : 39|4@0+ (1,0) [0|5] ""  IO
    SG_ CurrentRearLeft : 23|16@0+ (1,0) [0|65535] ""  IO
    SG_ DisconnectStateRearLeftTarget : 8|1@0+ (1,0) [0|1] ""  IO
    SG_ TargetSpeedRearLeft : 7|15@0+ (0.125,-2048) [-2048|2047.875] "rad/s"  IO
    
    BA_DEF_ "BusType" STRING ;
    BA_DEF_ BO_  "GenMsgSendType" ENUM  "None","Cyclic","OnEvent";
    BA_DEF_ BO_ "GenMsgCycleTime" INT 0 0;
    BA_DEF_ SG_  "FieldType" STRING ;
    BA_DEF_ SG_  "GenSigStartValue" INT 0 10000;
    BA_DEF_DEF_ "BusType" "CAN";
    BA_DEF_DEF_ "FieldType" "";
    BA_DEF_DEF_ "GenMsgCycleTime" 0;
    BA_DEF_DEF_ "GenSigStartValue" 0;

    BA_ "GenMsgSendType" BO_ 1 0;
    BA_ "GenMsgSendType" BO_ 100 1;
    BA_ "GenMsgCycleTime" BO_ 100 1000;
    BA_ "GenMsgSendType" BO_ 101 1;
    BA_ "GenMsgCycleTime" BO_ 101 100;
    BA_ "GenMsgSendType" BO_ 200 1;
    BA_ "GenMsgCycleTime" BO_ 200 100;
    BA_ "GenMsgSendType" BO_ 400 1;
    BA_ "GenMsgCycleTime" BO_ 400 100;
    BA_ "GenMsgSendType" BO_ 500 2;
    BA_ "FieldType" SG_ 100 Command "Command";
    BA_ "FieldType" SG_ 500 TestEnum "TestEnum";
    BA_ "GenSigStartValue" SG_ 500 TestEnum 2;
    
    VAL_ 1530 DisconnectStateRearRight 0 "Undefined" 1 "Locked" 2 "Unlocked" 3 "Locking" 4 "Unlocking" 5 "Faulted" ;
    VAL_ 1530 DisconnectStateRearLeft 0 "Undefined" 1 "Locked" 2 "Unlocked" 3 "Locking" 4 "Unlocking" 5 "Faulted" ;
    `)

    // Create payload from hex string
    byteStringHex := "8000000420061880000005200600"
    p, _ := can.PayloadFromHex(byteStringHex)

    // Load example dbc file
    c, _ := generate.Compile("test.dbc", dbcFile)
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
