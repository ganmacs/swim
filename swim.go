package swim

import (
	"os"
	"sync"
	"sync/atomic"

	"github.com/ganmacs/swim/logger"
)

var log = logger.NewSimpleLogger(os.Stdout)

type Swim struct {
	Name string

	nodeLock   *sync.RWMutex
	nodeMap    map[string]*Node
	nodeNum    uint32
	nodes      []*Node
	probeIndex int

	Incarnation uint64

	session *session
	config  *Config
}

func New(config *Config) (*Swim, error) {
	sw, err := newSwim(config)
	if err != nil {
		return nil, err
	}

	if err := sw.setAliveState(); err != nil {
		return nil, err
	}

	sw.start()
	return sw, nil
}

func newSwim(config *Config) (*Swim, error) {
	l := config.Logger
	if l != nil {
		log = l
	}

	tp := config.transport
	if tp == nil {
		dtr, err := newTransport(config.BindAddr, config.BindPort)
		if err != nil {
			return nil, err
		}
		tp = dtr
	}

	se := &session{
		transport:   tp,
		mutex:       new(sync.Mutex),
		ackHandlers: make(map[int]*ackHandler),
	}

	ru := &Swim{
		Name:     joinHostPort(config.BindAddr, config.BindPort),
		nodeMap:  make(map[string]*Node),
		nodeLock: new(sync.RWMutex),
		session:  se,
		config:   config,
	}

	return ru, nil
}

func (sw *Swim) Join(addrs string) (int, error) {
	a := &alive{
		Name: addrs,
		From: sw.Name,
	}

	sw.setAliveNode(a)
	return 0, nil // XXX
}

func (sw *Swim) start() {
	log.Infof("Starting Node... %s\n", sw.Name)
	go sw.session.listen(sw)
	go tick(sw.config.ProbeInterval, sw.probe)
}

func (sw *Swim) nextIncanation() uint64 {
	return atomic.AddUint64(&sw.Incarnation, 1)
}

func (sw *Swim) setAliveState() error {
	a := &alive{
		Name:        sw.Name,
		Incarnation: sw.nextIncanation(),
	}

	sw.setAliveNode(a)
	return nil
}

func (sw *Swim) setAliveNode(a *alive) {
	sw.nodeLock.Lock()
	defer sw.nodeLock.Unlock()
	node, ok := sw.nodeMap[a.Name]

	if !ok {
		node = newNode(a.Name)
		sw.nodeMap[node.name] = node

		l := sw.NodeSize()
		offset := randInt(l)
		sw.nodes = append(sw.nodes, node)
		sw.nodes[offset], sw.nodes[l] = sw.nodes[l], sw.nodes[offset]

		atomic.AddUint32(&sw.nodeNum, 1)
	}

	// An old message
	if a.Incarnation < node.incarnation {
		log.Debugf("Receive old <ALIVE> message about %s\n", a.Name)
		return
	}

	if a.Name == sw.Name {
		// XXX
	} else {
		if node.status != aliveState {
			node.asAliveNode()
		}

		if a.Incarnation == node.incarnation && !ok {
			log.Debugf("Receive same incarnation <ALIVE> message about %s\n", a.Name)
			return
		}

		// TODO: re-broadcast

		node.incarnation = a.Incarnation
	}

	// set alive message
}

func (sw *Swim) NodeSize() int {
	return int(atomic.LoadUint32(&sw.nodeNum))
}

func (sw *Swim) handlePingMsg(p *ping) {
	ack := ack{Id: p.Id, Name: p.From, From: sw.Name}
	sw.session.sendAckMsg(p.From, &ack)
}

func (sw *Swim) handleAckMsg(a *ack) {
	sw.session.receiveAck(a)
}

func (sw *Swim) handlePingReqMsg(pr *pingReq) {
}

func (sw *Swim) handleAliveMsg(a *alive) {
}

func (sw *Swim) handleSuspectedMsg(s *suspected) {
}

func (sw *Swim) handleDeadMsg(d *dead) {
}

func (sw *Swim) probe() {
	log.Debug("Start Probing...")

	sw.nodeLock.Lock()
	n := sw.selectNode()
	sw.nodeLock.Unlock()

	if n == nil {
		log.Debug("Available node is not found")
		return
	}

	// FIXME
	p := &ping{Name: n.name, From: sw.Name}
	ch, ok := sw.session.sendPingAndRecvAck(p, sw.config.ProbeTimeout)
	if ok {
		return
	}

	log.Debug("Start <PING_REQ>")
	sw.nodeLock.Lock()
	nodes := sw.selectNodes(sw.config.SwimNodeCount, func(nd *Node) bool {
		if nd.name == sw.Name || nd.name == n.name {
			return false
		}
		return nd.status == aliveState
	})
	sw.nodeLock.Unlock()

	if len(nodes) == 0 {
		log.Error("No available nodes for <PING_REQ>")
		select {
		case _, ok := <-ch:
			if !ok {
				log.Errorf("Timeout, <PING> message, make %s suspect", n.name)
			}
		}
		return
	}

	pr := &pingReq{Id: p.Id, From: sw.Name, To: p.Name}
	if err := sw.session.sendPingPeqAndRecvAck(pr, nodes, ch); err != nil {
		log.Error(err)
		// suspect!
	}
}

func (sw *Swim) selectNodes(n int, fn func(*Node) bool) (nodes []*Node) {
	v := 0
	for _, node := range sw.nodeMap {
		if fn(node) {
			v += 1
			nodes = append(nodes, node)
		}

		if v >= n {
			return
		}
	}
	return
}

func (sw *Swim) selectNode() *Node {
	tryCount := 0

START:
	if tryCount > sw.NodeSize() {
		return nil
	}

	tryCount++

	if sw.probeIndex == sw.NodeSize() {
		sw.probeIndex = 0
	}

	n := sw.nodes[sw.probeIndex]
	sw.probeIndex++

	if n.name == sw.Name {
		goto START
	}

	if n.status == deadState {
		goto START
	}

	return n
}
