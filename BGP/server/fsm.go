package server

import (
	"mybgp/bgp"
	"mybgp/config"
	"net"
)

// fsm gobgp実装
// https://github.com/osrg/gobgp/blob/da00912b2fbc96ed272947a4005afff15826e526/pkg/server/fsm.go#L173

type fsm struct {
	config *config.Neighbor
	state  bgp.FSM_State
}

type fsmHandler struct {
	fsm  *fsm
	conn net.Conn
}

func newFSM(pConf *config.Neighbor) *fsm {
	fsm := &fsm{
		config: pConf,
		state:  bgp.Idle,
	}
	return fsm
}

// msg fsm
func (fsm *fsm) bgpMessageStateUpdate(msgType uint8) {

	switch msgType {
	case bgp.MSG_Open:

	}
}
