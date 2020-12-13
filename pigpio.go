package gpio

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	gpioSystemBasePath = "/sys/class/gpio"
	gpioIn             = "in"
	gpioOut            = "out"
	gpioLow            = "0"
	gpioHigh           = "1"
)

type DirectionType int

const (
	In DirectionType = iota
	Out
)

func (direction DirectionType) String() string {
	return [...]string{gpioIn, gpioOut}[direction]
}

type ValueType int

const (
	Low ValueType = iota
	High
)

func (value ValueType) String() string {
	return [...]string{gpioLow, gpioHigh}[value]
}

type Pin struct {
	Port        int
	Direction   DirectionType
	Value       ValueType
	Debug       bool
	Initialized bool
	Error       error
}

func Open(port int, direction DirectionType, debug bool) (gpio *Pin) {

	// limit port to [1..27]
	if port < 1 || 27 < port {
		fmt.Printf("GPIO: invalid port number %d\n", port)
		return
	}

	gpio = &Pin{
		Port:        port,
		Direction:   In,
		Value:       High,
		Debug:       debug,
		Initialized: false,
		Error:       nil,
	}

	if gpio.export(); gpio.Error != nil {
		fmt.Printf("GPIO %2d: error creating device: %e\n", gpio.Port, gpio.Error)
		return
	}

	if gpio.SetDirection(direction); gpio.Error != nil {
		fmt.Printf("GPIO %2d: error setting direction: %e\n", gpio.Port, gpio.Error)
		return
	}

	if gpio.SetValue(Low); gpio.Error != nil {
		fmt.Printf("GPIO %2d: error setting value: %e\n", gpio.Port, gpio.Error)
		return
	}

	gpio.Initialized = true
	gpio.Error = nil
	return
}

func (g *Pin) port() string {
	return strconv.Itoa(g.Port)
}

func (g *Pin) exportFilename() string {
	return fmt.Sprintf("%s/export", gpioSystemBasePath)
}

func (g *Pin) unexportFilename() string {
	return fmt.Sprintf("%s/unexport", gpioSystemBasePath)
}

func (g *Pin) directionFilename() string {
	return fmt.Sprintf("%s/gpio%s/direction", gpioSystemBasePath, g.port())
}

func (g *Pin) valueFilename() string {
	return fmt.Sprintf("%s/gpio%s/value", gpioSystemBasePath, g.port())
}

func closeHelper(f *os.File) {
	if err := f.Close(); err != nil {
		fmt.Printf("GPIO: error closing file '%s': %e\n", f.Name(), err)
	}
}

func (g *Pin) waitForExport() {
	// wait for creation of control files
	if g.Debug {
		fmt.Printf("GPIO %2d: waiting for initialization ", g.Port)
	}
	for {
		fd, err := os.OpenFile(g.directionFilename(), os.O_RDWR, 0)
		if err == nil {
			closeHelper(fd)
			if g.Debug {
				fmt.Println(". done.")
			}
			break
		} else {
			if g.Debug {
				fmt.Print(".")
			}
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func (g *Pin) export() {
	if g.Debug {
		fmt.Printf("GPIO %2d: exporting...\n", g.Port)
	}

	f, err := os.OpenFile(g.exportFilename(), os.O_WRONLY, 0)
	if err != nil {
		fmt.Printf("GPIO %2d: can't open export file, error: %e\n", g.Port, err)
		g.Error = err
		return
	}
	defer closeHelper(f)

	if _, err = f.WriteString(g.port()); err != nil {
		fmt.Printf("GPIO %2d: can't write to export file, error: %e\n", g.Port, err)
		g.Error = err
		return
	}

	g.waitForExport()
	g.Error = nil
}

func (g *Pin) SetDirection(direction DirectionType) {
	if g.Error != nil {
		fmt.Printf("GPIO %2d: can't perform Output: Not initialized, error was: %e\n", g.Port, g.Error)
		return
	}

	if g.Direction == direction && g.Initialized {
		if g.Debug {
			fmt.Printf("GPIO %2d: direction: %s already set\n", g.Port, direction.String())
		}
	} else {
		if g.Debug {
			fmt.Printf("GPIO %2d: direction: %s\n", g.Port, direction.String())
		}

		f, err := os.OpenFile(g.directionFilename(), os.O_RDWR, 0)
		if err != nil {
			fmt.Printf("GPIO %2d: can't open direction file, error: %e\n", g.Port, err)
			g.Error = err
			return
		}
		defer closeHelper(f)

		if _, err = f.WriteString(direction.String()); err != nil {
			fmt.Printf("GPIO %2d: can't write to direction file, error: %e\n", g.Port, err)
			g.Error = err
			return
		}

		g.Direction = direction
	}

	g.Error = nil
}

func (g *Pin) SetValue(value ValueType) {
	if g.Error != nil {
		fmt.Printf("GPIO %2d: can't perform Value: Not initialized, error was: %e\n", g.Port, g.Error)
		return
	}

	if g.Value == value && g.Initialized {
		if g.Debug {
			fmt.Printf("GPIO %2d: value: %s already set\n", g.Port, value.String())
		}
	} else {
		if g.Debug {
			fmt.Printf("GPIO %2d: value: %s\n", g.Port, value.String())
		}

		f, err := os.OpenFile(g.valueFilename(), os.O_RDWR, 0)
		if err != nil {
			fmt.Printf("GPIO %2d: can't open value file, error: %e\n", g.Port, err)
			g.Error = err
			return
		}
		defer closeHelper(f)

		if _, err = f.WriteString(value.String()); err != nil {
			fmt.Printf("GPIO %2d: can't write to value file, error: %e\n", g.Port, err)
			g.Error = err
			return
		}

		g.Value = value
	}

	g.Error = nil
}

func (g *Pin) ToggleValue() {
	if g.Error != nil {
		fmt.Printf("GPIO %2d: can't perform Toggle: Not initialized, error was: %e\n", g.Port, g.Error)
		return
	}

	g.SetValue(map[ValueType]ValueType{High: Low, Low: High}[g.Value])
}

func (g *Pin) Close() {
	if g.Debug {
		fmt.Printf("GPIO %2d: shutdown...\n", g.Port)
	}

	g.SetValue(Low)

	f, err := os.OpenFile(g.unexportFilename(), os.O_WRONLY, 0)
	if err != nil {
		fmt.Printf("GPIO %2d: can't open unexport file, error: %e\n", g.Port, err)
		g.Error = err
		return
	}
	defer closeHelper(f)

	if _, err = f.WriteString(g.port()); err != nil {
		fmt.Printf("GPIO %2d: can't write to unexport file, error: %e\n", g.Port, err)
		g.Error = err
		return
	}

	g.Error = nil
}
