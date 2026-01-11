package ctl

import (
	"fmt"

	"go.bug.st/serial"
)

func TempByteToKelvin(b byte) int {
	return 2900 + int(b-1)*4100/40
}

func PowerOff(port string) error {
	return Send(port, MakeFrame(false, 0, 2900))
}

func kelvinToTemp(k int) byte {
	if k < 2900 {
		k = 2900
	}
	if k > 7000 {
		k = 7000
	}
	return byte(((k - 2900) * 40 / 4100) + 1)
}

func MakeFrame(on bool, brightness int, kelvin int) []byte {
	if brightness < 0 {
		brightness = 0
	}
	if brightness > 100 {
		brightness = 100
	}

	temp := kelvinToTemp(kelvin)
	pwr := byte(0)
	if on {
		pwr = 1
	}

	frame := []byte{0x3A, 0x02, 0x03, pwr, byte(brightness), temp, 0x00}

	sum := 0
	for i := 0; i < 6; i++ { // 3A 02 03 PWR BRIGHT TEMP ONLY
		sum += int(frame[i])
	}
	frame = append(frame, byte(sum&0xFF))
	return frame
}

func Send(port string, frame []byte) error {
	mode := &serial.Mode{BaudRate: 115200}
	p, err := serial.Open(port, mode)
	if err != nil {
		return fmt.Errorf("Make sure the device is connected and the port is correct.\nError opening port %s: %v", port, err)
	}
	defer p.Close()
	_, err = p.Write(frame)
	if err != nil {
		return fmt.Errorf("Error writing to port %s: %v", port, err)
	}
	return nil
}

func GetProfileValues(profile string) (int, int, error) {
	switch profile {
	case "cold":
		return 7000, 100, nil
	case "sunlight":
		return 5600, 28, nil
	case "afternoon":
		return 5000, 16, nil
	case "sunset":
		return 4500, 16, nil
	case "candle":
		return 3400, 28, nil
	default:
		return 0, 0, fmt.Errorf("invalid profile: %s", profile)
	}
}
