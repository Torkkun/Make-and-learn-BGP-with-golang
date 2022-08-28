package utils

import (
	"strconv"
	"strings"
)

func Iptobyte(ip string) []byte {
	var ipbyte []byte
	for _, v := range strings.Split(ip, ".") {
		i, _ := strconv.ParseUint(v, 10, 8)
		ipbyte = append(ipbyte, byte(i))
	}
	return ipbyte
}
