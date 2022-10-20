package bgp

import (
	"encoding/binary"
	"fmt"
	"net"
)

type BGPMessage struct {
	Header BgpMessageHeader
	Body   BGPBody
}

type BGPBody interface {
	Deserialize([]byte) error   // packetからstructにする
	Serialize() ([]byte, error) // packetにする
}

// Message type
const (
	_ = iota
	MSG_Open
	MSG_UPDATE
	MSG_NOTIFICATION
	MSG_KEEPALIVE
)

// RFC
// https://www.rfc-editor.org/rfc/rfc4271#section-4.1
/*
0                   1                   2                   3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                                                               |
+                                                               +
|                                                               |
+                                                               +
|                           Marker                              |
+                                                               +
|                                                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          Length               |      Type     |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ */

type BgpMessageHeader struct {
	Marker []byte // 16-octet
	Length uint16 // 2-octet
	Type   uint8  // 1-octet
}

const (
	BGP_HEADER_LENGTH      = 10
	BGP_MAX_MESSAGE_LENGTH = 4096
)

func (bgpMsgH *BgpMessageHeader) Serialize() ([]byte, error) {
	buf := make([]byte, BGP_HEADER_LENGTH)
	// BGPコネクションに認証を使用する場合に使う。もし、認証を使用しない場合はすべて"1"が入る。
	for i := range buf[:16] {
		buf[i] = 0xff
	}
	// ヘッダを含むメッセージ全体の長さをオクテット単位で定義されます。フィールドの値は最小で19byte、最大で4096byte
	binary.BigEndian.PutUint16(buf[16:18], bgpMsgH.Length)
	// メッセージのタイプコードが定義されます。
	//1：OPEN
	//2：UPDATE
	//3：NOTIFICATION
	//4：KEEPALIVE
	buf[18] = bgpMsgH.Type
	return buf, nil
}

func (bgpMsgH *BgpMessageHeader) Deserialize(data []byte) error {
	if uint16(len(data)) < BGP_HEADER_LENGTH {
		return fmt.Errorf("not all BGP message header")
	}
	bgpMsgH.Length = binary.BigEndian.Uint16(data[16:18])
	if int(bgpMsgH.Length) < BGP_HEADER_LENGTH {
		return fmt.Errorf("unknown message type")
	}
	bgpMsgH.Type = data[18]
	return nil
}

// BGP OPEN MESSEAGE
// GoBGP
// https://github.com/osrg/gobgp/blob/da00912b2fbc96ed272947a4005afff15826e526/pkg/packet/bgp/bgp.go#L1118
// RFC
// https://www.rfc-editor.org/rfc/rfc4271#section-4.2
/*
0                   1                   2                   3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+
|    Version    |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|     My Autonomous System      |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|           Hold Time           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                         BGP Identifier                        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| Opt Parm Len  |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                                                               |
|             Optional Parameters (variable)                    |
|                                                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ */

type BgpOpenMessage struct {
	Version            uint8 // 現在のBGPはRFC4271で定義されるBGP-4を指している
	MyAutonomousSystem uint16
	// HoldTime
	// https://www.infraexpert.com/study/bgpz07.html
	HoldTime                uint16
	BgpIdentifier           net.IP
	OptionalParameterLength uint8
	OptionalParameters      []uint8 // 今のところ実装と使用をしていない
}

func NewBgpOpenMessage(myas uint16, holdtime uint16, id string, optparams []uint8) *BGPMessage {
	return &BGPMessage{
		Header: BgpMessageHeader{Type: MSG_Open},
		Body:   &BgpOpenMessage{4, myas, holdtime, net.ParseIP(id).To4(), 0, optparams},
	}
}

func (opmsg *BgpOpenMessage) Serialize() ([]byte, error) {
	// versionからOptparameterLengthまでをバッファに入れる
	buf := make([]byte, 10)
	buf[0] = opmsg.Version
	binary.BigEndian.PutUint16(buf[1:3], opmsg.MyAutonomousSystem)
	binary.BigEndian.PutUint16(buf[3:5], opmsg.HoldTime)
	copy(buf[5:9], opmsg.BgpIdentifier.To4()) // 4byte分をコピー
	// optparamは現在使用していないので、そのまま何もせずにlengthを返す
	opmsg.OptionalParameterLength = uint8(len(opmsg.OptionalParameters))
	buf[9] = opmsg.OptionalParameterLength
	return append(buf, opmsg.OptionalParameters...), nil
}

func (opmsg *BgpOpenMessage) Deserialize(data []byte) error {
	if len(data) < 10 {
		// 今はとりあえずのエラーメッセージ
		return fmt.Errorf("not all BGP Open message bytes available")
	}
	opmsg.Version = data[0]
	opmsg.MyAutonomousSystem = binary.BigEndian.Uint16(data[1:3])
	opmsg.HoldTime = binary.BigEndian.Uint16(data[3:5])
	opmsg.BgpIdentifier = net.IP(data[5:9]).To4()
	opmsg.OptionalParameterLength = data[9]
	data = data[10:]
	if len(data) < int(opmsg.OptionalParameterLength) {
		return fmt.Errorf("not all BGP Open message bytes available")
	}
	// optional parameterの処理が入るが省略
	return nil
}
