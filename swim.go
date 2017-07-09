package swim

import (
	"os"
	"sync"

	"github.com/ganmacs/swim/logger"
)

type Swim struct {
	Name string

	nodeLock *sync.RWMutex
	nodeMap  map[string]*Node
	nodes    []*Node

	Incarnation uint64

	transport *Transport
	config    *Config

	logger *logger.Logger
}

func New(config *Config) (*Swim, error) {
	sw, err := newSwim(config)
	if err != nil {
		return nil, err
	}

	sw.Start()
	return sw, nil
}

func newSwim(config *Config) (*Swim, error) {
	log := config.Logger
	if log == nil {
		log = logger.NewSimpleLogger(os.Stdout)
	}

	tp := config.transport
	if tp == nil {
		dtr, err := newTransport(config.BindAddr, config.BindPort, log)
		if err != nil {
			return nil, err
		}
		tp = dtr
	}

	ru := &Swim{
		Name:      joinHostPort(config.BindAddr, config.BindPort),
		nodeMap:   make(map[string]*Node),
		nodeLock:  new(sync.RWMutex),
		transport: tp,
		logger:    log,
		config:    config,
	}

	return ru, nil
}

func (sw *Swim) Join(addrs string) (int, error) {
	a := &alive{
		Name: addrs,
		From: sw.Name,
	}

	sw.readAliveMessage(a)
	return 0, nil // XXX
}

func (sw *Swim) readAliveMessage(amsg *alive) {
	// set alive message
}

func (sw *Swim) Start() {
	sw.logger.Infof("Starting Node... %s\n", sw.Name)
	go runMessageHandler(sw, sw.transport)
}

func (sw *Swim) handlePingMsg(p *ping) {
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
