package main

import (
	"fmt"
	"time"

	"github.com/apinchouk/smart-home/dispatcher"
	"github.com/apinchouk/smart-home/gpio"
)

type ledState int

const (
	stateLedDisable         ledState = 0
	stateLedTemporaryEnable          = 1
	stateLedEnable                   = 2
)

type ledEvent int

const (
	evEntry     ledEvent = 0
	evButtonOn           = 1
	evButtonOff          = 2
	tmProtect            = 3
	tmEnable             = 4
	tmDisable            = 5
)

type StairWayLed struct {
	dispatcher     dispatcher.IDispatcher
	gpioButtonLed  gpio.IGpioIn
	gpioControlLed gpio.IGpioOut
	state          ledState
	timer1         dispatcher.ITimer
	timer2         dispatcher.ITimer
	protect        bool
}

func (this *StairWayLed) Init() {
	this.timer1 = this.dispatcher.CreateTimer(this)
	this.timer2 = this.dispatcher.CreateTimer(this)
	this.onLedDisable(evEntry, nil)
}
func (this *StairWayLed) OnGpioEvent(pin gpio.EnPINS2BCM, value int) {

	if value == 0 {
		this.onEvent(evButtonOff, nil)
	} else {
		this.onEvent(evButtonOn, nil)
	}

}

// not implemented
func (this *StairWayLed) OnDispatcherEvent(data dispatcher.IEventData) {

}
func (this *StairWayLed) OnDispatcherTimer(id int, data dispatcher.IEventData) {
	this.onEvent(ledEvent(data.(int)), nil)
}
func (this *StairWayLed) onEvent(event ledEvent, data dispatcher.IEventData) {
	fmt.Printf("%v, State=%d, onEvent=%d\n", time.Now(), this.state, event)

	switch this.state {
	case stateLedDisable:
		this.onLedDisable(event, data)
	case stateLedTemporaryEnable:
		this.onLedTemporatyEnable(event, data)
	case stateLedEnable:
		this.onLedEnable(event, data)
	}

}
func (this *StairWayLed) onLedDisable(event ledEvent, data dispatcher.IEventData) {
	switch event {
	case evEntry:
		fmt.Println("Disable")
		this.state = stateLedDisable
		this.gpioControlLed.SetValue(0)
		this.timer1.Start(time.Second*2, tmProtect, dispatcher.TimerSingle)
		this.protect = true
	case evButtonOn:
		if this.protect == false {
			this.onLedTemporatyEnable(evEntry, nil)
		}
	case tmProtect:
		this.protect = false
	}
}
func (this *StairWayLed) onLedTemporatyEnable(event ledEvent, data dispatcher.IEventData) {
	switch event {
	case evEntry:
		fmt.Println("TemporaryEnable")
		this.state = stateLedTemporaryEnable
		this.gpioControlLed.SetValue(1)
		this.timer1.Start(time.Second*2, tmProtect, dispatcher.TimerSingle)
		this.timer2.Start(time.Second*60, tmDisable, dispatcher.TimerSingle)
		this.protect = true
	case evButtonOn:
		if this.protect == false {
			this.onLedDisable(evEntry, nil)
		}
	case tmProtect:
		this.protect = false
		if this.gpioButtonLed.Value() == 1 {
			this.onLedEnable(evEntry, nil)
		}
	case evButtonOff:
		if this.protect == false {
			this.timer1.Stop()
		}

	case tmDisable:
		this.onLedDisable(evEntry, nil)
	}

}
func (this *StairWayLed) onLedEnable(event ledEvent, data dispatcher.IEventData) {
	switch event {
	case evEntry:
		fmt.Println("Enable")
		this.state = stateLedEnable
		this.gpioControlLed.SetValue(1)
		this.timer1.Start(time.Second*2, tmProtect, dispatcher.TimerSingle)
		this.protect = true
	case evButtonOn:
		if this.protect == false {
			this.onLedDisable(evEntry, nil)
		}
	case tmProtect:
		this.protect = false
	}
}
