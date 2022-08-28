package bgp

type EventQueue []Event

func NewEventQueue() *EventQueue {
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
