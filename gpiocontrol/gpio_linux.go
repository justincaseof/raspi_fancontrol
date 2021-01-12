package gpiocontrol

import (
	"fmt"
	"errors"
	gpioperiph "periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/physic"
	iohost "periph.io/x/periph/host"
)

var fanPin gpioperiph.PinOut

func initGPIONative(gpioConfig *GPIOConfig) error {
	fmt.Println("* initializing gpio lib *")
	if _, err := iohost.Init(); err != nil {
		fmt.Println("error initializing gpio lib")
		return err
	}
	fmt.Println("done.")

	fmt.Println("* setting up GPIO pins *")

	// ### FAN PIN ###
	fmt.Println("\t--> initializing pin: ", gpioConfig.FanPin)
	fanPin = gpioreg.ByName(gpioConfig.FanPin)
	if fanPin == nil {
		return errors.New("unable to set up fanPin")
	}
	// we're using 'RisingEdge' to trigger interrupt upon release of pushed button
	//ledPin.PWM(gpioperiph.DutyHalf, 10 * physic.Hertz)
	//ledPin.Out(gpioperiph.High)
	//ledPin.Out(gpioperiph.Low)
	fmt.Println("\t--> done.")

	return nil
}

func pinPWMnative(dutyPercentage uint32, hertz uint64) {
	fanPin.PWM(gpioperiph.Duty(uint32(gpioperiph.DutyMax)*dutyPercentage/100), physic.Frequency(hertz*uint64(physic.Hertz)))
}

func pinHIGHnative() {
	fanPin.Out(gpioperiph.High)
}

func pinLOWnative() {
	fanPin.Out(gpioperiph.Low)
}
