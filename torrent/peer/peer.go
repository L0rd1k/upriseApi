package peer

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

type Peer struct {
	Ip   net.IP
	Port uint16
}

func Deforme(buffer []byte) ([]Peer, error) {
	const peerSize = 6
	peersNum := len(buffer) / peerSize
	if len(buffer)%peerSize != 0 {
		error := fmt.Errorf("received deformed peers")
		return nil, error
	}
	peers := make([]Peer, peersNum)
	for i := 0; i < peersNum; i++ {
		offset := i * peerSize
		peers[i].Ip = net.IP(buffer[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16([]byte(buffer[offset+4 : offset+6]))
	}
	return peers, nil
}

func (peer Peer) String() string {
	fmt.Println("HostPort: ", peer.Ip.String(), ":", strconv.Itoa(int(peer.Port)))
	return net.JoinHostPort(peer.Ip.String(), strconv.Itoa((int(peer.Port))))
}
