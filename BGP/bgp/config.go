package bgp

// 状態遷移ハンドラはserver/fsmに移行する
// コンフィグはconfigに移行する

import (
	"log"
	"strconv"
	"strings"
)

// Mode
const (
	_ = iota
	Passive
	Active
)

type Config struct {
	Local_as_number   uint16
	Local_ip_address  string
	Remote_as_number  uint16
	Remote_ip_address string
	Mode              int
}

func ModeParseFromStr(s string) int {
	switch s {
	case "active":
		return Active
	case "passive":
		return Passive
	default:
		log.Fatalln("MODE ERROR")
		return 0
	}
}

func ConfigParseFromStr(configMessage string) (*Config, error) {
	confStr := strings.Split(configMessage, " ")
	localAsNumber, err := strconv.ParseUint(confStr[0], 10, 16)
	if err != nil {
		return nil, err
	}
	remoteAsNumber, err := strconv.ParseUint(confStr[2], 10, 16)
	if err != nil {
		return nil, err
	}
	mode := ModeParseFromStr(confStr[4])

	return &Config{
		Local_as_number: uint16(localAsNumber),
		//Local_ip_address:  utils.Iptobyte(confStr[1]),
		Local_ip_address: confStr[1],
		Remote_as_number: uint16(remoteAsNumber),
		//Remote_ip_address: utils.Iptobyte(confStr[3]),
		Remote_ip_address: confStr[3],
		Mode:              mode,
	}, nil
}
