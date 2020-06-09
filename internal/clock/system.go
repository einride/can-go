package clock

import (
	"time"
)

// System returns a Clock implementation that delegate to the time package.
func System() Clock {
	return &systemClock{}
}

type systemClock struct{}

var _ Clock = &systemClock{}

func (c systemClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (c systemClock) NewTicker(d time.Duration) Ticker {
	return &systemTicker{Ticker: *time.NewTicker(d)}
}

func (c systemClock) Now() time.Time {
	return time.Now()
}

type systemTicker struct {
	time.Ticker
}

func (t systemTicker) C() <-chan time.Time {
	return t.Ticker.C
}
