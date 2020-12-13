/*
A blinker example using pi-go-gpio library.
Toggles two LEDs on physical pin 17/27.
Connect two LEDs with resistors from pin 17 and pin 27 to ground.
*/
package main

import (
	"github.com/wolfgang-werner/pi-go-gpio"
	"time"
)

func main() {
	gpio17 := gpio.Open(17, gpio.Out, false)
	defer gpio17.Close()

	gpio27 := gpio.Open(27, gpio.Out, false)
	defer gpio27.Close()

	gpio17.SetValue(gpio.High)
	for i := 0; i < 50; i++ {
		time.Sleep(200 * time.Millisecond)
		gpio17.ToggleValue()
		gpio27.ToggleValue()
	}
}
