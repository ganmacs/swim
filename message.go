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
