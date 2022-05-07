package tor

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/L0rd1k/uprise-api/torrent/p2p"
	"github.com/jackpal/bencode-go"
)

const hashSize int = 20
const peerSize int = 20

var localPort int = 8080

func setLocalPort(port int) {
	localPort = port
}

type Torrent struct {
	Announce        string
	Comment         string
	InfoHash        [hashSize]byte
	InfoPieceHashes [][hashSize]byte
	InfoPieceLength int
	InfoLength      int
	InfoName        string
}

type initBencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type initBencodeMain struct {
	Announce string `bencode:"announce"`
	Comment  string `bencode:"comment"`
	Info     initBencodeInfo
}

func (info *initBencodeInfo) hash() ([hashSize]byte, error) {
	var buffer bytes.Buffer
	error := bencode.Marshal(&buffer, *info)
	if error != nil {
		return [hashSize]byte{}, error
	}
	hashSum := sha1.Sum(buffer.Bytes())
	return hashSum, nil
}

func (info *initBencodeInfo) splitHash() ([][hashSize]byte, error) {
	shaHashLength := hashSize
	buffer := []byte(info.Pieces)
	if (len(buffer) % shaHashLength) != 0 {
		error := fmt.Errorf("error hash length %d", len(buffer))
		return nil, error
	}
	hashNumber := len(buffer) / shaHashLength
	hashes := make([][hashSize]byte, hashNumber)
	for elem := 0; elem < hashNumber; elem++ {
		copy(hashes[elem][:], buffer[elem*shaHashLength:(elem+1)*shaHashLength])
	}
	return hashes, nil
}

func (tor *initBencodeMain) toTorrent() (Torrent, error) {
	hashSum, error := tor.Info.hash()
	if error != nil {
		return Torrent{}, nil
	}
	hashPieces, error := tor.Info.splitHash()
	if error != nil {
		return Torrent{}, error
	}
	t := Torrent{
		Announce:        tor.Announce,
		Comment:         tor.Comment,
		InfoHash:        hashSum,
		InfoPieceHashes: hashPieces,
		InfoPieceLength: tor.Info.PieceLength,
		InfoLength:      tor.Info.Length,
		InfoName:        tor.Info.Name,
	}

	return t, nil
}

func Open(fileName string) (Torrent, error) {
	file, error := os.Open(fileName)
	if error != nil {
		return Torrent{}, nil
	}
	defer file.Close()

	benc := initBencodeMain{}
	error = bencode.Unmarshal(file, &benc)
	if error != nil {
		return Torrent{}, error
	}
	return benc.toTorrent()
}

func (tor *Torrent) SaveToFile(fileName string) error {
	var id [peerSize]byte
	_, err := rand.Read(id[:])
	if err != nil {
		return err
	}

	peers, err := tor.requestPeers(id, uint16(localPort))
	fmt.Println("Peers count: ", len(peers))
	for i := 0; i < len(peers); i++ {
		fmt.Println(peers[i].String())
	}

	if err != nil {
		return err
	}

	torrent := p2p.Torrent{
		Peers:       peers,
		PeerId:      id,
		InfoHash:    tor.InfoHash,
		PieceHashes: tor.InfoPieceHashes,
		PieceLength: tor.InfoPieceLength,
		Length:      tor.InfoLength,
		Name:        tor.InfoName,
	}

	buffer, err := torrent.Download()
	if err != nil {
		return err
	}
	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = outFile.Write(buffer)
	if err != nil {
		return err
	}
	return nil
}
