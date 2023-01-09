//go:build !linux || !go1.18

package candevice

import (
	"fmt"
	"runtime"
)

type Device struct{}

func New(_ string) (*Device, error) {
	return nil, fmt.Errorf("candevice is not supported on OS %s and runtime %s", runtime.GOOS, runtime.Version())
}
