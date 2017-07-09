package swim

import (
	"errors"
	"fmt"
	"net"
)

const (
	udpPacketBufSize = 2048
)

type transport interface {
	PacketCh() chan *packet
	sendData(string, []byte) error
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

	go tr.listenUDP()
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

func (tr *Transport) listenUDP() {
	for {
		buf := make([]byte, udpPacketBufSize)

		n, addr, err := tr.udpListener.ReadFromUDP(buf)
		if err != nil {
			log.Error("could not read %v", err)
			continue
		}

		tr.packetCh <- &packet{buf: buf[:n], from: addr}
	}
}

func (tr *Transport) PacketCh() chan *packet {
	return tr.packetCh
}

func (tr *Transport) sendData(addr string, data []byte) error {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return err
	}

	size, err := conn.Write(data)
	if err != nil {
		return err
	}

	if len(data) != size {
		msg := fmt.Sprintf("failed writing data %d/%d\n", size, len(data))
		return errors.New(msg)
	}

	log.Debugf("Sucess sending data %d bytes to %s\n", size, addr)
	return nil
}
