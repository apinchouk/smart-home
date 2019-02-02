package gpio

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

type gpioInImplemented struct {
	param    GpioInParam
	fd       int
	bcmPin   EnPINS2BCM
	callback OnGpioEvent
	value    int
	ch       chan IChannelData
}

func (this *gpioInImplemented) init() bool {
	//	writeInFile("/tmp/test", "%d", 1)

	exportPin(this.bcmPin)
	setPinDirection(this.bcmPin, "in")
	if this.param.PullUp == PullUpOn {
		setPinDirection(this.bcmPin, "high")
	}
	switch this.param.Trigger {
	case TrigNone:
		setPinEdge(this.bcmPin, "none")
	case TrigRising:
		setPinEdge(this.bcmPin, "rising")
	case TrigFalling:
		setPinEdge(this.bcmPin, "falling")
	case TrigBoth:
		setPinEdge(this.bcmPin, "both")
	}

	fileName := fmt.Sprintf(sysFsValue, this.bcmPin)
	fd, err := syscall.Open(fmt.Sprintf(sysFsValue, this.bcmPin), syscall.O_RDONLY, 0)
	if err != nil {
		fmt.Println("syscall.Open: ", err)
		return false
	}
	fmt.Println("Open", fileName)

	if e := syscall.SetNonblock(fd, true); e != nil {
		fmt.Println("SetNonblock:", e)
		os.Exit(1)
	}

	this.fd = int(fd)
	this.ch = make(chan IChannelData, 256)
	go detected(this)
	return true
}

func (this *gpioInImplemented) Channel() chan IChannelData {
	return this.ch
}
func (this *gpioInImplemented) OnEvent(ch IChannelData) {
	this.callback(this.bcmPin, ch.(int))
}
func (this *gpioInImplemented) Value() int {
	return this.value
}

func (this *gpioInImplemented) setValue(value int) {
	this.value = value
	this.ch <- value
	fmt.Printf("value=%d\n", value)
}

func detected(this *gpioInImplemented) {
	defer syscall.Close(this.fd)

	epfd, e := syscall.EpollCreate1(0)
	if e != nil {
		fmt.Println("epoll_create1: ", e)
		os.Exit(1)
	}
	defer syscall.Close(epfd)

	var event syscall.EpollEvent
	var events [1]syscall.EpollEvent

	event.Events = syscall.EPOLLIN | syscall.EPOLLET
	event.Fd = int32(this.fd)
	if e = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, this.fd, &event); e != nil {
		fmt.Println("epoll_ctl: ", e)
		os.Exit(1)
	}

	var buf = make([]byte, 16)
	var value int
	timeout := -1
	for {
		nevents, e := syscall.EpollWait(epfd, events[:], timeout)
		if e != nil {
			fmt.Println("epoll_wait: ", e)
			return
		}
		if nevents == 0 {
			timeout = -1
			this.setValue(value)

			continue
		}
		if nevents == 1 {
			syscall.Seek(this.fd, 0, 0)
			n, err := syscall.Read(this.fd, buf)
			if err == nil && n > 0 {
				value = int(buf[0] - '0')
				if this.param.Invert {
					value ^= 1
				}
				if this.param.JitterMs > 0 {
					timeout = int(this.param.JitterMs)
				} else {
					timeout = -1
					this.setValue(value)
				}
			}
		}
	}
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
