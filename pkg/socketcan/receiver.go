package socketcan

import (
	"bufio"
	"io"

	"go.einride.tech/can"
)

type ReceiverOption func(*receiverOpts)

type receiverOpts struct {
	frameInterceptor FrameInterceptor
}

type Receiver struct {
	opts  receiverOpts
	rc    io.ReadCloser
	sc    *bufio.Scanner
	frame frame
}

func NewReceiver(rc io.ReadCloser, opt ...ReceiverOption) *Receiver {
	opts := receiverOpts{}
	for _, f := range opt {
		f(&opts)
	}
	sc := bufio.NewScanner(rc)
	sc.Split(scanFrames)
	return &Receiver{
		rc:   rc,
		opts: opts,
		sc:   sc,
	}
}

func scanFrames(data []byte, _ bool) (int, []byte, error) {
	if len(data) < lengthOfFrame {
		// not enough data for a full frame
		return 0, nil, nil
	}
	return lengthOfFrame, data[0:lengthOfFrame], nil
}

func (r *Receiver) Receive() bool {
	ok := r.sc.Scan()
	r.frame = frame{}
	if ok {
		r.frame.unmarshalBinary(r.sc.Bytes())
		if r.opts.frameInterceptor != nil {
			r.opts.frameInterceptor(r.frame.decodeFrame())
		}
	}
	return ok
}

func (r *Receiver) HasErrorFrame() bool {
	return r.frame.isError()
}

func (r *Receiver) Frame() can.Frame {
	return r.frame.decodeFrame()
}

func (r *Receiver) ErrorFrame() ErrorFrame {
	return r.frame.decodeErrorFrame()
}

func (r *Receiver) Err() error {
	return r.sc.Err()
}

func (r *Receiver) Close() error {
	return r.rc.Close()
}

// ReceiverFrameInterceptor returns a ReceiverOption that sets the FrameInterceptor for the
// receiver. Only one frame interceptor can be installed.
func ReceiverFrameInterceptor(i FrameInterceptor) ReceiverOption {
	return func(o *receiverOpts) {
		o.frameInterceptor = i
	}
}
