package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

type DNSHeader struct {
	ID      uint16
	QR      uint8
	OPCODE  uint8
	AA      uint8
	TC      uint8
	RD      uint8
	RA      uint8
	Z       uint8
	RCODE   uint8
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

type DNSQuestion struct {
	NAME  string
	TYPE  uint16
	CLASS uint16
}

type DNSAnswer struct {
	NAME     string
	TYPE     uint16
	CLASS    uint16
	TTL      uint32
	RDLENGTH uint16
	RDATA    []byte
}

type DNSMessage struct {
	DNSHeader *DNSHeader
	Questions []*DNSQuestion
	Answers   []*DNSAnswer
}

func (m *DNSMessage) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	// Read Header
	header := &DNSHeader{}
	if err := binary.Read(buf, binary.BigEndian, &header.ID); err != nil {
		return err
	}
	flags := uint16(0)
	if err := binary.Read(buf, binary.BigEndian, &flags); err != nil {
		return err
	}
	header.QR = uint8(flags >> 15)
	header.OPCODE = uint8((flags >> 11) & 0xF)
	header.AA = uint8((flags >> 10) & 0x1)
	header.TC = uint8((flags >> 9) & 0x1)
	header.RD = uint8((flags >> 8) & 0x1)
	header.RA = uint8((flags >> 7) & 0x1)
	header.Z = uint8((flags >> 4) & 0x7)
	header.RCODE = uint8(flags & 0xF)

	if err := binary.Read(buf, binary.BigEndian, &header.QDCOUNT); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &header.ANCOUNT); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &header.NSCOUNT); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &header.ARCOUNT); err != nil {
		return err
	}
	m.DNSHeader = header

	// Read Questions
	m.Questions = make([]*DNSQuestion, header.QDCOUNT)
	for i := uint16(0); i < header.QDCOUNT; i++ {
		name, err := readDNSName(buf)
		if err != nil {
			return err
		}
		var qType, qClass uint16
		if err := binary.Read(buf, binary.BigEndian, &qType); err != nil {
			return err
		}
		if err := binary.Read(buf, binary.BigEndian, &qClass); err != nil {
			return err
		}
		m.Questions[i] = &DNSQuestion{
			NAME:  name,
			TYPE:  qType,
			CLASS: qClass,
		}
	}

	// For simplicity, Answers parsing is omitted in this example.

	return nil
}

func (m *DNSMessage) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}

	// Write Header
	if err := binary.Write(buf, binary.BigEndian, m.DNSHeader.ID); err != nil {
		return nil, err
	}
	flags := uint16(m.DNSHeader.QR)<<15 |
		uint16(m.DNSHeader.OPCODE)<<11 |
		uint16(m.DNSHeader.AA)<<10 |
		uint16(m.DNSHeader.TC)<<9 |
		uint16(m.DNSHeader.RD)<<8 |
		uint16(m.DNSHeader.RA)<<7 |
		uint16(m.DNSHeader.Z)<<4 |
		uint16(m.DNSHeader.RCODE)
	if err := binary.Write(buf, binary.BigEndian, flags); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, m.DNSHeader.QDCOUNT); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, m.DNSHeader.ANCOUNT); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, m.DNSHeader.NSCOUNT); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, m.DNSHeader.ARCOUNT); err != nil {
		return nil, err
	}

	// Write Questions
	for _, question := range m.Questions {
		if err := writeDNSName(buf, question.NAME); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.BigEndian, question.TYPE); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.BigEndian, question.CLASS); err != nil {
			return nil, err
		}
	}

	// For simplicity, Answers marshaling is omitted in this example.

	return buf.Bytes(), nil
}
func readDNSName(buf *bytes.Buffer) (string, error) {
	var nameParts []string
	for {
		length, err := buf.ReadByte()
		if err != nil {
			return "", err
		}
		if length == 0 {
			break
		}
		part := make([]byte, length)
		if _, err := buf.Read(part); err != nil {
			return "", err
		}
		nameParts = append(nameParts, string(part))
	}
	return strings.Join(nameParts, "."), nil
}

func writeDNSName(buf *bytes.Buffer, name string) error {
	parts := bytes.Split([]byte(name), []byte("."))
	for _, part := range parts {
		if err := buf.WriteByte(byte(len(part))); err != nil {
			return err
		}
		if _, err := buf.Write(part); err != nil {
			return err
		}
	}
	return buf.WriteByte(0) // End with zero-length byte
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
		size, _, err := udpConn.ReadFromUDP(buf)
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
		fmt.Printf("Received DNS message: %+v\n", receivedMessage)
	}
}
