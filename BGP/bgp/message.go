package bgp

import (
	"time"
)

type BgpMessageInterface interface {
	Serialize()
	Deserialize() []uint8
}

// Message type
const (
	_ = iota
	Open
)

// https://www.rfc-editor.org/rfc/rfc4271#section-4.1
type BgpMessageHeader struct {
	Length uint16
	Type   int
}

func (bgpMsgH *BgpMessageHeader) Serialize() []uint8 {
	var bytes []uint8
	// BGPコネクションに認証を使用する場合に使う。もし、認証を使用しない場合はすべて"1"が入る。
	var marker []uint8 = []uint8{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	// ヘッダを含むメッセージ全体の長さをオクテット単位で定義されます。フィールドの値は最小で19byte、最大で4096byte
	var length = bgpMsgH.Length
	// メッセージのタイプコードが定義されます。
	//1：OPEN
	//2：UPDATE
	//3：NOTIFICATION
	//4：KEEPALIVE
	var msgtype = bgpMsgH.Type
	bytes = append(bytes, marker...)
	bytes = append(bytes, uint8(length))
	bytes = append(bytes, uint8(msgtype))
	return bytes
}

func (bgpMsgH *BgpMessageHeader) Deserialize(bytes []uint8) *BgpMessageHeader {
	var length = bytes[16:18]
	return &BgpMessageHeader{
		Length: uint16(length),
		Type:   int(bytes[18]),
	}
}

type BgpVersion uint8

type BgpOpenMessage struct {
	Header                      BgpMessageHeader
	Version                     BgpVersion
	My_autonomous_system_number AutonomousSystemNumber
	// HoldTime
	// https://www.infraexpert.com/study/bgpz07.html
	HoldTime                time.Time
	BgpIdentifier           string
	OptionalParameterLength uint8
	OptionalParameters      []uint8
}

func NewBgpOpenMessage() *BgpOpenMessage {

}
