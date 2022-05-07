package msg

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	MsgChoke         uint8 = 0 //> chokes the receiver
	MsgUnchoke       uint8 = 1 //> unchokes the receiver
	MsgInterested    uint8 = 2 //> expresses interest in receiving data
	MsgNotInterested uint8 = 3 //> expresses disinterest in receiving data
	MsgHave          uint8 = 4 //> alerts the receiver that the sender has downloaded a piece
	MsgBitField      uint8 = 5 //> encodes which pieces that the sender has downloaded
	MsgRequest       uint8 = 6 //> requests a block of data from the receiver
	MsgPiece         uint8 = 7 //> delivers a block of data to fulfill a request
	MsgCancel        uint8 = 8 //> cancels a request
)

type Message struct {
	Id      uint8
	Payload []byte
}

func Read(reader io.Reader) (*Message, error) {
	// Allocate 4 bit for buffer length
	buffLength := make([]byte, 4)
	//> Read data to buffer
	_, err := io.ReadFull(reader, buffLength)
	if err != nil {
		return nil, err
	}
	//> Convert buffer to big endian
	length := binary.BigEndian.Uint32(buffLength)
	if length == 0 {
		return nil, nil
	}
	// Allocate buffer for peer data
	messageBuff := make([]byte, length)
	_, err = io.ReadFull(reader, messageBuff)
	if err != nil {
		return nil, err
	}
	m := Message{
		Id:      uint8(messageBuff[0]),
		Payload: messageBuff[1:],
	}
	return &m, nil
}

func FormatRequest(index, begin, length int) *Message {
	//> Create request message
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return &Message{Id: MsgRequest, Payload: payload}
}

func FormatHave(index int) *Message {
	//> Create have message
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))
	return &Message{Id: MsgHave, Payload: payload}
}

func ParsePiece(index int, buf []byte, msg *Message) (int, error) {
	if msg.Id != MsgPiece {
		return 0, fmt.Errorf("expected PIECE (ID %d), got ID %d", MsgPiece, msg.Id)
	}
	if len(msg.Payload) < 8 {
		return 0, fmt.Errorf("payload too short. %d < 8", len(msg.Payload))
	}
	parsedIndex := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
	if parsedIndex != index {
		return 0, fmt.Errorf("expected index %d, got %d", index, parsedIndex)
	}
	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	if begin >= len(buf) {
		return 0, fmt.Errorf("begin offset too high. %d >= %d", begin, len(buf))
	}
	data := msg.Payload[8:]
	if begin+len(data) > len(buf) {
		return 0, fmt.Errorf("data too long [%d] for offset %d with length %d", len(data), begin, len(buf))
	}
	copy(buf[begin:], data)
	return len(data), nil
}

func ParseHave(msg *Message) (int, error) {
	if msg.Id != MsgHave {
		return 0, fmt.Errorf("expected HAVE (ID %d), got ID %d", MsgHave, msg.Id)
	}
	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("expected payload length 4, got length %d", len(msg.Payload))
	}
	index := int(binary.BigEndian.Uint32(msg.Payload))
	return index, nil
}

func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}
	//> Length
	length := uint32(len(m.Payload) + 1)
	buffer := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buffer[0:4], length)
	//> Message ID
	buffer[4] = byte(m.Id)
	//> Payload
	copy(buffer[5:], m.Payload)
	return buffer
}

func (m *Message) name() string {
	if m == nil {
		return "KeepAlive"
	}
	switch m.Id {
	case MsgChoke:
		return "Choke"
	case MsgUnchoke:
		return "Unchoke"
	case MsgInterested:
		return "Interested"
	case MsgNotInterested:
		return "NotInterested"
	case MsgHave:
		return "Have"
	case MsgBitField:
		return "Bitfield"
	case MsgRequest:
		return "Request"
	case MsgPiece:
		return "Piece"
	case MsgCancel:
		return "Cancel"
	default:
		return fmt.Sprintf("Unknown#%d", m.Id)
	}
}

func (m *Message) String() string {
	if m == nil {
		return m.name()
	}
	return fmt.Sprintf("%s [%d]", m.name(), len(m.Payload))
}
