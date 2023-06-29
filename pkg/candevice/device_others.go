//go:build !linux || !go1.18

package candevice

import (
	"fmt"
	"runtime"
)

type NotSupportedError struct{}

func (e NotSupportedError) Error() string {
	return fmt.Sprintf("candevice is not supported on OS %s and runtime %s", runtime.GOOS, runtime.Version())
}

type Device struct{}

func New(_ string) (*Device, error) {
	return nil, NotSupportedError{}
}

func (d *Device) IsUp() (bool, error) {
	return false, NotSupportedError{}
}

func (d *Device) SetUp() error {
	return NotSupportedError{}
}

func (d *Device) SetDown() error {
	return NotSupportedError{}
}

func (d *Device) Bitrate() (uint32, error) {
	return 0, NotSupportedError{}
}

func (d *Device) SetBitrate(_ uint32) error {
	return NotSupportedError{}
}

func (d *Device) SetListenOnlyMode(mode bool) error {
	return NotSupportedError{}
}

type Info struct{}

func (d *Device) Info() (Info, error) {
	return Info{}, NotSupportedError{}
}
