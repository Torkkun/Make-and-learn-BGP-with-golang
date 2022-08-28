package bgp

import (
	"fmt"
	"strconv"
	"strings"
)

type Mode int

const (
	Passive Mode = iota
	Active
)

type Config struct {
	Local_as_number   uint16
	Local_ip_address  string
	Remote_as_number  uint16
	Remote_ip_address string
	Mode              Mode
}

func ModeParseFromStr(s string) (Mode, error) {
	switch s {
	case "active":
		return Active, nil
	case "passive":
		return Passive, nil
	default:
		// 0はactiveなので
		return -1, fmt.Errorf("Mode Parse Error")
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
	mode, err := ModeParseFromStr(confStr[4])
	if err != nil {
		return nil, err
	}
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
