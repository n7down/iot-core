package main

import (
	"net"
)

func main() {
	Conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{127, 0, 0, 1}, Port: 10001, Zone: ""})
	defer Conn.Close()
	Conn.Write([]byte("hello"))
}
