package swim

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type ackHandler struct {
	ch    chan bool
	timer *time.Timer
}

// session is correct name?
type session struct {
	transport   transport
	mutex       *sync.Mutex
	ackHandlers map[int]*ackHandler

	// seqId is used for Id  of PING
	seqId int32
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

func (se *session) nextId() int {
	return int(atomic.AddInt32(&se.seqId, 1))
}

func (se *session) handlePacket(mh messageHandler, pack *packet) {
	switch pack.messageType() {
	case pingMsg:
		var p ping
		if err := Decode(pack.body(), &p); err != nil {
			log.Error(err)
			return
		}
		log.Debugf("Receive <PING> message from %s", p.Name)
		mh.handlePingMsg(&p)
	case ackMsg:
		var a ack
		if err := Decode(pack.body(), &a); err != nil {
			log.Error(err)
			return
		}
		log.Debugf("Receive <ACK> message from %s", a.Name)
		mh.handleAckMsg(&a)
	}
}

func (se *session) sendAckMsg(toName string, a *ack) {
	if err := se.sendMessage(toName, ackMsg, a); err != nil {
		log.Error(err)
	}
}

func (se *session) sendPingAndRecvAck(p *ping, timeout time.Duration) (chan bool, bool) {
	p.Id = se.nextId()

	ch := make(chan bool)
	se.setAckHandler(p.Id, ch, timeout*3) // timeout*3 is for waiting PingReq
	if err := se.sendMessage(p.Name, pingMsg, p); err != nil {
		log.Error(err)
		return nil, true
	}

	select {
	case _, ok := <-ch:
		if ok {
			log.Debug("PING is success")
			return nil, true
		} else {
			log.Error("Timeout, PING is failure")
		}
	case <-time.After(timeout):
		log.Error("Timeout, PING is failure")
	}

	return ch, false
}

func (se *session) receiveAck(a *ack) {
	v, ok := se.ackHandlers[a.Id]
	if !ok {
		log.Errorf("Unknow Ack Id: %d", a.Id)
		return
	}
	v.ch <- true
}

func (se *session) sendPingPeqAndRecvAck(pr *pingReq, nodes []*Node, ch chan bool) error {
	for _, n := range nodes {
		log.Debugf("Send <PING_REQ> messaeg to %s", n.name)
		if err := se.sendMessage(n.name, pingReqMsg, pr); err != nil {
			return err
		}
	}

	select {
	case _, ok := <-ch:
		if ok {
			log.Debug("PING is success")
			return nil
		} else {
			return errors.New("Timeout, PING is failure")
		}
	}
}

func (se *session) setAckHandler(id int, ch chan bool, timeout time.Duration) {
	ah := &ackHandler{ch: ch}
	ah.timer = time.AfterFunc(timeout, func() {
		se.mutex.Lock()
		delete(se.ackHandlers, id)
		close(ch)
		se.mutex.Unlock()
	})

	se.mutex.Lock()
	se.ackHandlers[id] = ah
	se.mutex.Unlock()
}

func (se *session) sendMessage(addr string, msgType messageType, msg interface{}) error {
	emsg, err := Encode(msgType, msg)
	if err != nil {
		return err
	}

	bmsg := emsg.Bytes()
	return se.transport.sendData(addr, bmsg)
}
