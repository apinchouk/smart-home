package main

import (
	"fmt"
	"smart_home/gpio"
	"time"
)

type Application struct {
	gpioKey gpio.IGpioIn
}

func init() {
	fmt.Printf("init\n")
}
func Run() {
	param := gpio.GpioInParam{
		Trigger:  gpio.TrigFalling,
		PullUp:   gpio.PullUpOff,
		Invert:   false,
		JitterMs: 0,
	}

	i := gpio.GetGpioIn(gpio.Pin16, onColdKey, param)

	for {
		select {
		case p := <-i.Channel():
			i.OnEvent(p)
		}
	}
//	time.Sleep(time.Second * 10)
}

func onColdKey(pin gpio.EnPINS2BCM, value int) {
	fmt.Printf("coldKey=%d\n", value)
}
