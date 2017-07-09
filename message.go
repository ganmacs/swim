package swim

type messageType uint8

const (
	pingMsg messageType = iota
	ackMsg
	pingReqMsg

	// state messge
	aliveMsg
	suspectedMsg
	deadMsg

	compoundMsg
)

type ping struct {
	Id   int
	Name string
	From string
}

type ack struct {
	Id   int
	Name string
	From string
}

type pingReq struct {
	Id   int
	From string
	To   string
}

type alive struct {
	Name        string
	From        string
	Incarnation uint64
}

type suspected struct {
	Name        string
	From        string
	Incarnation uint64
}

type dead struct {
	Name        string
	From        string
	Incarnation uint64
}

type messageHandler interface {
	handlePingMsg(*ping)
	handleAckMsg(*ack)
	handlePingReqMsg(*pingReq)
	handleAliveMsg(*alive)
	handleSuspectedMsg(*suspected)
	handleDeadMsg(*dead)
}

func runMessageHandler(mh messageHandler, tp transport) {
	for {
		select {
		case packet := <-tp.PacketCh():
			go handlePacket(mh, packet)
			// shutdonw channel
		}
	}
}

func handlePacket(mh messageHandler, p *packet) {
	switch p.messageType() {
	case pingMsg:
		var p ping
		// decode
		mh.handlePingMsg(&p)
	}
}
