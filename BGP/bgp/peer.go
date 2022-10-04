package bgp

import (
	"log"
	"net"
)

type State int

const (
	Idle State = iota
	Connect
	OpenSent
	OpenConfirm
)

type Event int

const (
	ManualStart Event = iota
	TcpCrAcked
	TcpConnectionConfirmed
)

type Peer struct {
	Config        *Config
	Event_queue   *EventQueue
	Now_state     State
	TcpConnection net.Conn
}

func NewPeer(config *Config) *Peer {
	eventqueue := NewEventQueue()
	return &Peer{
		Config:      config,
		Event_queue: eventqueue,
		Now_state:   Idle,
	}
}

func (peer *Peer) start() {
	peer.Event_queue.enqueue(ManualStart)
}

func (peer *Peer) nextStep() {
	event := peer.Event_queue.dequeue()
	peer.handleEvent(event)
}

func (peer *Peer) createTcpConnectionToRemoteIp() (net.Conn, error) {
	if peer.Config.Mode == Active {
		conn, err := net.Dial("tcp", peer.Config.Remote_ip_address+":179")
		if err != nil {
			log.Println("接続できてない")
			return nil, err
		}
		peer.Event_queue.enqueue(TcpCrAcked)
		return conn, nil
	} else {
		addr := peer.Config.Local_ip_address + ":179"
		listner, err := net.Listen("tcp", addr)
		if err != nil {
			log.Println("179 port can't bind")
			return nil, err
		}

		host, port, err := net.SplitHostPort(listner.Addr().String())
		if err != nil {
			panic(err)
		}
		log.Printf("Listening on host: %s, port: %s\n", host, port)

		conn, err := listner.Accept()
		if err != nil {
			return nil, err
		}
		peer.Event_queue.enqueue(TcpConnectionConfirmed)
		return conn, nil
	}
}

func (peer *Peer) handleEvent(event Event) {
	switch peer.Now_state {
	case Idle:
		switch event {
		case ManualStart:
			conn, err := peer.createTcpConnectionToRemoteIp()
			if err != nil {
				log.Fatalln(err)
			}
			peer.TcpConnection = conn
			peer.Now_state = Connect
		default:
			return
		}
	case Connect:
		switch event {
		case TcpConnectionConfirmed | TcpCrAcked:
			peer.Now_state = OpenSent
		default:
			return
		}
	default:
		return
	}
}
