package swim

// session is correct name?
type session struct {
	transport transport
}

type messageHandler interface {
	handlePingMsg(*ping)
	handleAckMsg(*ack)
	handlePingReqMsg(*pingReq)
	handleAliveMsg(*alive)
	handleSuspectedMsg(*suspected)
	handleDeadMsg(*dead)
}

func (se *session) listen(mh messageHandler) {
	for {
		select {
		case packet := <-se.transport.PacketCh():
			go se.handlePacket(mh, packet)
			// shutdonw channel
		}
	}
}

func (se *session) handlePacket(mh messageHandler, pack *packet) {
	switch pack.messageType() {
	case pingMsg:
		// var p ping
		// if err := Decode(pack.body(), &p); err != nil {
		// 	log.Error(err)
		// 	return
		// }
		// log.Debugf("Receive <PING> message from %s", p.Name)
		// mh.handlePingMsg(&p)
	}
}
