package gpio

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

type IChannelData interface {
}

//
type IGpioIn interface {
	init() bool
	Value() int
	Channel() chan IChannelData
	OnEvent(ch IChannelData)
}

const (
	sysFsExport    string = "/sys/class/gpio/export"
	sysFsUnexport  string = "/sys/class/gpio/unexport"
	sysFsDirection string = "/sys/class/gpio/gpio%d/direction"
	sysFsValue     string = "/sys/class/gpio/gpio%d/value"
	sysFsEdge      string = "/sys/class/gpio/gpio%d/edge"
)

type OnGpioEvent func(EnPINS2BCM, int)

//getGpioIn
func GetGpioIn(pin EnPINS2BCM, callback OnGpioEvent, param GpioInParam) IGpioIn {
	gpio := gpioInImplemented{bcmPin: pin, param: param, callback: callback}

	gpio.init()
	return &gpio
}
