package ctl

import (
	"time"

	"go.bug.st/serial"
)

type Device struct {
	port serial.Port
}

func Open(dev string) (*Device, error) {
	p, err := serial.Open(dev, &serial.Mode{
		BaudRate: 115200,
	})
	if err != nil {
		return nil, err
	}

	p.SetReadTimeout(500 * time.Millisecond)
	p.ResetInputBuffer()
	p.ResetOutputBuffer()
	p.Write([]byte{0, 0, 0, 0})
	time.Sleep(80 * time.Millisecond)

	return &Device{port: p}, nil
}

func (n *Device) Send(frame []byte) error {
	time.Sleep(60 * time.Millisecond)
	_, err := n.port.Write(frame)
	return err
}

func (n *Device) Close() {
	n.port.Close()
}
