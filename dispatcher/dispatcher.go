package dispatcher

import (
	"fmt"
	"time"
)

type IEventData interface {
}

type ICallbackEvent interface {
	OnDispatcherEvent(data IEventData)
	OnDispatcherTimer(id int, data IEventData)
}

type EnTimerType int

const (
	TimerSingle     EnTimerType = 0
	TimerPeriodical             = 1
)

type ITimer interface {
	Start(duration time.Duration, data IEventData, param EnTimerType)
	Stop()
	Kill()
}
type IDispatcher interface {
	SendEvent(callback ICallbackEvent, data IEventData)
	CreateTimer(callback ICallbackEvent) ITimer
	Run()

	dispatchTimers()
	removeTimer(id int)
}
type TimerImp struct {
	dispatcher *DispatherImplemented
	param      EnTimerType
	callback   ICallbackEvent
	data       IEventData
	duration   time.Duration
	left       time.Duration
	id         int
}
type EventItem struct {
	callback ICallbackEvent
	data     IEventData
}
type DispatherImplemented struct {
	eventQueue   chan EventItem
	timerQueue   []*TimerImp
	timerQuality time.Duration
}

func GetDispatcher() IDispatcher {
	dispatcher := DispatherImplemented{}
	dispatcher.Init()
	return &dispatcher
}
func (this *DispatherImplemented) Init() {
	this.eventQueue = make(chan EventItem, 256)
	this.timerQueue = make([]*TimerImp, 256)
	this.timerQuality = time.Millisecond * 100
}
func (this *DispatherImplemented) SendEvent(callback ICallbackEvent, data IEventData) {

	this.eventQueue <- EventItem{callback, data}
}
func (this *DispatherImplemented) CreateTimer(callback ICallbackEvent) ITimer {
	timer := TimerImp{dispatcher: this, callback: callback}
	for index := range this.timerQueue {
		if this.timerQueue[index] == nil {
			this.timerQueue[index] = &timer
			timer.id = index
			fmt.Println("Assign timer in ", index)
			return &timer
		}
	}
	return nil
}
func (this *DispatherImplemented) removeTimer(id int) {
	this.timerQueue[id] = nil
}
func (this *TimerImp) Start(duration time.Duration, data IEventData, param EnTimerType) {
	this.duration = duration / this.dispatcher.timerQuality
	this.data = data
	this.param = param
	this.left = this.duration
}
func (this *TimerImp) Stop() {
	this.left = 0
}
func (this *TimerImp) Kill() {
	if this.id != -1 {
		this.dispatcher.removeTimer(this.id)
		this.id = -1
	}
}

func (this *DispatherImplemented) Run() {

	tm := make(chan bool, 16)

	go func() {
		for {
			time.Sleep(this.timerQuality)
			tm <- true
		}
	}()

	for {
		select {
		case p := <-this.eventQueue:
			p.callback.OnDispatcherEvent(p.data)
		case <-tm:
			this.dispatchTimers()
		}
	}
}
func (this *DispatherImplemented) dispatchTimers() {
	for _, tm := range this.timerQueue {
		if tm != nil {
			if tm.left > 0 {
				tm.left--
				if tm.left == 0 {
					tm.callback.OnDispatcherTimer(tm.id, tm.data)
					if tm.param == TimerPeriodical {
						tm.left = tm.duration
					}
				}
			}
		}
	}
}
