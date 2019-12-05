package sensorled

import (
	"fmt"
	"time"
	"github.com/apinchouk/smart-home/dispatcher"
	"github.com/apinchouk/smart-home/gpio"
//	"github.com/kelvins/sunrisesunset"
)

type ledState int

const (
	stateLedDisable         ledState = 0
	stateLedTemporaryEnable          = 1
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

type SensorLed struct {
	Dispatcher     dispatcher.IDispatcher
	GpioButtonLed  gpio.IGpioIn
	GpioControlLed gpio.IGpioOut
	state          ledState
	timer          dispatcher.ITimer
}
func isNight() bool {
	t:=time.Now()
	sunset,_:=time.Parse("15:04:05", "16:00:00")
	sunrise,_:=time.Parse("15:04:05", "09:00:00")
	//fmt.Println("Time:", t.Format("15:04:05")) // Sunset: 18:14:27
/*
	p := sunrisesunset.Parameters{
		Latitude:  53.195682,
		Longitude: 45.010420,
		UtcOffset: 3.0,
		Date:      t,
	  }
	  sunrise, sunset, err := p.GetSunriseSunset()

	  if err == nil {
		fmt.Println("Sunrise:", sunrise.Format("15:04:05")) // Sunrise: 06:11:44
		fmt.Println("Sunset:", sunset.Format("15:04:05")) // Sunset: 18:14:27
	} else {
		fmt.Println(err)
	}
	*/
	fmt.Printf("Time:%v %v %v\n", t.Hour(),sunrise.Hour(),sunset.Hour())
	if (t.Hour()>sunset.Hour() || (t.Hour()==sunset.Hour() && t.Minute()>=sunset.Minute())) || 
	(t.Hour()<sunrise.Hour() || (t.Hour()==sunrise.Hour() && t.Minute()<=sunrise.Minute())) {
		fmt.Println("Night\n")
		return true
	} else {
		fmt.Println("Day\n")
	}
	return false
}
func (this *SensorLed) Init() {
	
	this.timer = this.Dispatcher.CreateTimer(this)
	this.onLedDisable(evEntry, nil)
}
func (this *SensorLed) OnGpioEvent(pin gpio.EnPINS2BCM, value int) {

	if value == 0 {
		this.onEvent(evButtonOff, nil)
	} else {
		this.onEvent(evButtonOn, nil)
	}

}

// not implemented
func (this *SensorLed) OnDispatcherEvent(data dispatcher.IEventData) {

}
func (this *SensorLed) OnDispatcherTimer(id int, data dispatcher.IEventData) {
	this.onEvent(ledEvent(data.(int)), nil)
}
func (this *SensorLed) onEvent(event ledEvent, data dispatcher.IEventData) {
	fmt.Printf("SensorLed: %v, State=%d, onEvent=%d\n", time.Now(), this.state, event)

	switch this.state {
	case stateLedDisable:
		this.onLedDisable(event, data)
	case stateLedTemporaryEnable:
		this.onLedTemporatyEnable(event, data)
	
	}

}
func (this *SensorLed) onLedDisable(event ledEvent, data dispatcher.IEventData) {
	switch event {
	case evEntry:
		fmt.Println("Disable")
		this.state = stateLedDisable
		this.GpioControlLed.SetValue(0)
	case evButtonOn:
		if (isNight()) {
			this.onLedTemporatyEnable(evEntry, nil)
		}
	}
}
func (this *SensorLed) onLedTemporatyEnable(event ledEvent, data dispatcher.IEventData) {
	switch event {
	case evEntry:
		fmt.Println("TemporaryEnable")
		this.state = stateLedTemporaryEnable
		this.GpioControlLed.SetValue(1)
		this.timer.Start(time.Minute*10, tmDisable, dispatcher.TimerSingle)
	case evButtonOn:
		this.timer.Start(time.Minute*10, tmDisable, dispatcher.TimerSingle)
	case evButtonOff:
		this.timer.Start(time.Second*30, tmDisable, dispatcher.TimerSingle)
	case tmDisable:
		this.onLedDisable(evEntry, nil)
	}
}
