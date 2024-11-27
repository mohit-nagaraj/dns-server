package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type DNSHeader struct {
	ID      uint16
	QR      uint8
	Opcode  uint8
	AA      uint8
	TC      uint8
	RD      uint8
	RA      uint8
	Z       uint8
	Rcode   uint8
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

type DNSQuestion struct {
	Name  []byte
	Type  uint16
	Class uint16
}

type DNSMessage struct {
	Header    DNSHeader
	Questions []DNSQuestion
}

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

		responseMessage := DNSMessage{
			Header: DNSHeader{
				ID:      receivedMessage.Header.ID,
				QR:      1,
				Opcode:  receivedMessage.Header.Opcode,
				AA:      0,
				TC:      0,
				RD:      receivedMessage.Header.RD,
				RA:      0,
				Z:       0,
				Rcode:   0,
				QDCOUNT: 1,
				ANCOUNT: 0,
				NSCOUNT: 0,
				ARCOUNT: 0,
			},
			Questions: receivedMessage.Questions,
		}

		if responseMessage.Header.Opcode != 0 {
			responseMessage.Header.Rcode = 4
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

func (m *DNSMessage) UnmarshalBinary(data []byte) error {
	m.Header.ID = binary.BigEndian.Uint16(data[0:2])
	m.Header.QR = data[2] >> 7
	m.Header.Opcode = (data[2] >> 3) & 0x0F
	m.Header.AA = (data[2] >> 2) & 0x01
	m.Header.TC = (data[2] >> 1) & 0x01
	m.Header.RD = data[2] & 0x01
	m.Header.RA = data[3] >> 7
	m.Header.Z = (data[3] >> 4) & 0x07
	m.Header.Rcode = data[3] & 0x0F
	m.Header.QDCOUNT = binary.BigEndian.Uint16(data[4:6])
	m.Header.ANCOUNT = binary.BigEndian.Uint16(data[6:8])
	m.Header.NSCOUNT = binary.BigEndian.Uint16(data[8:10])
	m.Header.ARCOUNT = binary.BigEndian.Uint16(data[10:12])

	offset := 12
	for i := 0; i < int(m.Header.QDCOUNT); i++ {
		var question DNSQuestion
		question.Name, offset = readName(data, offset)
		question.Type = binary.BigEndian.Uint16(data[offset : offset+2])
		question.Class = binary.BigEndian.Uint16(data[offset+2 : offset+4])
		offset += 4
		m.Questions = append(m.Questions, question)
	}

	return nil
}

func (m *DNSMessage) MarshalBinary() ([]byte, error) {
	data := make([]byte, 12)
	binary.BigEndian.PutUint16(data[0:2], m.Header.ID)
	data[2] = m.Header.QR<<7 | m.Header.Opcode<<3 | m.Header.AA<<2 | m.Header.TC<<1 | m.Header.RD
	data[3] = m.Header.RA<<7 | m.Header.Z<<4 | m.Header.Rcode
	binary.BigEndian.PutUint16(data[4:6], m.Header.QDCOUNT)
	binary.BigEndian.PutUint16(data[6:8], m.Header.ANCOUNT)
	binary.BigEndian.PutUint16(data[8:10], m.Header.NSCOUNT)
	binary.BigEndian.PutUint16(data[10:12], m.Header.ARCOUNT)

	for _, question := range m.Questions {
		data = append(data, question.Name...)
		qType := make([]byte, 2)
		binary.BigEndian.PutUint16(qType, question.Type)
		data = append(data, qType...)
		qClass := make([]byte, 2)
		binary.BigEndian.PutUint16(qClass, question.Class)
		data = append(data, qClass...)
	}

	return data, nil
}

func readName(data []byte, offset int) ([]byte, int) {
	var name []byte
	for {
		length := int(data[offset])
		if length == 0 {
			offset++
			break
		}
		name = append(name, data[offset:offset+length+1]...)
		offset += length + 1
	}
	return name, offset
}
