package handshake

import (
	"fmt"
	"io"
)

type Handshake struct {
	PrtMsg   string   //> handshake identificator
	InfoHash [20]byte //> Hash array
	PeerId   [20]byte //> Identificator of current peer
}

//> Create Handshake struct
func New(infoHash, peerId [20]byte) *Handshake {
	// Create Handshake instance
	return &Handshake{
		PrtMsg:   "BitTorrent protocol",
		InfoHash: infoHash,
		PeerId:   peerId,
	}
}

//> Convert Handshake struct to byte representation
func (hsh *Handshake) Serialize() []byte {
	//> create buffer slice
	//> 1 + msg + 8 + 20 + 20
	buffer := make([]byte, len(hsh.PrtMsg)+49)
	//> Protocol Message length
	buffer[0] = byte(len(hsh.PrtMsg))
	currPtr := 1
	//> Protocol message
	currPtr += copy(buffer[currPtr:], []byte(hsh.PrtMsg))
	//> Reserve 8 bytes
	currPtr += copy(buffer[currPtr:], make([]byte, 8))
	//> Info hash slice
	currPtr += copy(buffer[currPtr:], hsh.InfoHash[:])
	//> Peer identification slice
	currPtr += copy(buffer[currPtr:], hsh.PeerId[:])
	return buffer
}

func (hsh *Handshake) Deserialize(reader *io.Reader, buffer []byte) error {
	//> Read message size
	msgLength := int(buffer[0])
	if msgLength == 0 {
		err := fmt.Errorf("protocol message can't be 0")
		return err
	}
	//> retrieve handshake data
	handshakeBuf := make([]byte, 48+msgLength)
	_, err := io.ReadFull(*reader, handshakeBuf)
	if err != nil {
		return err
	}
	var infoHash, peerId [20]byte
	copy(infoHash[:], handshakeBuf[msgLength+8:msgLength+8+20])
	copy(peerId[:], handshakeBuf[msgLength+8+20:])
	hsh.PrtMsg = string(handshakeBuf[0:msgLength])
	hsh.InfoHash = infoHash
	hsh.PeerId = peerId
	return nil
}

func Read(reader io.Reader) (*Handshake, error) {
	//> Create buffer for reading message size
	buff := make([]byte, 1)
	//> Read data length to buffer
	_, err := io.ReadFull(reader, buff)
	if err != nil {
		return nil, err
	}
	//> Create handshake
	hShake := Handshake{}
	//> Deserialize received data
	hShake.Deserialize(&reader, buff)
	return &hShake, nil
}
