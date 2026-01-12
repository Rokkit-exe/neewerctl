package ctl

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/Rokkit-exe/neewerctl/models"
	"github.com/Rokkit-exe/neewerctl/utils"
	"go.bug.st/serial"
)

var (
	stopMagic     = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	getStateMagic = []byte{0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	socketPath    = "/run/neewer.sock"
	socketPerm    = os.FileMode(0o666)
	socketType    = "unix"
)

func RunDeamon(port string) error {
	os.Remove(socketPath)

	n, err := Open(port)
	if err != nil {
		return fmt.Errorf("failed to open Neewer device: %v", err)
	}
	defer n.Close()

	var (
		stateMu sync.RWMutex
		state   = &models.State{
			Port: port,
		}
		stateReady = make(chan struct{}) // Signal when first state is read
	)

	// Background goroutine to continuously read device state
	go ReadState(n, &stateMu, state, stateReady)

	// Wait for first state reading before accepting connections
	fmt.Println("Waiting for initial state...")
	<-stateReady
	fmt.Println("Initial state received!")

	ln, err := net.Listen(socketType, socketPath)
	if err != nil {
		return fmt.Errorf("failed to start UNIX %s socket listener: %v", socketPath, err)
	}
	defer ln.Close()
	err = os.Chmod(socketPath, socketPerm)
	if err != nil {
		return fmt.Errorf("failed to set socket permissions: %v", err)
	}

	fmt.Println("Neewer daemon running")

	for {
		fmt.Println("Waiting for connection...")
		c, err := ln.Accept()
		if err != nil {
			fmt.Println("Listener error:", err)
			break
		}

		fmt.Println("Connection accepted!")
		go handleClient(c, n, &stateMu, state)
	}

	return nil
}

func ReadState(n *Neewer, stateMu *sync.RWMutex, state *models.State, stateReady chan struct{}) {
	buf := make([]byte, 8)
	fmt.Println("State reader started")
	firstRead := true

	for {
		if _, err := n.port.Read(buf); err != nil {
			fmt.Printf("Serial read error: %v\n", err)
			continue
		}

		// Validate frame
		if buf[0] != 0x3A || buf[1] != 0x02 || buf[2] != 0x03 {
			continue
		}

		// Update state
		stateMu.Lock()
		state.Power = buf[3] == 1
		state.Brightness = int(buf[4])
		state.Temperature = utils.TempByteToKelvin(buf[5])
		fmt.Printf("State updated: Power=%v, Brightness=%d, Temp=%dK\n",
			state.Power, state.Brightness, state.Temperature)
		stateMu.Unlock()

		if firstRead {
			close(stateReady)
			firstRead = false
		}
	}
}

func handleClient(c net.Conn, n *Neewer, stateMu *sync.RWMutex, state *models.State) {
	defer c.Close()

	buf := make([]byte, 8)

	// Stop command
	if bytes.Equal(buf, stopMagic) {
		fmt.Println("Stop command received")
		c.Write([]byte{0x01})
		os.Remove("/run/neewer.sock")
		return
	}

	// Get state command
	if buf[0] == 0xFF && buf[1] == 0x00 {
		fmt.Println("Get state command received")
		stateMu.RLock()
		response := []byte{
			byte(utils.BoolToInt(state.Power)),
			byte(state.Brightness),
			utils.KelvinToTemp(state.Temperature),
		}
		stateMu.RUnlock()

		fmt.Printf("Sending state: device=%s, power=%v, brightness=%d, temp=%d\n",
			state.Port, state.Power, state.Brightness, state.Temperature)

		written, err := c.Write(response)
		if err != nil {
			fmt.Printf("Error writing response: %v\n", err)
		} else {
			fmt.Printf("Wrote %d bytes\n", written)
		}

		// Give client time to read before closing
		time.Sleep(50 * time.Millisecond)
		return
	}

	// Set command - send to device and acknowledge
	fmt.Printf("Set command received: %x\n", buf)
	n.Send(buf)
	c.Write([]byte{0x01})
}

func Connect() (net.Conn, error) {
	c, err := net.Dial("unix", "/run/neewer.sock")
	if err != nil {
		return nil, fmt.Errorf("neewerd is not running\nTry starting it with:\n	sudo systemctl start neewerd.service\n	or\n sudo neewerctl deamon start")
	}
	return c, nil
}

func GetState(port string) (*models.State, error) {
	c, err := Connect()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	if _, err := c.Write(getStateMagic); err != nil {
		return nil, fmt.Errorf("write error: %w", err)
	}

	buf := make([]byte, 3)
	n, err := c.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("read error (got %d bytes): %w", n, err)
	}

	return &models.State{
		Port:        port,
		Power:       buf[0] == 1,
		Brightness:  int(buf[1]),
		Temperature: utils.TempByteToKelvin(buf[2]),
	}, nil
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

func PowerOff(port string) error {
	return Send(port, MakeFrame(false, 0, 2900))
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

func Send(_ string, frame []byte) error {
	c, err := Connect()
	if err != nil {
		return err
	}
	defer c.Close()

	if _, err = c.Write(frame); err != nil {
		return err
	}

	ack := make([]byte, 1)
	_, err = c.Read(ack)
	return err
}
