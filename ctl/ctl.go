package ctl

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"go.bug.st/serial"
)

var stopMagic = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

func copyFile(src, dst string) error {
	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destFile.Close()

	// Copy contents
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy contents: %v", err)
	}

	// Optionally copy file permissions
	info, err := sourceFile.Stat()
	if err == nil {
		os.Chmod(dst, info.Mode())
	}

	return nil
}

func InstallService() error {
	projectRoot, _ := os.Getwd() // or use BinaryRoot() if installed
	src := filepath.Join(projectRoot, "service", "neewerd.service")
	dst := "/etc/systemd/system/neewerd.service"

	// Open source
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Create destination
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy contents
	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	// Optional: set permissions
	if err := os.Chmod(dst, 0o644); err != nil {
		return err
	}

	// Reload systemd to pick up the new service
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		return err
	}

	fmt.Println("Service installed successfully")
	return nil
}

func RunDeamon(port string) error {
	os.Remove("/run/neewer.sock")

	n, err := Open(port)
	if err != nil {
		return fmt.Errorf("failed to open Neewer device: %v", err)
	}
	defer n.Close()

	ln, err := net.Listen("unix", "/run/neewer.sock")
	if err != nil {
		return fmt.Errorf("failed to start UNIX socket listener: %v", err)
	}
	defer ln.Close()
	os.Chmod("/run/neewer.sock", 0o666)

	fmt.Println("Neewer daemon running...") // optional logging

	// Blocking accept loop
	for {
		c, err := ln.Accept()
		if err != nil {
			// Only break on fatal errors
			fmt.Println("Listener error:", err)
			break
		}

		go func(c net.Conn) {
			defer c.Close()
			buf := make([]byte, 8)
			if _, err := c.Read(buf); err != nil {
				return
			}

			if bytes.Equal(buf, stopMagic) {
				os.Remove("/run/neewer.sock")
				ln.Close()
				return
			}

			n.Send(buf)
		}(c)
	}

	return nil
}

func StartDeamon() error {
	out, err := exec.Command("systemctl", "start", "neewerd.service").CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start neewerd service: %v, output: %s", err, string(out))
	}
	return nil
}

func StopDeamon() error {
	out, err := exec.Command("systemctl", "stop", "neewerd.service").CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop neewerd service: %v, output: %s", err, string(out))
	}
	return nil
}

// func StopDeamon() error {
// 	c, err := net.Dial("unix", "/run/neewer.sock")
// 	if err != nil {
// 		return fmt.Errorf("neewerd is not running")
// 	}
// 	defer c.Close()
// 	_, err = c.Write([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
// 	return err
// }

type Neewer struct {
	port serial.Port
}

func Open(dev string) (*Neewer, error) {
	p, err := serial.Open(dev, &serial.Mode{BaudRate: 115200})
	if err != nil {
		return nil, err
	}

	p.ResetInputBuffer()
	p.ResetOutputBuffer()
	p.Write([]byte{0, 0, 0, 0})
	time.Sleep(80 * time.Millisecond)

	return &Neewer{port: p}, nil
}

func (n *Neewer) Send(frame []byte) error {
	time.Sleep(60 * time.Millisecond)
	_, err := n.port.Write(frame)
	return err
}

func (n *Neewer) Close() {
	n.port.Close()
}

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

func Send(_ string, frame []byte) error {
	c, err := net.Dial("unix", "/run/neewer.sock")
	if err != nil {
		return fmt.Errorf("neewerd is not running")
	}
	defer c.Close()
	_, err = c.Write(frame)
	return err
}

// func Send(port string, frame []byte) error {
// 	mode := &serial.Mode{BaudRate: 115200}
// 	p, err := serial.Open(port, mode)
// 	if err != nil {
// 		return fmt.Errorf("Make sure the device is connected and the port is correct.\nError opening port %s: %v", port, err)
// 	}
// 	defer p.Close()
// 	_, err = p.Write(frame)
// 	if err != nil {
// 		return fmt.Errorf("Error writing to port %s: %v", port, err)
// 	}
// 	return nil
// }

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
