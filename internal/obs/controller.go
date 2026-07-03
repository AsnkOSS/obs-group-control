// Package obs wraps goobs to control a group of OBS instances.
package obs

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/andreykaipov/goobs"

	"obs-group-control/internal/config"
)

const (
	verifyTimeout  = 10 * time.Second
	verifyInterval = 300 * time.Millisecond
)

// Controller drives recording on all configured devices.
type Controller struct {
	devices []config.Device
}

// NewController creates a controller for the given devices.
func NewController(devices []config.Device) *Controller {
	return &Controller{devices: devices}
}

// DeviceCount returns the number of configured devices.
func (c *Controller) DeviceCount() int {
	return len(c.devices)
}

// StartAll connects to every device, starts recording, and blocks until
// every instance confirms it is actively recording. Any failure aborts
// with an error so the UI never shows "recording" prematurely.
func (c *Controller) StartAll() error {
	return c.forEach(func(client *goobs.Client, d config.Device) error {
		status, err := client.Record.GetRecordStatus()
		if err != nil {
			return fmt.Errorf("%s: query status: %w", d.Addr, err)
		}
		if !status.OutputActive {
			if _, err := client.Record.StartRecord(); err != nil {
				return fmt.Errorf("%s: start record: %w", d.Addr, err)
			}
		}
		return c.waitRecordState(client, d, true)
	})
}

// StopAll stops recording on every device and waits for confirmation.
func (c *Controller) StopAll() error {
	return c.forEach(func(client *goobs.Client, d config.Device) error {
		status, err := client.Record.GetRecordStatus()
		if err != nil {
			return fmt.Errorf("%s: query status: %w", d.Addr, err)
		}
		if status.OutputActive {
			if _, err := client.Record.StopRecord(); err != nil {
				return fmt.Errorf("%s: stop record: %w", d.Addr, err)
			}
		}
		return c.waitRecordState(client, d, false)
	})
}

// forEach runs op against every device concurrently and joins errors.
func (c *Controller) forEach(op func(*goobs.Client, config.Device) error) error {
	var wg sync.WaitGroup
	errs := make([]error, len(c.devices))

	for i, d := range c.devices {
		wg.Add(1)
		go func(i int, d config.Device) {
			defer wg.Done()
			client, err := connect(d)
			if err != nil {
				errs[i] = fmt.Errorf("%s: connect: %w", d.Addr, err)
				return
			}
			defer client.Disconnect()
			errs[i] = op(client, d)
		}(i, d)
	}
	wg.Wait()

	var failed []string
	for _, err := range errs {
		if err != nil {
			failed = append(failed, err.Error())
		}
	}
	if len(failed) > 0 {
		return fmt.Errorf("%d/%d devices failed:\n%s",
			len(failed), len(c.devices), strings.Join(failed, "\n"))
	}
	return nil
}

// waitRecordState polls until the record output reaches the wanted state.
func (c *Controller) waitRecordState(client *goobs.Client, d config.Device, active bool) error {
	deadline := time.Now().Add(verifyTimeout)
	for time.Now().Before(deadline) {
		status, err := client.Record.GetRecordStatus()
		if err != nil {
			return fmt.Errorf("%s: verify status: %w", d.Addr, err)
		}
		if status.OutputActive == active {
			return nil
		}
		time.Sleep(verifyInterval)
	}
	return fmt.Errorf("%s: timed out waiting for recording=%v", d.Addr, active)
}

func connect(d config.Device) (*goobs.Client, error) {
	opts := []goobs.Option{}
	if d.Password != "" {
		opts = append(opts, goobs.WithPassword(d.Password))
	}
	return goobs.New(d.Addr, opts...)
}
