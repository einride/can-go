package socketcan

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.einride.tech/can/internal/gen/mock/mocksocketcan"
)

func TestUnwrapPathError(t *testing.T) {
	innerErr := fmt.Errorf("inner error")
	for _, tt := range []struct {
		msg      string
		err      error
		expected error
	}{
		{
			msg:      "no path error",
			err:      innerErr,
			expected: innerErr,
		},
		{
			msg:      "single path error",
			err:      &os.PathError{Op: "read", Err: innerErr},
			expected: innerErr,
		},
		{
			msg:      "double path error",
			err:      &os.PathError{Op: "read", Err: &os.PathError{Op: "read", Err: innerErr}},
			expected: &os.PathError{Op: "read", Err: innerErr},
		},
	} {
		tt := tt
		t.Run(tt.msg, func(t *testing.T) {
			require.Equal(t, tt.expected, unwrapPathError(tt.err))
		})
	}
}

func TestFileConn_ReadWrite(t *testing.T) {
	for _, tt := range []struct {
		op     string
		fn     func(file, []byte) (int, error)
		mockFn func(*mocksocketcan.MockfileMockRecorder, interface{}) *gomock.Call
	}{
		{
			op:     "read",
			fn:     file.Read,
			mockFn: (*mocksocketcan.MockfileMockRecorder).Read,
		},
		{
			op:     "write",
			fn:     file.Write,
			mockFn: (*mocksocketcan.MockfileMockRecorder).Write,
		},
	} {
		tt := tt
		t.Run(tt.op, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := mocksocketcan.NewMockfile(ctrl)
			fc := &fileConn{f: f, net: "can", ra: &canRawAddr{device: "can0"}}
			t.Run("no error", func(t *testing.T) {
				var data []byte
				tt.mockFn(f.EXPECT(), data).Return(42, nil)
				n, err := tt.fn(fc, data)
				require.Equal(t, 42, n)
				require.NoError(t, err)
			})
			t.Run("error", func(t *testing.T) {
				var data []byte
				cause := fmt.Errorf("boom")
				tt.mockFn(f.EXPECT(), data).Return(0, &os.PathError{Err: cause})
				n, err := tt.fn(fc, data)
				require.Equal(t, 0, n)
				require.Error(t, &net.OpError{Op: tt.op, Net: fc.net, Addr: fc.RemoteAddr(), Err: err})
			})
		})
	}
}

func TestFileConn_Addr(t *testing.T) {
	fc := &fileConn{la: &canRawAddr{device: "can0"}, ra: &canRawAddr{device: "can1"}}
	t.Run("local", func(t *testing.T) {
		require.Equal(t, fc.la, fc.LocalAddr())
	})
	t.Run("remote", func(t *testing.T) {
		require.Equal(t, fc.ra, fc.RemoteAddr())
	})
}

func TestFileConn_SetDeadlines(t *testing.T) {
	for _, tt := range []struct {
		op     string
		fn     func(file, time.Time) error
		mockFn func(*mocksocketcan.MockfileMockRecorder, interface{}) *gomock.Call
	}{
		{
			op:     "set deadline",
			fn:     file.SetDeadline,
			mockFn: (*mocksocketcan.MockfileMockRecorder).SetDeadline,
		},
		{
			op:     "set read deadline",
			fn:     file.SetReadDeadline,
			mockFn: (*mocksocketcan.MockfileMockRecorder).SetReadDeadline,
		},
		{
			op:     "set write deadline",
			fn:     file.SetWriteDeadline,
			mockFn: (*mocksocketcan.MockfileMockRecorder).SetWriteDeadline,
		},
	} {
		tt := tt
		t.Run(tt.op, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := mocksocketcan.NewMockfile(ctrl)
			fc := &fileConn{f: f, net: "can", ra: &canRawAddr{device: "can0"}}
			t.Run("no error", func(t *testing.T) {
				tt.mockFn(f.EXPECT(), time.Unix(0, 1)).Return(nil)
				require.NoError(t, tt.fn(fc, time.Unix(0, 1)))
			})
			t.Run("error", func(t *testing.T) {
				cause := fmt.Errorf("boom")
				tt.mockFn(f.EXPECT(), time.Unix(0, 1)).Return(&os.PathError{Err: cause})
				err := tt.fn(fc, time.Unix(0, 1))
				require.Equal(t, &net.OpError{Op: tt.op, Net: fc.net, Addr: fc.RemoteAddr(), Err: cause}, err)
			})
		})
	}
}

func TestFileConn_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	f := mocksocketcan.NewMockfile(ctrl)
	fc := &fileConn{f: f, net: "can", ra: &canRawAddr{device: "can0"}}
	t.Run("no error", func(t *testing.T) {
		f.EXPECT().Close().Return(nil)
		require.NoError(t, fc.Close())
	})
	t.Run("error", func(t *testing.T) {
		cause := fmt.Errorf("boom")
		f.EXPECT().Close().Return(&os.PathError{Err: cause})
		err := fc.Close()
		require.Equal(t, &net.OpError{Op: "close", Net: fc.net, Addr: fc.RemoteAddr(), Err: cause}, err)
	})
}
