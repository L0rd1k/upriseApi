package tor

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/L0rd1k/uprise-api/torrent/peer"
	"github.com/jackpal/bencode-go"
)

type TorrentResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (tor *Torrent) obtainTrackerQuery(id [20]byte, port uint16) (string, error) {
	base, error := url.Parse(tor.Announce)
	if error != nil {
		return "", error
	}
	params := url.Values{
		"info_hash":  []string{string(tor.InfoHash[:])},
		"peer_id":    []string{string(id[:])},
		"port":       []string{strconv.Itoa(int(localPort))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(tor.InfoLength)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}

func (tor *Torrent) requestPeers(id [20]byte, localPort uint16) ([]peer.Peer, error) {
	query, err := tor.obtainTrackerQuery(id, localPort)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	response, err := client.Get(query)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	trackerResponse := TorrentResponse{}
	err = bencode.Unmarshal(response.Body, &trackerResponse)
	if err != nil {
		return nil, err
	}
	return peer.Deforme([]byte(trackerResponse.Peers))
}
