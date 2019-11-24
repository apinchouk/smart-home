package main

import (
	"fmt"

	"github.com/apinchouk/smart-home/dispatcher"
	"github.com/apinchouk/smart-home/gpio"
	"github.com/apinchouk/smart-home/sensorled"
)

type IApplication interface {
	init() bool
	run()
}
type Application struct {
	dispatcher dispatcher.IDispatcher
	stairWay   StairWayLed
	sensorLed  sensorled.SensorLed
}

func (this *Application) init() bool {

	fmt.Printf("init\n")
	this.dispatcher = dispatcher.GetDispatcher()

	relay := gpio.GpioOutParam{
		PullUp:    gpio.PullUpOff,
		ActiveLow: true,
	}

	button := gpio.GpioInParam{
		Trigger:  gpio.TrigBoth,
		PullUp:   gpio.PullUpOn,
		Invert:   true,
		JitterMs: 100,
	}

	this.stairWay = StairWayLed{dispatcher: this.dispatcher}
	this.stairWay.gpioButtonLed = gpio.GetGpioIn(app.dispatcher, gpio.Pin40, &this.stairWay, button)
	this.stairWay.gpioControlLed = gpio.GetGpioOut(gpio.Pin3, relay)
	this.stairWay.Init()

	hallkey := gpio.GpioInParam{
		Trigger:  gpio.TrigBoth,
		PullUp:   gpio.PullUpOff,
		Invert:   false,
		JitterMs: 100,
	}

	this.sensorLed = sensorled.SensorLed{Dispatcher: this.dispatcher}
	this.sensorLed.GpioButtonLed = gpio.GetGpioIn(app.dispatcher, gpio.Pin13, &this.sensorLed, hallkey)
	this.sensorLed.GpioControlLed = gpio.GetGpioOut(gpio.Pin15, relay)
	this.sensorLed.Init()

	return true
}

func (this *Application) run() {

	this.dispatcher.Run()
}
