package mocks

//go:generate mockgen -destination gen/mockclock/mocks.go -package mockclock go.einride.tech/can/internal/clock Clock,Ticker
//go:generate mockgen -destination gen/mocksocketcan/mocks.go -package mocksocketcan -source ../../pkg/socketcan/fileconn.go
//go:generate mockgen -destination gen/mockcanrunner/mocks.go -package mockcanrunner go.einride.tech/can/pkg/canrunner Node,TransmittedMessage,ReceivedMessage,FrameTransmitter,FrameReceiver
