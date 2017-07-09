package swim

import (
	"math/rand"
	"net"
	"strconv"
	"time"
)

func joinHostPort(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}

func randInt(n int) int {
	if n == 0 {
		return n
	}

	return (rand.Int() % n)
}

func tick(interval time.Duration, fn func()) {
	t := time.NewTicker(interval)

	for {
		select {
		case <-t.C:
			fn()
			// shutdown channel
		}
	}
}
