package main

import (
	"flag"
	"fmt"

	"go.bug.st/serial"
)

func coldProfile(port string) error {
	return setProfile(port, 7000, 100)
}

func sunlightProfile(port string) error {
	return setProfile(port, 5600, 28)
}

func afternoonProfile(port string) error {
	return setProfile(port, 5000, 16)
}

func sunsetProfile(port string) error {
	return setProfile(port, 4500, 16)
}

func candleProfile(port string) error {
	return setProfile(port, 3400, 28)
}

func PowerOff(port string) error {
	// 0x00 = init/power-on frame
	frame := []byte{0x3A, 0x02, 0x03, 0x00, 0x64, 0x00, 0x00}
	// checksum
	sum := 0
	for _, b := range frame {
		sum += int(b)
	}
	frame = append(frame, byte(sum&0xFF))

	mode := &serial.Mode{BaudRate: 115200}
	p, err := serial.Open(port, mode)
	if err != nil {
		return err
	}
	defer p.Close()

	_, err = p.Write(frame)
	return err
}

func PowerOn(port string) error {
	// Brightness = 0
	frame := []byte{0x3A, 0x02, 0x03, 0x01, 0x00, 0x29}
	sum := 0
	for _, b := range frame {
		sum += int(b)
	}
	frame = append(frame, 0x00, byte(sum&0xFF))

	mode := &serial.Mode{BaudRate: 115200}
	p, err := serial.Open(port, mode)
	if err != nil {
		return err
	}
	defer p.Close()

	_, err = p.Write(frame)
	return err
}

func brightnessFrame(val byte) []byte {
	frame := []byte{0x3A, 0x02, 0x03, 0x01, val, 0x29}
	sum := 0
	for _, b := range frame {
		sum += int(b)
	}
	frame = append(frame, 0x00, byte(sum&0xFF))
	return frame
}

func KelvinToTempFrame(k int) []byte {
	// clamp
	val := int(float64(k-2700) * 41.0 / 4300.0)
	frame := []byte{0x3A, 0x02, 0x03, 0x01, 0x64, byte(val), 0x00}
	sum := 0
	for _, b := range frame {
		sum += int(b)
	}
	frame = append(frame, byte(sum&0xFF))
	return frame
}

func setProfile(port string, k, b int) error {
	tempFrame := KelvinToTempFrame(k)
	brightFrame := brightnessFrame(byte(b))
	err := send(port, tempFrame)
	if err != nil {
		return err
	}
	err = send(port, brightFrame)
	return err
}

func send(port string, frame []byte) error {
	mode := &serial.Mode{BaudRate: 115200}
	p, err := serial.Open(port, mode)
	if err != nil {
		fmt.Println("Make sure the device is connected and the port is correct.")
		fmt.Println("Port:", port)
		fmt.Println("Error opening port:", err)
		return err
	}
	defer p.Close()
	_, err = p.Write(frame)
	if err != nil {
		fmt.Println("Port:", port)
		fmt.Println("Error writing to port:", err)
	}
	return err
}

func main() {
	portName := flag.String("port", "/dev/ttyUSB0", "Serial port (default: /dev/ttyUSB0)")
	brightness := flag.Int("brightness", -1, "Brightness 0-100")
	temperature := flag.Int("temperature", -1, "Temperature 0-41 (2700k - 7000k)")
	on := flag.Bool("on", false, "Power on")
	off := flag.Bool("off", false, "Soft power off")
	profile := flag.String("profile", "", "Set profile [cold, sunlight, afternoon, sunset, candle] (default: cold)")

	flag.Parse()

	if *on {
		err := PowerOn(*portName)
		if err != nil {
			panic(err)
		}
		fmt.Println("Powered on")
	}

	if *off {
		err := PowerOff(*portName)
		if err != nil {
			panic(err)
		}
		fmt.Println("Powered off")
	}

	if *temperature >= 0 {
		v := *temperature
		if v < 2700 {
			v = 2700
		}
		if v > 7000 {
			v = 7000
		}
		frame := KelvinToTempFrame(v)
		err := send(*portName, frame)
		if err != nil {
			panic(err)
		}
		fmt.Println("Temperature set to", v)
	}

	if *brightness >= 0 {
		v := *brightness
		if v > 100 {
			v = 100
		}

		frame := brightnessFrame(byte(v))
		err := send(*portName, frame)
		if err != nil {
			panic(err)
		}
		fmt.Println("Brightness set to", v)
	}

	if *profile != "" {
		var err error
		switch *profile {
		case "cold":
			err = coldProfile(*portName)
		case "sunlight":
			err = sunlightProfile(*portName)
		case "afternoon":
			err = afternoonProfile(*portName)
		case "sunset":
			err = sunsetProfile(*portName)
		case "candle":
			err = candleProfile(*portName)
		default:
			fmt.Println("Unknown profile:", *profile)
			return
		}
		if err != nil {
			fmt.Println("Error setting profile:", err)
			panic(err)
		}
		fmt.Println("Profile set to", *profile)
	}
}
