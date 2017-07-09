package swim

import (
	"time"
)

type stateType int

const (
	aliveState stateType = iota
	deadState
)

type Node struct {
	state
	name string
}

type state struct {
	status      stateType
	lastChange  time.Time
	incarnation uint64
}

func newNode(name string, inc uint64) *Node {
	return &Node{
		name: name,
		state: state{
			status:      deadState,
			lastChange:  time.Now(),
			incarnation: inc,
		},
	}
}

func (nd *Node) asAliveNode() {
	if nd.state.status == aliveState {
		return
	}

	nd.state.status = aliveState
	nd.state.lastChange = time.Now()
}
