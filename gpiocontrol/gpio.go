package gpiocontrol

// GPIOConfig cfg
type GPIOConfig struct {
	FanPin      string `yaml:"fan-pin"`
}

func InitGPIO(gpioConfig *GPIOConfig) error {
	return initGPIONative(gpioConfig)
}

func LEDpwm(dutyPercentage uint32, hertz uint64) {
	pinPWMnative(dutyPercentage, hertz)
}

func FanOn() {
  pinHIGHnative()
}

func FanOff() {
  pinLOWnative()
}
