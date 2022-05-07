package msg

import (
	"encoding/binary"
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
