package messages

import (
	"encoding/binary"
	"errors"
)

type MessageType uint8

var (
	ErrInvalidMessage = errors.New("invalid message")
	ErrBadMessageType = errors.New("message type is wrong")
)

const (
	MessageTypeKeyGen1 MessageType = iota
	MessageTypeKeyGen2
	MessageTypeKeyGen3
	MessageTypeSign1
	MessageTypeSign2
	MessageTypeSign3
)

const HeaderLength = 1
const HeaderLengthFrom = HeaderLength + 4
const HeaderLengthFromTo = HeaderLength + 8

type Message struct {
	Type     MessageType
	From, To uint32
	KeyGen1  *KeyGen1
	KeyGen2  *KeyGen2
	Sign1    *Sign1
	Sign2    *Sign2
	Sign3    *Sign3
}

func (m *Message) MarshalBinary() ([]byte, error) {
	switch {
	case m.KeyGen1 != nil:
		var buf []byte
		buf = make([]byte, HeaderLengthFrom, HeaderLengthFrom+m.KeyGen1.Size())
		buf[0] = byte(MessageTypeKeyGen1)
		binary.BigEndian.PutUint32(buf[1:5], m.From)
		return m.KeyGen1.BytesAppend(buf[:HeaderLengthFrom])
	case m.KeyGen2 != nil:
		var buf [HeaderLengthFromTo + KeyGenSize2]byte
		buf[0] = byte(MessageTypeKeyGen2)
		binary.BigEndian.PutUint32(buf[1:5], m.From)
		binary.BigEndian.PutUint32(buf[5:9], m.To)
		return m.KeyGen2.BytesAppend(buf[:HeaderLengthFromTo])
	case m.Sign1 != nil:
		var buf [HeaderLengthFrom + SignSize1]byte
		buf[0] = byte(MessageTypeSign1)
		binary.BigEndian.PutUint32(buf[1:5], m.From)
		return m.Sign1.BytesAppend(buf[:HeaderLengthFrom])
	case m.Sign2 != nil:
		var buf [HeaderLengthFrom + SignSize2]byte
		buf[0] = byte(MessageTypeSign2)
		binary.BigEndian.PutUint32(buf[1:5], m.From)
		return m.Sign2.BytesAppend(buf[:HeaderLengthFrom])
	case m.Sign3 != nil:
		var buf [HeaderLength + SignSize2]byte
		buf[0] = byte(MessageTypeSign3)
		return m.Sign3.BytesAppend(buf[:1])
	}
	return nil, errors.New("message does not contain any data")
}

func (m *Message) UnmarshalBinary(data []byte) error {
	msgType := MessageType(data[0])
	m.Type = msgType

	switch msgType {
	case MessageTypeKeyGen1:
		m.From = binary.BigEndian.Uint32(data[1:])
		m.KeyGen1 = new(KeyGen1)
		return m.KeyGen1.UnmarshalBinary(data[5:])

	case MessageTypeKeyGen2:
		m.From = binary.BigEndian.Uint32(data[1:])

		m.To = binary.BigEndian.Uint32(data[5:])

		m.KeyGen2 = new(KeyGen2)
		return m.KeyGen2.UnmarshalBinary(data[9:])

	case MessageTypeSign1:
		m.From = binary.BigEndian.Uint32(data[1:])

		m.Sign1 = new(Sign1)
		return m.Sign1.UnmarshalBinary(data[5:])

	case MessageTypeSign2:
		m.From = binary.BigEndian.Uint32(data[1:])

		m.Sign2 = new(Sign2)
		return m.Sign2.UnmarshalBinary(data[5:])

	case MessageTypeSign3:
		m.Sign3 = new(Sign3)
		return m.Sign3.UnmarshalBinary(data[1:])
	}
	return errors.New("message type not recognized")
}