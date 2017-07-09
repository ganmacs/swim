package swim

import (
	"net"
)

type packet struct {
	buf []byte

	from net.Addr
}

func (p *packet) messageType() messageType {
	return messageType(p.buf[0])
}

func (p *packet) body() []byte {
	return p.buf[1:]
}

// func (p *packet) Read() []byte {
// 	return p.buf[1:]
// }

// func (p *packet) Write(data []byte) {
// }
