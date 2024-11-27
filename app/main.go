package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

// Header represents the DNS message header structure
type Header struct {
	ID      uint16
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

// ParseHeader parses the header section of a DNS packet
func ParseHeader(data []byte) *Header {
	return &Header{
		ID:      binary.BigEndian.Uint16(data[0:2]),
		Flags:   binary.BigEndian.Uint16(data[2:4]),
		QDCount: binary.BigEndian.Uint16(data[4:6]),
		ANCount: binary.BigEndian.Uint16(data[6:8]),
		NSCount: binary.BigEndian.Uint16(data[8:10]),
		ARCount: binary.BigEndian.Uint16(data[10:12]),
	}
}

// BuildResponseHeader creates a DNS response header based on the request header
func BuildResponseHeader(requestHeader *Header, rcode uint16) []byte {
	flags := combineFlags(1, // QR (response)
		uint((requestHeader.Flags>>11)&0xF), // OPCODE
		0,                                   // AA (authoritative answer)
		0,                                   // TC (truncation)
		uint((requestHeader.Flags>>8)&0x1),  // RD (recursion desired)
		0,                                   // RA (recursion available)
		0,                                   // Z (reserved)
		uint(rcode))                         // RCODE

	responseHeader := &Header{
		ID:      requestHeader.ID,
		Flags:   flags,
		QDCount: requestHeader.QDCount,
		ANCount: 1, // One answer record
		NSCount: 0,
		ARCount: 0,
	}
	return responseHeader.ToBytes()
}

// ToBytes converts the Header struct into a byte slice
func (h *Header) ToBytes() []byte {
	headerBytes := make([]byte, 12)
	binary.BigEndian.PutUint16(headerBytes[0:2], h.ID)
	binary.BigEndian.PutUint16(headerBytes[2:4], h.Flags)
	binary.BigEndian.PutUint16(headerBytes[4:6], h.QDCount)
	binary.BigEndian.PutUint16(headerBytes[6:8], h.ANCount)
	binary.BigEndian.PutUint16(headerBytes[8:10], h.NSCount)
	binary.BigEndian.PutUint16(headerBytes[10:12], h.ARCount)
	return headerBytes
}

// combineFlags creates the DNS flags field from individual components
func combineFlags(qr, opcode, aa, tc, rd, ra, z, rcode uint) uint16 {
	return uint16(qr<<15 | opcode<<11 | aa<<10 | tc<<9 | rd<<8 | ra<<7 | z<<4 | rcode)
}

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

		if size < 12 {
			fmt.Println("Invalid DNS packet received")
			continue
		}

		// Parse the received DNS packet header
		requestHeader := ParseHeader(buf[:12])
		opcode := (requestHeader.Flags >> 11) & 0xF

		// Determine the response code based on the OPCODE
		var rcode uint16
		if opcode == 0 {
			rcode = 0 // No error
		} else {
			rcode = 4 // Not implemented
		}

		// Build the response header
		responseHeader := BuildResponseHeader(requestHeader, rcode)

		// Construct a minimal DNS response (header only)
		response := responseHeader

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
