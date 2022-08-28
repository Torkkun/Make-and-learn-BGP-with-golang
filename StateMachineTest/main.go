package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"
)

type State int

const (
	PowerOn State = iota
	PowerOff
)

type Event int

const (
	PushedPowerButton Event = iota
	PushedVolumeIncreaseButton
	PushedVolumeDecreaseButton
)

type TV struct {
	now_state   State
	event_queue *EventQueue
	volume      uint8
}

type EventQueue []Event

func TvNew() *TV {
	return &TV{
		now_state:   PowerOff,
		event_queue: EventQueueNew(),
		volume:      10,
	}
}

func (tv *TV) be_pushed_power_button() {
	tv.event_queue.enqueue(PushedPowerButton)
}

func (tv *TV) be_pushed_volume_increase_button() {
	tv.event_queue.enqueue(PushedVolumeIncreaseButton)
}

func (tv *TV) be_pushed_volume_decrease_button() {
	tv.event_queue.enqueue(PushedVolumeDecreaseButton)
}

func (tv *TV) handler_event(event Event) {
	switch tv.now_state {
	case PowerOn:
		switch event {
		case PushedPowerButton:
			tv.now_state = PowerOff
		case PushedVolumeIncreaseButton:
			tv.volume += 1
		case PushedVolumeDecreaseButton:
			tv.volume -= 1
		}
	case PowerOff:
		switch event {
		case PushedPowerButton:
			tv.now_state = PowerOn
		default:
			return
		}
	}
}

func EventQueueNew() *EventQueue {
	evq := new(EventQueue)
	return evq
}

func (evq *EventQueue) dequeue() Event {
	result := (*evq)[0]
	*evq = (*evq)[1:]
	return result
}

func (evq *EventQueue) enqueue(event Event) {
	*evq = append(*evq, event)
}

func push_random_button_of_tv(tv *TV) {
	num, err := rand.Int(rand.Reader, big.NewInt(3))
	if err != nil {
		log.Println(err.Error())
		return
	}
	switch num.Int64() {
	case 0:
		tv.be_pushed_power_button()
	case 1:
		tv.be_pushed_volume_increase_button()
	case 2:
		tv.be_pushed_volume_decrease_button()
	default:
		return
	}
}

func main() {
	tv := TvNew()
	tv.be_pushed_power_button()
	for {
		push_random_button_of_tv(tv)
		event := tv.event_queue.dequeue()
		fmt.Printf("tv information: now_state=%v, volume=%v\n input_event: %v\n",
			tv.now_state, tv.volume, event)
		tv.handler_event(event)
		time.Sleep(2 * time.Second)
	}
}
