package main

import (
	"fmt"
	"net"
)

func main() {
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
		var receivedMessage DNSMessage
		err = receivedMessage.UnmarshalBinary(buf[:size])
		if err != nil {
			fmt.Println("Failed to unmarshal received DNS message:", err)
			continue
		}
		// fmt.Printf("Received DNS message: %+v\n", receivedMessage)
		var responseMessage = DNSMessage{
			DNSHeader: &DNSHeader{
				ID:      receivedMessage.ID,
				QR:      1,
				OPCODE:  0,
				OPCODE:  receivedMessage.OPCODE,
				AA:      0,
				TC:      0,
				RD:      0,
				RD:      receivedMessage.RD,
				RA:      0,
				Z:       0,
				RCODE:   0,
				QDCOUNT: 1,
				ANCOUNT: 1,
				NSCOUNT: 0,
				ARCOUNT: 0,
			},
		}
		if responseMessage.OPCODE != 0 {
			responseMessage.RCODE = 4
		}
		responseMessage.Questions = make([]*DNSQuestion, 1)
		responseMessage.Questions[0] = &DNSQuestion{
			NAME:  receivedMessage.Questions[0].NAME,
			TYPE:  1,
			CLASS: 1,
		}
		responseMessage.Answers = make([]*DNSAnswer, 1)
		responseMessage.Answers[0] = &DNSAnswer{
			NAME:     receivedMessage.Questions[0].NAME,
			TYPE:     1,
			CLASS:    1,
			TTL:      60,
			RDLENGTH: 4,
			RDATA:    []byte{8, 8, 8, 8},
		}
		response, err := responseMessage.MarshalBinary()
		if err != nil {
			fmt.Println("Failed to marshal response DNS message:", err)
			continue
		}
		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
