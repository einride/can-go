package can

// Message is anything that can marshal and unmarshal itself to/from a CAN frame.
type Message interface {
	FrameMarshaler
	FrameUnmarshaler
}

// FrameMarshaler can marshal itself to a CAN frame.
type FrameMarshaler interface {
	MarshalFrame() (Frame, error)
}

// FrameUnmarshaler can unmarshal itself from a CAN frame.
type FrameUnmarshaler interface {
	UnmarshalFrame(Frame) error
}
