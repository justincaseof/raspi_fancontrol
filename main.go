package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	gpiocontrol "raspi_fancontrol/gpiocontrol"
	readtemp "raspi_fancontrol/readtemp"
	"syscall"
	"time"
	"gopkg.in/yaml.v2"
)

const CONFIG_FILENAME = "config.yml"

var processing bool

type AppConfig struct {
	IsSimulation   bool                   `yaml:"is-simulation"`
	GPIOconfig     gpiocontrol.GPIOConfig `yaml:"gpio-config"`
	FanOnTemp      float32		      `yaml:"fan-on-temp"`
        FanOffTemp     float32                `yaml:"fan-off-temp"`
	CheckInterval  uint32                 `yaml:"check-interval"`
}

var appConfig = AppConfig {CheckInterval: 5}	// defaults go here

func main() {
	fmt.Println("### STARTUP")

	// READ CONFIG
	readConfig(&appConfig)

	// INIT
	err := gpiocontrol.InitGPIO(&appConfig.GPIOconfig)
	if err != nil {
		fmt.Println("Cannot set up GPIO ", err)
		panic("Cannot set up GPIO")
	}

	// GO
	go mainLoop()

	// wait indefinitely until external abortion
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) // Ctrl + c
	<-sigs
	fmt.Println("### EXIT")
}

// ==== I/O and properties ====

func readConfig(appConfig *AppConfig) {
	var err error
	var bytes []byte
	bytes, err = ioutil.ReadFile(CONFIG_FILENAME)
	if err != nil {
		fmt.Println("Cannot open config file: ", CONFIG_FILENAME)
		panic(err)
	}
	err = yaml.Unmarshal(bytes, appConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println("Config parsed.")
}

func mainLoop() {
	processing = false
	handleState()

	for {
		select {
		case <-time.After(time.Duration(appConfig.CheckInterval) * time.Second):
			{
				handleState()
			}
		}
	}
}

func handleState() {
	temperatureInfo, err := readtemp.GetTemp()
	if err!=nil {
		fmt.Println("Cannot read temperature: ", err)
	} else {
		fmt.Println("Temp: ", temperatureInfo.Value)
		if temperatureInfo.Value >= appConfig.FanOnTemp {
                  fmt.Println("Fan ON")
		  gpiocontrol.FanOn()
		}
                if temperatureInfo.Value <= appConfig.FanOffTemp {
		  fmt.Println("Fan OFF")
                  gpiocontrol.FanOff()
                }

	}
}
