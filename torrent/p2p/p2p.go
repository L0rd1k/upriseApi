package p2p

import (
	"log"

	"github.com/L0rd1k/uprise-api/torrent/client"
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
	return nil, nil
}

func (t *Torrent) startDownload(peer peer.Peer, queue chan *pieceWork, results chan *pieceResult) {
	client, err := client.New(peer, t.PeerId, t.InfoHash)
	if err != nil {
		log.Printf("Could not handshake with %s. Disconnecting\n", peer.Ip)
	}
	defer client.Connection.Close()
}
