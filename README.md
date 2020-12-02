# :electric_plug: CAN Go [![GoDoc][doc-img]][doc]

CAN toolkit for Go programmers.

[doc-img]: https://godoc.org/go.einride.tech/can?status.svg
[doc]: https://godoc.org/go.einride.tech/can


can-go makes use of the linux socketcan abstraction to talk to a CAN bus (See [socketcan](https://www.kernel.org/doc/Documentation/networking/can.txt) for more details)

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

It is possible to generate Go code from a `.dbc` file

```
$ go run go.einride.tech/can/cmd/cantool generate <dbc file root folder> <output folder>
```

In order to generate Go code that makes sense, we currently perform some validations when
parsing the dbc file so there may need to be some changes on the dbc file to make it work

After generating go code we can easy marshal a message to a frame

```go
// import etruckcan "github.com/myproject/myrepo/gen"

auxMsg := etruckcan.NewAuxiliary().SetHeadLights(etruckcan.Auxiliary_HeadLights_LowBeam)
frame := auxMsg.Frame()
```

Or unmarshal a frame to a message

```go
// import etruckcan "github.com/myproject/myrepo/gen"

// Error handling omitted for simplicity
_ := recv.Receive()
frame := recv.Frame()

var auxMsg *etruckcan.Auxiliary
_ = auxMsg.UnmarshalFrame(frame)

```
