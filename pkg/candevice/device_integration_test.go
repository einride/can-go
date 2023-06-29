//go:build integration

package candevice

import (
	"fmt"
	"testing"
)

const (
	bitrate125K = 125000
	bitrate250K = 250000
)

func TestSetBitrate(t *testing.T) {
	d, err := New("can0")
	if err != nil {
		t.Fatal("couldn't set up device:", err)
	}
	defer d.SetDown()

	if err := setBitrate(d, bitrate125K); err != nil {
		t.Fatal(err)
	}
	if err := setBitrate(d, bitrate250K); err != nil {
		t.Fatal(err)
	}

	// Set bitrate on device which is up
	if err := d.SetUp(); err != nil {
		t.Fatal(err)
	}
	if err := setBitrate(d, bitrate125K); err == nil {
		t.Fatal("setting bitrate on device which is up succeeded")
	}
	if err := d.SetDown(); err != nil {
		t.Fatal(err)
	}

	// Invalid bitrate
	if err := setBitrate(d, 0); err == nil {
		t.Fatal("setting invalid bitrate succeeded")
	}
}

func TestSetUpDown(t *testing.T) {
	d, err := New("can0")
	if err != nil {
		t.Fatal("couldn't set up device:", err)
	}
	defer d.SetDown()

	if err := d.SetBitrate(bitrate125K); err != nil {
		t.Fatal(err)
	}

	// Set up twice and set down twice. This checks that calling it twice has no effect
	if err := setUp(d); err != nil {
		t.Fatal(err)
	}
	if err := setUp(d); err != nil {
		t.Fatal(err)
	}
	if err := setDown(d); err != nil {
		t.Fatal(err)
	}
	if err := setDown(d); err != nil {
		t.Fatal(err)
	}
}

func TestSetListenOnlyMode(t *testing.T) {
	d, err := New("can0")
	if err != nil {
		t.Fatal("couldn't set up device:", err)
	}
	defer d.SetDown()

	if err := d.SetListenOnlyMode(true); err != nil {
		t.Fatal(err)
	}

	// Set ListenOnly mode on device which is up
	if err := d.SetUp(); err != nil {
		t.Fatal(err)
	}
	if err := d.SetListenOnlyMode(false); err == nil {
		t.Fatal("setting ListenOnly mode on device which is up succeeded")
	}
	if err := d.SetDown(); err != nil {
		t.Fatal(err)
	}

	if err := d.SetListenOnlyMode(false); err != nil {
		t.Fatal(err)
	}
}

func setBitrate(d *Device, bitrate uint32) error {
	if err := d.SetBitrate(bitrate); err != nil {
		return err
	}
	if err := setUp(d); err != nil {
		return err
	}
	actualBitrate, err := d.Bitrate()
	if err != nil {
		return err
	}
	if err := setDown(d); err != nil {
		return err
	}
	if actualBitrate != bitrate {
		return fmt.Errorf("expected bitrate: %d, actual: %d", bitrate, bitrate)
	}
	return nil
}

func setUp(d *Device) error {
	if err := d.SetUp(); err != nil {
		return err
	}

	isUp, err := d.IsUp()
	if err != nil {
		return err
	}
	if !isUp {
		return fmt.Errorf("device not up after calling SetUp()")
	}
	return nil
}

func setDown(d *Device) error {
	if err := d.SetDown(); err != nil {
		return err
	}
	isUp, err := d.IsUp()
	if err != nil {
		return err
	}
	if isUp {
		return fmt.Errorf("device not down after calling SetDown()")
	}
	return nil
}
