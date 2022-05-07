package p2p

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/L0rd1k/uprise-api/torrent/client"
	"github.com/L0rd1k/uprise-api/torrent/msg"
	"github.com/L0rd1k/uprise-api/torrent/peer"
)

const peerSize int = 20
const hashSize int = 20

type Torrent struct {
	Peers       []peer.Peer
	PeerId      [peerSize]byte
	InfoHash    [hashSize]byte
	PieceHashes [][hashSize]byte
	PieceLength int
	Length      int
	Name        string
}

type pieceWork struct {
	index  int
	hash   [hashSize]byte
	length int
}

type pieceResult struct {
	index  int
	buffer []byte
}

type pieceProgress struct {
	index     int
	client    *client.Client
	buf       []byte
	dowloaded int
	requested int
	backlog   int
}

func (state *pieceProgress) readMessage() error {
	mesg, err := state.client.Read()
	if err != nil {
		return err
	}
	if mesg == nil {
		return nil
	}
	switch mesg.Id {
	case msg.MsgUnchoke:
		state.client.Choked = false
	case msg.MsgChoke:
		state.client.Choked = true
	case msg.MsgHave:
		index, err := msg.ParseHave(mesg)
		if err != nil {
			return err
		}
		state.client.Bits.SetBit(index)
	case msg.MsgPiece:
		n, err := msg.ParsePiece(state.index, state.buf, mesg)
		if err != nil {
			return err
		}
		state.dowloaded += n
		state.backlog--
	}
	return nil
}

func (t *Torrent) calculatePiecesRange(idx int) (begin int, end int) {
	begin = idx * t.PieceLength
	end = begin + t.PieceLength
	if end > t.Length {
		end = t.Length
	}
	return begin, end
}

func (t *Torrent) calculatePieceSize(idx int) int {
	begin, end := t.calculatePiecesRange(idx)
	return end - begin
}

func (t *Torrent) Download() ([]byte, error) {
	queue := make(chan *pieceWork, len(t.PieceHashes))
	results := make(chan *pieceResult)
	for idx, hash := range t.PieceHashes {
		length := t.calculatePieceSize(idx)
		queue <- &pieceWork{idx, hash, length}
	}
	for _, peer := range t.Peers {
		go t.startDownload(peer, queue, results)
	}

	buf := make([]byte, t.Length)
	donePieces := 0
	for donePieces < len(t.PieceHashes) {
		res := <-results
		begin, end := t.calculatePiecesRange(res.index)
		copy(buf[begin:end], res.buffer)
		donePieces++

		percent := float64(donePieces) / float64(len(t.PieceHashes)) * 100
		numWorkers := runtime.NumGoroutine() - 1
		log.Printf("(%0.2f%%) Downloaded piece #%d from %d peers\n", percent, res.index, numWorkers)
	}
	close(queue)
	return buf, nil
}

func (t *Torrent) startDownload(peer peer.Peer, queue chan *pieceWork, results chan *pieceResult) {
	client, err := client.New(peer, t.PeerId, t.InfoHash)
	if err != nil {
		log.Printf("Could not handshake with %s. Disconnecting\n", peer.Ip)
	}
	defer client.Connection.Close()

	client.SendUnchoke()
	client.SendInterested()

	for piece := range queue {
		if !client.Bits.HasBit(piece.index) {
			queue <- piece //> put piece to queue
			continue
		}
		buf, err := tryDownloadPiece(client, piece)
		if err != nil {
			log.Println("Exiting, err")
			queue <- piece
			return
		}

		err = checkIntegrity(piece, buf)
		if err != nil {
			log.Printf("Piece #%d failed integrity check\n", piece.index)
			queue <- piece
			continue
		}

		client.SendHave(piece.index)
		results <- &pieceResult{piece.index, buf}
	}
}

func checkIntegrity(piece *pieceWork, buf []byte) error {
	hash := sha1.Sum(buf)
	if !bytes.Equal(hash[:], piece.hash[:]) {
		return fmt.Errorf("index %d failed integrity check", piece.index)
	}
	return nil
}

func tryDownloadPiece(c *client.Client, ps *pieceWork) ([]byte, error) {
	state := pieceProgress{
		index:  ps.index,
		client: c,
		buf:    make([]byte, ps.length),
	}
	c.Connection.SetDeadline(time.Now().Add(30 * time.Second))
	defer c.Connection.SetDeadline(time.Time{})
	for state.dowloaded < ps.length {
		if !state.client.Choked {
			for state.backlog < 5 && state.requested < ps.length {
				blockSize := 5
				if ps.length-state.requested < blockSize {
					blockSize = ps.length - state.requested
				}
				err := c.SendRequest(ps.index, state.requested, blockSize)
				if err != nil {
					return nil, err
				}
				state.backlog++
				state.requested += blockSize
			}
		}
		err := state.readMessage()
		if err != nil {
			return nil, err
		}
	}
	return state.buf, nil
}
