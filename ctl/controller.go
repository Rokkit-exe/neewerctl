package controller

import (
	"fmt"
	"time"

	"github.com/Rokkit-exe/neewerctl/utils"
)

type Ctl struct {
	dev   *Device
	state *State
}

func (c *Ctl) setState(on bool, brightness int, kelvin int) {
	c.state.Power = on
	c.state.Brightness = brightness
	c.state.Temperature = kelvin
}

func (c *Ctl) GetState() *State {
	return c.state
}

func NewCtl(state *State) (*Ctl, error) {
	device, err := Open(state.Port)
	if err != nil {
		return nil, err
	}

	// Set read timeout so Read() doesn't block forever
	if err := device.port.SetReadTimeout(100 * time.Millisecond); err != nil {
		device.Close()
		return nil, fmt.Errorf("failed to set read timeout: %w", err)
	}

	c := &Ctl{
		dev:   device,
		state: state,
	}
	return c, nil
}

func (c *Ctl) Close() {
	if c.dev != nil {
		c.dev.Close()
	}
}

func MakeFrame(on bool, brightness int, kelvin int) []byte {
	b := utils.ClampInt(brightness, 0, 100)

	t := utils.KelvinToTemp(kelvin)
	pwr := byte(0)
	if on {
		pwr = 1
	}

	frame := []byte{0x3A, 0x02, 0x03, pwr, byte(b), t, 0x00}

	sum := 0
	for _, v := range frame[:6] {
		sum += int(v)
	}
	frame = append(frame, byte(sum&0xFF))
	return frame
}

func (c *Ctl) Send(on bool, brightness int, kelvin int) error {
	frame := MakeFrame(on, brightness, kelvin)

	time.Sleep(60 * time.Millisecond)
	err := c.dev.Send(frame)

	if err == nil {
		time.Sleep(100 * time.Millisecond)
		c.setState(on, brightness, kelvin)

		err = c.state.SaveState()
		if err != nil {
			return fmt.Errorf("failed to save state: %w", err)
		}
	}
	return err
}
