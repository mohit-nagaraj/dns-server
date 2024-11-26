package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		response := make([]byte, 28) // 12 bytes for header + 16 bytes for question
		copy(response, buf[:12])
		response[2] = flipIndicator(response[2])

		// Set the ID to 1234
		binary.BigEndian.PutUint16(response[0:2], 1234)

		// Set QDCOUNT to 1
		binary.BigEndian.PutUint16(response[4:6], 1)

		// Encode the domain name "codecrafters.io"
		domainName := []byte{0x0c, 'c', 'o', 'd', 'e', 'c', 'r', 'a', 'f', 't', 'e', 'r', 's', 0x02, 'i', 'o', 0x00}
		copy(response[12:], domainName)

		// Set Type to 1 (A record)
		binary.BigEndian.PutUint16(response[28:30], 1)

		// Set Class to 1 (IN)
		binary.BigEndian.PutUint16(response[30:32], 1)

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func flipIndicator(b byte) byte {
	return b | 0b10000000
}
