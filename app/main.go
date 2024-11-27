package main

import "encoding/binary"

type Message struct {
	Header   []byte
	Question []byte
	Answer   []byte
	//Authority
	//Additional
}

func NewMessage() *Message {
	return &Message{
		Header:   make([]byte, 12),
		Question: []byte{},
		Answer:   []byte{},
	}
}
func (m *Message) SetHeader() {
	binary.BigEndian.PutUint16((*m).Header[0:2], 1234)
	binary.BigEndian.PutUint16((*m).Header[2:4], combineFlags(1, 0, 0, 0, 0, 0, 0, 0))
	binary.BigEndian.PutUint16((*m).Header[4:6], 1)
	binary.BigEndian.PutUint16((*m).Header[6:8], 1)
	binary.BigEndian.PutUint16((*m).Header[8:10], 0)
	binary.BigEndian.PutUint16((*m).Header[10:12], 0)
}
func (m *Message) SetQuestion() {
	(*m).Question = []byte("\x0ccodecrafters\x02io\x00")
	(*m).Question = binary.BigEndian.AppendUint16((*m).Question, 1)
	(*m).Question = binary.BigEndian.AppendUint16((*m).Question, 1)
}
func (m *Message) SetAnswer() {
	answer := []byte{}
	answer = append(answer, []byte("\x0ccodecrafters\x02io\x00")...)
	answer = binary.BigEndian.AppendUint16(answer, 1)
	answer = binary.BigEndian.AppendUint16(answer, 1)
	answer = binary.BigEndian.AppendUint32(answer, 60)
	answer = binary.BigEndian.AppendUint16(answer, 4)
	answer = binary.BigEndian.AppendUint32(answer, binary.BigEndian.Uint32([]byte("\x08\x08\x08\x08")))
	(*m).Answer = answer
}
func PrepareMessage() *Message {
	message := NewMessage()
	message.SetHeader()
	message.SetQuestion()
	message.SetAnswer()
	return message
}
func combineFlags(qr, opcode, aa, tc, rd, ra, z, rcode uint) uint16 {
	return uint16(qr<<15 | opcode<<11 | aa<<10 | tc<<9 | rd<<8 | ra<<7 | z<<4 | rcode)
}
