package swim

import (
	"net"
)

type transport interface {
	PacketCh() chan *packet
}

// will be an interface
type Transport struct {
	packetCh   chan *packet
	shutdownCh chan int

	udpListener *net.UDPConn
}

func newTransport(addr string, port int) (*Transport, error) {
	tr := &Transport{
		packetCh:   make(chan (*packet)),
		shutdownCh: make(chan (int)),
	}

	if err := tr.setupUDPListener(addr, port); err != nil {
		log.Error(err)
	}

	return tr, nil
}

func (tr *Transport) setupUDPListener(addr string, port int) error {
	ip := net.ParseIP(addr)
	udpAddr := &net.UDPAddr{IP: ip, Port: port}
	udpLn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Error(err)
		return err
	}

	tr.udpListener = udpLn
	return nil
}

func (tr *Transport) PacketCh() chan *packet {
	return tr.packetCh
}
