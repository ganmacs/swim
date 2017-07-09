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

	nodeLock *sync.RWMutex
	nodeMap  map[string]*Node
	nodeNum  uint32
	nodes    []*Node

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

	se := &session{transport: tp}

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
		Name:        addrs,
		From:        sw.Name,
		Incarnation: 0, // TODO
	}

	sw.setAliveNode(a)
	return 0, nil // XXX
}

func (sw *Swim) start() {
	log.Infof("Starting Node... %s\n", sw.Name)
	go sw.session.listen(sw)
}

func (sw *Swim) setAliveState() error {
	a := &alive{
		Name:        sw.Name,
		Incarnation: 0, // TODO:
	}

	sw.setAliveNode(a)
	return nil
}

func (sw *Swim) setAliveNode(amsg *alive) {
	sw.nodeLock.Lock()
	defer sw.nodeLock.Unlock()
	node, ok := sw.nodeMap[amsg.Name]

	if !ok {
		node = newNode(amsg.Name, amsg.Incarnation)
		sw.nodeMap[node.name] = node

		l := sw.NodeSize()
		offset := randInt(l)
		sw.nodes = append(sw.nodes, node)
		sw.nodes[offset], sw.nodes[l] = sw.nodes[l], sw.nodes[offset]

		atomic.AddUint32(&sw.nodeNum, 1)
	}

	// An old message
	if amsg.Incarnation < node.incarnation {
		log.Debugf("Receive old <ALIVE> message about %s\n", amsg.Name)
		return
	}

	if amsg.Name == sw.Name {
		// XXX
	} else {
		if amsg.Incarnation == node.incarnation {
			log.Debugf("Receive same incarnation <ALIVE> message about %s\n", amsg.Name)
			return
		}

		// TODO: re-broadcast

		node.incarnation = amsg.Incarnation
		if node.status != aliveState {
			node.asAliveNode()
		}
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
}

func (sw *Swim) handlePingReqMsg(pr *pingReq) {
}

func (sw *Swim) handleAliveMsg(a *alive) {
}

func (sw *Swim) handleSuspectedMsg(s *suspected) {
}

func (sw *Swim) handleDeadMsg(d *dead) {
}
