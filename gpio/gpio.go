package gpio

import (
	"fmt"
	"os"
	"time"

	"github.com/apinchouk/smart-home/dispatcher"
)

type EnPINS2BCM int
type EnPullState int
type EnTrigger int

const (
	Pin11 EnPINS2BCM = 17 // холодный контур
	Pin13            = 27
	Pin8             = 14 // зеленый
	Pin10            = 15 // красный
	Pin12            = 18 // синий
	Pin16            = 23 // датчик потока
	Pin40            = 21 // кнопки включения света
	Pin3             = 2  // реле управления светом
)

//
const (
	PullUpOff EnPullState = 0
	PullUpOn              = 1
)

// условие срабатывания триггера для входного GPIO
const (
	TrigNone    EnTrigger = 0
	TrigRising            = 1
	TrigFalling           = 2
	TrigBoth              = 3
)

type GpioInParam struct {
	Invert   bool
	Trigger  EnTrigger
	PullUp   EnPullState
	JitterMs uint32
}
type GpioOutParam struct {
	ActiveLow bool
	PullUp    EnPullState
}

type IChannelData interface {
}

//
type IGpioIn interface {
	init() bool
	Value() int
}
type IGpioOut interface {
	init() bool
	Value() int
	SetValue(value int)
}

func writeInFile(file string, format string, v ...interface{}) bool {

	fd, err := os.Create(file)
	if err != nil {
		return false
	}
	defer fd.Close()

	_, err = fmt.Fprintf(fd, fmt.Sprintf(format, v...))

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return false
	} else {
		//	fmt.Printf("write in %s=%v\n", file, fmt.Sprintf(format, v...))
	}

	return true
}

func exportPin(bcmPin EnPINS2BCM) bool {
	if writeInFile(sysFsExport, "%d", bcmPin) {
		time.Sleep(100 * time.Millisecond)
		return true
	}

	return false
}
func setPinDirection(bcmPin EnPINS2BCM, direction string) bool {
	return writeInFile(fmt.Sprintf(sysFsDirection, bcmPin), direction)
}
func setPinEdge(bcmPin EnPINS2BCM, edge string) bool {
	return writeInFile(fmt.Sprintf(sysFsEdge, bcmPin), edge)
}
func setPinValue(bcmPin EnPINS2BCM, value string) bool {
	return writeInFile(fmt.Sprintf(sysFsValue, bcmPin), value)
}
func setPinActiveLow(bcmPin EnPINS2BCM, value string) bool {
	return writeInFile(fmt.Sprintf(sysFsActiveLow, bcmPin), value)
}

const (
	sysFsExport    string = "/sys/class/gpio/export"
	sysFsUnexport  string = "/sys/class/gpio/unexport"
	sysFsDirection string = "/sys/class/gpio/gpio%d/direction"
	sysFsValue     string = "/sys/class/gpio/gpio%d/value"
	sysFsActiveLow string = "/sys/class/gpio/gpio%d/active_low"
	sysFsEdge      string = "/sys/class/gpio/gpio%d/edge"
)

type IGpioEvent interface {
	OnGpioEvent(pin EnPINS2BCM, value int)
}

//getGpioIn
func GetGpioIn(dispatcher dispatcher.IDispatcher, pin EnPINS2BCM, callback IGpioEvent, param GpioInParam) IGpioIn {
	gpio := gpioInImplemented{dispatcher: dispatcher, bcmPin: pin, param: param, callback: callback}

	gpio.init()
	return &gpio
}

func GetGpioOut(pin EnPINS2BCM, param GpioOutParam) IGpioOut {
	gpio := gpioOutImplemented{bcmPin: pin, param: param}

	gpio.init()
	return &gpio
}
