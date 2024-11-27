package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:2053")
	if err != nil {
		log.Fatal(err.Error())
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer udpConn.Close()
	buf := make([]byte, 512)
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err.Error())
		}
		message := string(buf[:size])
		fmt.Println(message)
		fmt.Println(buf[:2])
		response := []byte{}
		response = append(response, buf[:2]...)
		response = append(response, 0b10000000)
		response = append(response, 0b00000000)
		response = append(response, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00)
		response = binary.BigEndian.AppendUint16(response, uint16(1)) // QDCOUNT
		response = binary.BigEndian.AppendUint16(response, 0x0000)    // ANCOUNT
		response = binary.BigEndian.AppendUint16(response, 0x0000)    // NSCOUNT
		response = binary.BigEndian.AppendUint16(response, 0x0000)    // ARCOUNT
		response = append(response, encodeDomain("codecrafters.io")...)
		response = binary.BigEndian.AppendUint16(response, uint16(1))
		response = binary.BigEndian.AppendUint16(response, uint16(1))
		if _, err := udpConn.WriteToUDP(response, source); err != nil {
			log.Fatal(err.Error())
		}
	}
}
func encodeDomain(domain string) []byte {
	encodes := []byte{}
	for _, seg := range strings.Split(domain, ".") {
		n := len(seg)
		encodes = append(encodes, byte(n))
		encodes = append(encodes, []byte(seg)...)
	}
	return append(encodes, 0x00)
}
func bigendEncode(v int) []byte {
	encodes := []byte{}
	return binary.BigEndian.AppendUint16(encodes, uint16(v))
}
