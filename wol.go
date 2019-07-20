package wol

import (
	"encoding/binary"
	"fmt"
	"net"
)

type MagicPacket [102]byte

// new magic pkg
func NewMagicPacket(macAddr string) (packet MagicPacket, err error) {
	var mac net.HardwareAddr
	mac, err = net.ParseMAC(macAddr)
	if err != nil {
		return
	}

	if len(mac) != 6 {
		err = fmt.Errorf("invalid MAC address")
		return
	}
	copy(packet[0:], []byte{255, 255, 255, 255, 255, 255})
	offset := 6

	for i := 0; i < 16; i++ {
		copy(packet[offset:], mac)
		offset += 6
	}

	return
}


func (mp MagicPacket) Broadcast() {
	addresses := getAllInternalAddress()
	for _, addr := range addresses {
		sendUDPPacket(mp, addr+":9")
	}
}

func (mp MagicPacket) Send(addr string) error {
	return sendUDPPacket(mp, addr+":9")
}

func (mp MagicPacket) SendPort(addr string, port string) error {
	return sendUDPPacket(mp, addr+":"+port)
}

/*
 utils funcs
*/

// send udp pkg
func sendUDPPacket(mp MagicPacket, addr string) (err error) {
	var conn net.Conn
	conn, err = net.Dial("udp", addr)
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = conn.Write(mp[:])
	return
}

// convert a IP address to the broadcast address
func getIPv4BroadcastAddress(n *net.IPNet) (ip net.IP, err error) {
	if n.IP.To4() == nil {
		err = fmt.Errorf("does not support IPv6 addresses")
		return
	}
	ip = make(net.IP, len(n.IP.To4()))
	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(n.IP.To4())|^binary.BigEndian.Uint32(net.IP(n.Mask).To4()))
	return
}
// get all internal address
func getAllInternalAddress() (result []string) {
	set := map[string]bool{}
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addresses {
		if aIPNet, ok := addr.(*net.IPNet); ok {
			if aIPNet.IP.To4() == nil {
				continue // ipv6 without broadcast addr
			}
			if aIPNet.IP.IsLoopback() {
				continue // loopback
			}
			broadcastIPAddr, err := getIPv4BroadcastAddress(aIPNet)
			if err != nil {
				continue
			}
			if _, exist := set[broadcastIPAddr.String()]; !exist {
				result = append(result, broadcastIPAddr.String())
				set[broadcastIPAddr.String()] = true
			}
		}
	}
	return
}