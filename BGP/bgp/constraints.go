package bgp

// https://www.infraexpert.com/study/bgpz02.html
type FSM_State int

// Activeは今は無し
const (
	Idle FSM_State = iota
	Connect
	OpenSent
	OpenConfirm
	Established
)
