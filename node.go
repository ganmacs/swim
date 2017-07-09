package swim

import (
	"net"
	"strconv"
	"time"
)

type stateType int

const (
	aliveState stateType = iota
	deadState
)

type Node struct {
	addr string
	port int

	state state
}

type state struct {
	kind        stateType
	lastChange  time.Time
	incarnation uint64
}

func (nd *Node) Address() string {
	return net.JoinHostPort(nd.addr, strconv.Itoa(nd.port))
}

func (nd *Node) aliveNode() {
	if nd.state.kind == aliveState {
		return
	}

	nd.state.kind = aliveState
	nd.state.lastChange = time.Now()
}
