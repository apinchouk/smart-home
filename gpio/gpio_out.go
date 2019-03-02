package gpio

type gpioOutImplemented struct {
	param  GpioOutParam
	bcmPin EnPINS2BCM
	value  int
}

func (this *gpioOutImplemented) init() bool {

	exportPin(this.bcmPin)
	setPinDirection(this.bcmPin, "out")
	if this.param.PullUp == PullUpOn {
		setPinDirection(this.bcmPin, "high")
	}
	if this.param.ActiveLow == true {
		setPinActiveLow(this.bcmPin, "1")
	} else {
		setPinActiveLow(this.bcmPin, "0")
	}
	return true
}

func (this *gpioOutImplemented) Value() int {
	return this.value
}

func (this *gpioOutImplemented) SetValue(value int) {
	this.value = value
	if value > 0 {
		setPinValue(this.bcmPin, "1")
	} else {
		setPinValue(this.bcmPin, "0")
	}
}
