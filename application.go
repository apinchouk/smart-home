package main

import (
	"fmt"

	"github.com/apinchouk/smart-home/dispatcher"
	"github.com/apinchouk/smart-home/gpio"
)

type IApplication interface {
	init() bool
	run()
}
type Application struct {
	dispatcher dispatcher.IDispatcher
	stairWay   StairWayLed
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
	return true
}

func (this *Application) run() {

	this.dispatcher.Run()
}
