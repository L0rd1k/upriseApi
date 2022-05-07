package client

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/L0rd1k/uprise-api/torrent/bits"
	"github.com/L0rd1k/uprise-api/torrent/handshake"
	"github.com/L0rd1k/uprise-api/torrent/msg"
	"github.com/L0rd1k/uprise-api/torrent/peer"
)

type Client struct {
	Connection net.Conn
	Choked     bool
	Bits       bits.Bits
	peer       peer.Peer
	hash       [20]byte
	peerId     [20]byte
}

func New(peer peer.Peer, peerId, infoHash [20]byte) (*Client, error) {
	//> Set tcp connection timout 3 Seconds
	connection, err := net.DialTimeout("tcp", peer.String(), 20*time.Second)
	if err != nil {
		log.Fatal("Error set tcp connection timeout")
		return nil, err
	}

	//> Complete handshake
	_, err = performHandshake(connection, infoHash, peerId)
	if err != nil {
		//> Close connection if there is an error
		connection.Close()
		return nil, err
	}
	//> Read received data from peer
	buffer, err := recieveBits(connection)
	if err != nil {
		connection.Close()
		return nil, err
	}

	return &Client{
		Connection: connection,
		Choked:     true,
		Bits:       buffer,
		peer:       peer,
		hash:       infoHash,
		peerId:     peerId,
	}, nil
}

func recieveBits(connection net.Conn) (bits.Bits, error) {
	connection.SetDeadline(time.Now().Add(5 * time.Second))
	//> Disable timout at the end of the function
	defer connection.SetDeadline(time.Time{})
	message, err := msg.Read(connection)
	if err != nil {
		return nil, err
	}
	if message == nil {
		err := fmt.Errorf("bitfield is empty")
		return nil, err
	}
	if message.Id != msg.MsgBitField {
		err := fmt.Errorf("wrong bitfield")
		return nil, err
	}
	return message.Payload, nil
}

func performHandshake(connection net.Conn, infoHash, peerId [20]byte) (*handshake.Handshake, error) {
	//> Timeout for all input-output operation
	connection.SetDeadline(time.Now().Add(3 * time.Second))
	//> Disable timout at the end of the function
	defer connection.SetDeadline(time.Time{})
	//> Create handshake
	request := handshake.New(infoHash, peerId)
	//> Write handshake data
	_, err := connection.Write(request.Serialize())
	if err != nil {
		return nil, err
	}
	//> Resive output handshake
	response, err := handshake.Read(connection)
	if err != nil {
		return nil, err
	}
	// Make sure that hashes are equal
	if !bytes.Equal(response.InfoHash[:], infoHash[:]) {
		return nil, fmt.Errorf("info hashes not equal")
	}
	return response, nil
}

func (c *Client) Read() (*msg.Message, error) {
	//> Read message from the connection
	message, err := msg.Read(c.Connection)
	return message, err
}

func (c *Client) SendRequest(index, begin, length int) error {
	request := msg.FormatRequest(index, begin, length)
	_, err := c.Connection.Write(request.Serialize())
	return err
}

func (c *Client) SendInterested() error {
	msg := msg.Message{Id: msg.MsgInterested}
	_, err := c.Connection.Write(msg.Serialize())
	return err
}

func (c *Client) SendNotInterested() error {
	msg := msg.Message{Id: msg.MsgNotInterested}
	_, err := c.Connection.Write(msg.Serialize())
	return err
}

func (c *Client) SendUnchoke() error {
	msg := msg.Message{Id: msg.MsgUnchoke}
	_, err := c.Connection.Write(msg.Serialize())
	return err
}

func (c *Client) SendHave(index int) error {
	msg := msg.FormatHave(index)
	_, err := c.Connection.Write(msg.Serialize())
	return err
}
